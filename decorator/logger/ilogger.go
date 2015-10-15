package logger

// ILogger logger interface
type ILogger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Printf(format string, args ...interface{})

	Warningf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	Panicf(format string, args ...interface{})

	Debug(args ...interface{})

	Info(args ...interface{})

	Print(args ...interface{})

	Warning(args ...interface{})

	Error(args ...interface{})

	Fatal(args ...interface{})

	Panic(args ...interface{})
}
