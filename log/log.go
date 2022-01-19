package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const IsoZonedDateTime = "2006-01-02 15:04:05"

const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

type LogLevel int

const (
	ErrorLevel LogLevel = iota
	WarnLevel
	InfoLevel
	DebugLevel
)

type logger struct {
	Logger                   *log.Logger
	colorful                 bool
	logLevel                 LogLevel
	infoStr, warnStr, errStr string
	prefix                   string
}

// LogMode log mode
func (l *logger) LogMode(level LogLevel) *logger {
	newlogger := *l
	newlogger.logLevel = level
	return &newlogger
}

// Info print info
func (l *logger) Info(msg string, v ...interface{}) {
	if l.logLevel >= InfoLevel {
		l.Logger.Printf(l.infoStr+l.getPrefix()+msg, v...)
	}
}

// Warn print warn messages
func (l *logger) Warn(msg string, v ...interface{}) {
	if l.logLevel >= WarnLevel {
		l.Logger.Printf(l.warnStr+l.getPrefix()+msg, v...)
	}
}

// Error print error messages
func (l *logger) Error(msg string, v ...interface{}) {
	if l.logLevel >= ErrorLevel {
		l.Logger.Printf(l.errStr+l.getPrefix()+msg, v...)
	}
}

func (l *logger) Debug(msg string, v ...interface{}) {
	if l.logLevel >= DebugLevel {
		l.Logger.Printf(l.errStr+l.getPrefix()+msg, v...)
	}
}
func (l *logger) getPrefix() string {
	if l.colorful {
		switch l.logLevel {
		case InfoLevel:
			return Green + l.prefix + Reset
		case WarnLevel:
			return BlueBold + l.prefix + Reset
		case ErrorLevel:
			return Magenta + l.prefix + Reset
		case DebugLevel:
			return Green + l.prefix + Reset
		}
	}
	return l.prefix
}

func New(out io.Writer, logLevel LogLevel, colorful bool) *logger {
	var (
		infoStr = "[info]  "
		warnStr = "[warn]  "
		errStr  = "[error]  "
	)

	if colorful {
		infoStr = Green + "[info] " + Reset
		warnStr = Magenta + "[warn] " + Reset
		errStr = Red + "[error] " + Reset
	}

	return &logger{
		Logger:   log.New(out, "", log.LstdFlags),
		colorful: colorful,
		logLevel: logLevel,
		infoStr:  infoStr,
		warnStr:  warnStr,
		errStr:   errStr,
	}
}

func NewWithPrefix(out io.Writer, logLevel LogLevel, colorful bool, prefix string) *logger {
	l := New(out, logLevel, colorful)

	l.prefix = "[" + prefix + "]  "
	return l
}

var (
	Discard = New(ioutil.Discard, DebugLevel, false)
	std     = New(os.Stderr, DebugLevel, true)
)

func Info(s string, v ...interface{}) {
	std.Info(s, v...)
}

func Debug(s string, v ...interface{}) {
	std.Debug(s, v...)
}

func Warn(s string, v ...interface{}) {
	std.Warn(s, v...)
}

func Error(s string, v ...interface{}) {
	std.Error(s, v...)
}

func Msgf(s string, v ...interface{}) {
	Info(s, v...)
}

func Msg(v ...interface{}) {
	if std.logLevel >= InfoLevel {
		Print(v...)
	}
}

// SetPrefix sets the output prefix for the standard logger.
func SetPrefix(prefix string) {
	std.prefix += "[" + prefix + "]  "
}

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	std.Logger.Output(2, fmt.Sprint(v...))
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	std.Logger.Output(2, fmt.Sprintf(format, v...))
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	std.Logger.Output(2, fmt.Sprintln(v...))
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	std.Logger.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	std.Logger.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	std.Logger.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	std.Logger.Output(2, s)
	panic(s)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	std.Logger.Output(2, s)
	panic(s)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	std.Logger.Output(2, s)
	panic(s)
}
