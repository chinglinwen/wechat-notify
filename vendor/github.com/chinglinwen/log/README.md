# log

log package from upspin. ( by the same Golang authors )

## Usage

### Default level

`Info`
> use flag or environment variable to set different default level through init.
> 
> call **log.SetLevel("debug")** to change level.

### Log

```
log.Print...
log.Debug.Print...
```

### Where info ( added a extra function to the original log package )

Where function usage example:

```
log.Printf("%v some error info here\n", log.Where())
```

## Log file setting ( size, rotate, etc. )

See http://github.com/natefinch/lumberjack

```
log.SetOutput(&lumberjack.Logger{
    Filename:   "/var/log/myapp/foo.log",
    MaxSize:    500, // megabytes
    MaxBackups: 3,
    MaxAge:     28, //days
})
```

## Example code to support turn on debug output on the fly

See [doc/support-debug-on-the-fly](doc/support-debug-on-the-fly.go)

## Other packages considered (but not choose it) 
> (for the record here only, use above one )
* https://github.com/op/go-logging
* https://github.com/hashicorp/logutils
* https://github.com/golang/glog
* https://github.com/go-kit/kit/tree/master/log
* https://github.com/ScottMansfield/nanolog
* https://github.com/sirupsen/logrus
* https://github.com/apex/log
* https://github.com/uber-go/zap
* https://github.com/juju/loggo

## Reference article about log

> I believe that there are **only two things you should log**:
> 
* Things that developers care about when they are developing or debugging software.
* Things that users care about when using your software.
 
> Obviously these are debug and info levels, respectively.

https://dave.cheney.net/2015/11/05/lets-talk-about-logging