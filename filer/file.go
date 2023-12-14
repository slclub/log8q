package filer

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	unit_size = map[string]int{
		"K": 1024,
		"M": 1024 * 1024,
		"G": 1024 * 1024 * 1024,
	}
	zero_time time.Time
)

// interface implement
var _ Sizer = &outer{}
var _ FileRotator = &outer{}
var _ Rotator = &rotateRule{}

/**
 * file info
 */
type fileInfo struct {
	dir_log  string
	dir_abs  string // 绝对路径
	app_name string // 日志文件名
	ext      string // 文件后缀
}

/**
 * 文件滚动规则
 */
type rotateRule struct {
	create        time.Time
	size          int64  // 按多大开始分隔日志文件
	method        string // time: 按时间，按日； size 按文件大小
	keep_duration int64  // 文件保留时长，天数
	file_index    int64  // 当日文件id
}

/**
 * Outer
 */
type outer struct {
	fi          *fileInfo
	rule        Rotator
	file        *os.File
	file_mu     sync.Locker
	head_create string
	size        int64
	ctx         context.Context
}

// config
type Config struct {
	FileName   string // 日志名字
	Dir        string // 日志分类路径
	RotateTime int64  // 保留天数
	Method     string // time:按时间滚动文件；size：按日志大小滚动日志
	Size       any
	Head       string // 日志头，字符串
}

func New(ctx context.Context, config *Config, rotate Rotator) *outer {
	if config == nil {
		config = &Config{}
	}
	config.Init()
	if rotate == nil {
		rotate = &rotateRule{
			create:        time.Now(),
			method:        config.Method,
			keep_duration: config.RotateTime,
			size:          config.SizeInt64(),
		}
	}

	out := &outer{
		fi: &fileInfo{
			dir_abs:  config.Abs(),
			app_name: config.BaseName(),
			ext:      config.Ext(),
			dir_log:  config.Dir,
		},
		rule:        rotate,
		file:        nil,
		file_mu:     &sync.Mutex{},
		head_create: config.Head,
		size:        0,
		ctx:         ctx,
	}
	out.Init()
	return out
}

// ----------------------------------------------------------------------
// outer
// ----------------------------------------------------------------------
func (self *outer) Init() {
	go func() {
		t := time.NewTicker(time.Minute * 10)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				self.AutoDelete(self.fi.fullpath())
			case <-self.ctx.Done():
				return
			}
		}
	}()
}
func (self *outer) reset() {
	self.size = 0
}

func (self *outer) Close() {
	self.file.Close()
}
func (self *outer) Sync() {
	self.file.Sync()
}

func (self *outer) File() *os.File {
	return self.file
}
func (self *outer) movefile() error {
	self.rule.RotateIncr()
	fullname := self.fi.fullname(self.fi.logName(self.rule))
	if ok, _ := isFileExist(fullname); ok {
		self.rule.RotateIncr()
	}
	err := os.Rename(self.fi.fullname(self.fi.usingName()), self.fi.fullname(self.fi.logName(self.rule)))
	return err
}
func (self *outer) Flush() error {
	self.file.Sync()
	return nil
}

// write data to file
func (self *outer) Write(p []byte) (size int, err error) {

	size = len(p)
	if size == 0 {
		return 0, nil
	}
	err = self.rotate()
	if err != nil {
		return 0, err
	}
	//testcode
	//fmt.Println(" file.Write size:", fl.size)

	self.file.Write(p)
	self.size += int64(size)
	return
}

func (self *outer) Size() int64 {
	return self.size
}

func (self *outer) rotate() error {

	if !self.rule.Check(self) {
		return nil
	}
	self.Move()
	// When we start project, value of fl.t is zero time
	if isSameDate(self.rule.CreatedTime(), zero_time) {
		self.rule.CreatedTime(getFileModifyTime(self.fi.fullname(self.fi.usingName())))
	}
	err := self.Open()
	if err == nil { // 未打开文件才需要走下面的创建步骤；
		return err
	}

	err = self.Create()

	return err
}

func (self *outer) Move() error {
	if self.file != nil {
		self.Flush()
		self.file.Close()
		//fmt.Println("[REMOVE][LOG][PATH]", fl.fullpath(), "[LOG_NAME]", fl.fullname(zero_time), "[MODIFY]",isSameDate(fl.t, time.Now()) )
		self.file = nil
	}
	if self.rule.RotateBool(self) {
		self.movefile()
	}

	return nil
}

func (self *outer) Create() error {
	fi, _, err := self.fi.create_file()
	if err != nil {
		return err
	}
	self.file = fi
	//
	self.reset()
	self.rule.CreatedTime(time.Now())

	self.file.Write([]byte(self.head_create))

	return nil
}

func (self *outer) AutoDelete(readdir string) error {
	rdir, _ := ioutil.ReadDir(readdir)
	expire_time := time.Now().Unix() - self.rule.Keep()
	for _, fi := range rdir {
		full_file_name := filepath.Join(self.fi.fullpath(), fi.Name())
		if fi.IsDir() {
			self.AutoDelete(full_file_name)
		}
		if ok, _ := isFileExist(full_file_name); ok && getFileModifyTime(full_file_name).Unix() < expire_time {
			txt := full_file_name[len(full_file_name)-3:]
			if len(full_file_name) > 3 && (txt == "log" || txt == "bak") {
				os.Remove(full_file_name)
			}
		}
	}
	return nil
}

func (self *outer) Open() error {
	self.file_mu.Lock()
	defer self.file_mu.Unlock()
	fname := self.fi.fullname(self.fi.usingName())
	//fl.t = getFileModifyTime(fl.fullpath())
	if ok, _ := isFileExist(fname); !ok {
		return errors.New("Outer.Open file not exist")
	}
	if !isSameDate(self.rule.CreatedTime(), time.Now()) {
		return errors.New("Outer.Open file time is old")
	}
	if self.file != nil {
		// the log file had already been opened
		return nil
	}
	fi, err := os.OpenFile(fname, os.O_WRONLY|os.O_APPEND, 0666)
	fmt.Println("[OPEN_FILE][OK]", fname, ";", err)
	if err != nil {
		fmt.Println("[OPEN_FILE][ERROR]", fname, ";", err)
		panic(any("Log file maybe permission deny!" + err.Error()))
		//return err
	}
	self.file = fi
	//
	finfo, err := self.file.Stat()
	if err == nil {
		self.size = finfo.Size()
	}
	self.reset()
	return nil
}

// ----------------------------------------------------------------------
// fileInfo methods
// ----------------------------------------------------------------------
func (fl *fileInfo) logName(rule Rotator) string {
	name := fmt.Sprintf("%s.%s", fl.app_name, rule.Name())
	name += fl.ext
	return name
}

func (fl *fileInfo) usingName() string {
	return fl.app_name + fl.ext
}

// create log floder.
func (fl *fileInfo) createLogDir() {

	full_path := fl.fullpath()
	fmt.Println("path print create log dir:", full_path)
	if ok, _ := isFileExist(full_path); ok {
		//fmt.Println("path exist create error:", full_path)
		return
	}
	err := os.MkdirAll(full_path, os.ModePerm)
	if err != nil {
		fmt.Println("MKDIR:", err)
	}
	return
}

func (fl *fileInfo) fullpath() string {

	full_path := fl.dir_abs + "/" + fl.dir_log
	full_path = strings.Replace(full_path, "//", "/", 5)
	return full_path
}

// follow time.Time get a name
func (fl *fileInfo) fullname(filename string) string {
	full_path := fl.fullpath()
	return filepath.Join(full_path, filename)
}

func (fl *fileInfo) create_file() (fn *os.File, fname string, err error) {
	fl.createLogDir()
	name := fl.usingName()
	fname = filepath.Join(fl.fullpath(), name)
	fn, err = os.Create(fname)
	return
}

// ----------------------------------------------------------------------
// rotate methods
// ----------------------------------------------------------------------
// 检测是否可以 滚动 文件检查
func (self *rotateRule) Check(r Sizer) bool {
	if r.File() == nil {
		return true
	}
	switch self.method {
	case "time":
		if isSameDate(self.create, time.Now()) {
			return false
		}
	case "size":
		if r.Size() < self.size {
			return false
		}
	}

	return true
}

func (self *rotateRule) Name() string {
	switch self.method {
	case "time":
		return fmt.Sprintf("%04d%02d%02d", self.create.Year(), self.create.Month(), self.create.Day())
	case "size":
		return fmt.Sprintf("%04d%02d%02d.%d", self.create.Year(), self.create.Month(), self.create.Day(), self.file_index)
	}
	return ""
}

func (self *rotateRule) CreatedTime(args ...time.Time) time.Time {
	if len(args) > 0 {
		self.create = args[0]
	}
	return self.create
}

func (self *rotateRule) Keep() int64 {
	return self.keep_duration
}

func (self *rotateRule) RotateIncr() int64 {
	if isSameDate(self.create, time.Now()) {
		self.file_index = 0
	}
	self.file_index++
	return self.file_index
}

func (self *rotateRule) RotateBool(r Sizer) bool {
	switch self.method {
	case "time":
		if isSameDate(self.create, time.Now()) {
			return false
		}
	case "size":
		if r.Size() < self.size {
			return false
		}
	}
	return true
}

// ----------------------------------------------------------------------
// config
// ----------------------------------------------------------------------

func (config *Config) Init() {
	if config.Method == "" {
		config.Method = "time"
	}

	if config.RotateTime == 0 {
		config.RotateTime = 30 * 86400
	}

	if config.Size == 0 {
		config.Size = int64(unit_size["M"]) * 100
	}
	if config.FileName == "" {
		config.FileName = "log8q.log"
	}
}

func (config *Config) BaseName() string {
	p, err := filepath.Abs(config.FileName)
	if err != nil {
		log.Fatal(err)
	}
	filename := filepath.Base(p)
	return filename[:len(filename)-len(config.Ext())]
}

func (config *Config) Ext() string {

	p, err := filepath.Abs(config.FileName)
	if err != nil {
		log.Fatal(err)
	}
	ext := filepath.Ext(p)
	if ext == "" {
		ext = ".log"
	}
	return ext
}

func (config *Config) Abs() string {
	p, err := filepath.Abs(config.FileName)
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(p)
}

func (config *Config) SizeInt64() int64 {
	switch val := config.Size.(type) {
	case int64:
		return val
	case int:
		return int64(val)
	case string:
		if len(val) == 0 {
			return 0
		}
		m := val[len(val)-1:]
		unit, ok := unit_size[m]
		if !ok {
			d, _ := strconv.Atoi(val)
			return int64(d)
		}
		d, _ := strconv.Atoi(val[:len(val)-1])
		return int64(d * unit)
	}
	return 0
}
