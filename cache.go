package log8q

import (
	"github.com/slclub/go-tips/spinlock"
	"io"
	"math/rand"
	"sync"
	"time"
)

const (
	LEN_BUCKET = 1024 * 8
	USING      = 1
	UN_USING   = 0
)

var (
	_ io.ReadWriter = NewBucket(LEN_BUCKET)
)

type WriteMany interface {
	io.Writer
	WriteMany(...[]byte) (int, error)
}

var _ WriteMany = &Cache{}

type Bucket struct {
	data     []byte
	endoff   int
	offset   int
	lock     sync.Locker
	copyLock sync.Locker
	use      int
	Next     *Bucket
}

type Cache struct {
	list []*Bucket
	pool sync.Pool
	lock sync.Locker
}

func NewBucket(lenght int) *Bucket {
	return &Bucket{
		data:     make([]byte, lenght),
		endoff:   0,
		lock:     spinlock.New(),
		copyLock: &sync.Mutex{},
	}
}

func NewCache(cap, cap_bucket int) *Cache {
	cache := &Cache{
		pool: sync.Pool{
			New: func() any { return NewBucket(cap_bucket) },
		},
		list: []*Bucket{},
		lock: spinlock.New(),
	}
	for i := 0; i < cap; i++ {
		cache.list = append(cache.list, NewBucket(cap_bucket))
	}
	return cache
}

// ---------------------------------------------------
// Cache methods
// ---------------------------------------------------

func (cache *Cache) Write(data []byte) (int, error) {
	bucket := cache.Get()
	if bucket == nil {
		time.Sleep(time.Millisecond)
		return cache.Write(data)
	}
	defer bucket.Use(UN_USING)
	n, err := bucket.Write(data)
	return n, err
}

func (cache *Cache) WriteMany(args ...[]byte) (int, error) {
	bucket := cache.Get()
	if bucket == nil {
		time.Sleep(time.Millisecond)
		return cache.WriteMany(args...)
	}
	defer bucket.Use(UN_USING)
	nn := 0
	for _, v := range args {
		n, _ := bucket.Write(v)
		nn += n
	}
	return nn, nil
}

func (cache *Cache) Read(data []byte) (int, error) {
	n := 0
	//for _, v := range cache.list {
	//	if v.ReadSize() <= 0 || v.Use() == USING {
	//		continue
	//	}
	//	tn, err := v.Read(data[n:])
	//	n += tn
	//	if err != nil {
	//		return n, err
	//	}
	//	if n >= len(data) {
	//		return n, nil
	//	}
	//}
	l := len(cache.list)
	k := rand.Intn(l)
	for i := 0; i < len(cache.list); i++ {
		v := cache.list[(k+i)%l]
		if v.ReadSize() <= 0 || v.Use() == USING {
			continue
		}
		tn, err := v.Read(data[n:])
		n += tn
		if err != nil {
			return n, err
		}
		if n >= len(data) {
			return n, nil
		}
	}
	return n, nil
}

func (cache *Cache) Get() *Bucket {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	for _, v := range cache.list {
		if v.Use() == USING || v.Free() == 0 {
			continue
		}
		v.Use(USING)
		return v
	}
	return nil
}

func (cache *Cache) Size() int {
	n := 0
	for _, v := range cache.list {
		n += v.Size()
	}
	return n
}

func (cache *Cache) ReadSize() int {
	n := 0
	for _, v := range cache.list {
		n += v.ReadSize()
	}
	return n
}

func (cache *Cache) Cap() int {
	n := 0
	for _, v := range cache.list {
		n += v.Cap()
	}
	return n
}

// ---------------------------------------------------
// bucket methods
// ---------------------------------------------------
func (self *Bucket) Write(data []byte) (int, error) {
	//self.Use(USING)
	//defer self.Use(UN_USING)
	switch {
	case self.Free() >= len(data):
		n := copy(self.data[self.endoff:], data)
		self.endoff += n
		return n, nil
	case self.Free() < len(data):
		n := copy(self.data[self.endoff:], data)
		self.endoff += n
		self.Next = NewBucket(self.Cap())
		n1, _ := self.Next.Write(data[n:])
		return n + n1, nil
	}

	return 0, nil
}

func (self *Bucket) Read(p []byte) (int, error) {

	self.Use(USING)
	defer self.Use(UN_USING)
	// read lock 暂时使用相同的 唯一log
	self.copyLock.Lock()
	defer self.copyLock.Unlock()
	n := copy(p, self.data[self.offset:self.endoff])
	self.offset += n

	if self.offset >= self.endoff && self.Next == nil {
		self.reset()
	}

	if n < len(p) && self.Next != nil {
		n1 := copy(p[n:], self.Next.data[self.Next.offset:self.Next.endoff])
		self.Next.offset += n1
		if self.Next.offset >= self.Next.endoff {
			self.Next = nil
			self.reset()
		}
		n += n1
	}
	return n, nil
}

func (self *Bucket) ReadAll() ([]byte, error) {
	n := self.ReadSize()
	data := make([]byte, n)
	_, err := self.Read(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//
//func (self *Bucket) writeOnly(data []byte) (int, error) {
//	n := copy(self.data[self.endoff:], data)
//	return n, nil
//}

func (self *Bucket) Size() int {
	return self.endoff - self.offset
}

func (self *Bucket) ReadSize() int {
	if self == nil {
		return 0
	}
	n := self.endoff - self.offset
	if self.Next == nil {
		return n
	}
	return n + self.Next.ReadSize()
}

func (self *Bucket) Cap() int {
	return cap(self.data)
}

func (self *Bucket) Free() int {
	return self.Cap() - self.endoff
}

func (self *Bucket) Use(us ...int) int {
	self.lock.Lock()
	defer self.lock.Unlock()

	if len(us) > 0 {
		self.use = us[0]
	}
	return self.use
}

func (self *Bucket) reset() {
	self.offset = 0
	self.endoff = 0
}
