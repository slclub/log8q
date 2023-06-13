package filer

import (
	"os"
	"time"
)

type Sizer interface {
	Size() int64
	File() *os.File
}

// rotate 流程接口
type FileRotator interface {
	Create() error
	AutoDelete(dir string) error
	Move() error
	Open() error
}

// rotate 对象接口
type Rotator interface {
	Check(sizer Sizer) bool // 校验是否 滚动日志
	CreatedTime(...time.Time) time.Time
	Keep() int64                 // 日志保持天数
	Name() string                // 日志名字
	RotateIncr() int64           // 滚动文件id ++
	RotateBool(sizer Sizer) bool // rotate move 滚动文件时机
}
