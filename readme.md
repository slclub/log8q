# Log8q

## Overview
q = quick

Quick and high performance log system.

- Implementing offical log interface.
- Common and universal loggin methods.
- Customize log output.
- Log level output.

## Install

```go
go get github.com/slclub/log8q
```

## New

- Simple

```go
New(context.Background(), &Config{})
```

- Output logs to file

```go
l8 = log8q.New(context.Background(), &log8q.Config{
    Filename: "log/log8q.log",
})
```

- Output logs to Stdout

```go
l8 := New(context.Background(), &Config{
    Writer: os.Stdout,
})
```

## Record log

Each level of loggin function has two common calling methods. 
They are like fmt.Print and fmt.Printf, whatever it is function name and parameters
have the same format.

- Info
```go
l8.Info("stdout info", "b", "c", "d")
l8.Infof("stdout info name:%v id:%v", "xiaoming", 1)
```

- Debug
- Warn
- Error
- Fatal
- Print

## Customize

The writer of log8q can be replaced by Config.Writer. It is implement the io.Writer.
So, you can replace it with any object had implemented the io.Writer. os.Stdout, os.File etc. 

