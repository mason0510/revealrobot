package log

import (
	"fmt"
	"log"
	"os"
	"path"

	"time"

	"revealrobot/utils/env"

	"runtime"

	"github.com/natefinch/lumberjack"
)

const (
	EMERG = iota
	ALERT
	CRIT
	ERR
	WARN
	NOTICE
	INFO
	DEBUG
)

var LEVELS = map[string]uint{
	"emerg":  EMERG,
	"alert":  ALERT,
	"crit":   CRIT,
	"err":    ERR,
	"warn":   WARN,
	"notice": NOTICE,
	"info":   INFO,
	"debug":  DEBUG,
}

type Logger struct {
	path           string
	log            *log.Logger
	rlog           lumberjack.Logger
	rollingFile    bool
	lastRotateTime time.Time
	level          uint
	pid            []interface{}
}

func NewLogger(path string, level string) (*Logger, error) {
	tlog := new(Logger)

	tlog.path = path
	tlog.lastRotateTime = time.Now()
	tlog.level = LEVELS[level]
	tlog.pid = []interface{}{env.Pid}

	tlog.rlog.Filename = path
	tlog.rlog.MaxSize = 0x1000 * 2 // automatic rolling file on it increment than 2GB
	tlog.rlog.LocalTime = true

	l := log.New(&tlog.rlog, "", log.LstdFlags|log.Lshortfile)

	tlog.log = l

	return tlog, nil
}

func (tlog *Logger) checkRotate() {
	if !tlog.rollingFile {
		return
	}

	n := time.Now()
	if tlog.lastRotateTime.Year() != n.Year() ||
		tlog.lastRotateTime.Month() != n.Month() ||
		tlog.lastRotateTime.Day() != n.Day() {
		tlog.rlog.Rotate()
		tlog.lastRotateTime = n
	}
}

func (tlog *Logger) SetDailyFile() {
	tlog.rollingFile = true
}

func (tlog *Logger) Emerg(format string, v ...interface{}) {

	tlog.checkRotate()

	if tlog.level < EMERG {
		return
	}

	tlog.log.Printf("[EMERG] #%d "+format, append(tlog.pid, v...)...)
}

func (tlog *Logger) Alert(format string, v ...interface{}) {

	tlog.checkRotate()

	if tlog.level < ALERT {
		return
	}

	tlog.log.Printf("[ALERT] #%d "+format, append(tlog.pid, v...)...)
}

func (tlog *Logger) Crit(format string, v ...interface{}) {

	tlog.checkRotate()

	if tlog.level < CRIT {
		return
	}

	tlog.log.Printf("[CRIT] #%d "+format, append(tlog.pid, v...)...)
}

func (tlog *Logger) Err(format string, v ...interface{}) {

	tlog.checkRotate()

	if tlog.level < ERR {
		return
	}

	tlog.log.Printf("[ERROR] #%d "+format, append(tlog.pid, v...)...)
}

func (tlog *Logger) Warn(format string, v ...interface{}) {

	tlog.checkRotate()

	if tlog.level < WARN {
		return
	}

	tlog.log.Printf("[WARN] #%d "+format, append(tlog.pid, v...)...)
}

func (tlog *Logger) Notice(format string, v ...interface{}) {

	tlog.checkRotate()

	if tlog.level < NOTICE {
		return
	}

	tlog.log.Printf("[NOTICE] #%d "+format, append(tlog.pid, v...)...)
}

func (tlog *Logger) Info(format string, v ...interface{}) {

	tlog.checkRotate()

	if tlog.level < INFO {
		return
	}

	tlog.log.Printf("[INFO] #%d "+format, append(tlog.pid, v...)...)
}

func (tlog *Logger) Debug(format string, v ...interface{}) {

	tlog.checkRotate()

	if tlog.level < DEBUG {
		return
	}
	funcName, file, line, ok := runtime.Caller(2)
	if ok {
		//fmt.Println("func name: " + runtime.FuncForPC(funcName).Name())
		//fmt.Printf("file: %s, line: %d\n", file, line)
		_, filename := path.Split(file)
		_, funcshort := path.Split(runtime.FuncForPC(funcName).Name())

		sline := fmt.Sprintf("%s:%s %d", filename, funcshort, line)
		tlog.log.Printf("[DEBUG] "+sline+" #%d "+format, append(tlog.pid, v...)...)
	} else {
		tlog.log.Printf("[DEBUG] #%d "+format, append(tlog.pid, v...)...)
	}

}

var _logger *Logger

func GetDefault() *Logger {
	return _logger
}

func SetDefault(l *Logger) {
	_logger = l
}

func Stdout() {
	l := log.New(os.Stdout, "", log.LstdFlags)
	tlog := new(Logger)
	tlog.log = l
	tlog.level = DEBUG
	tlog.pid = []interface{}{env.Pid}
	SetDefault(tlog)
}

func Emerg(format string, v ...interface{}) {
	_logger.Emerg(format, v...)
}

func Alert(format string, v ...interface{}) {
	_logger.Alert(format, v...)
}

func Crit(format string, v ...interface{}) {
	_logger.Crit(format, v...)
}

func Err(format string, v ...interface{}) {
	_logger.Err(format, v...)
}

func Warn(format string, v ...interface{}) {
	_logger.Warn(format, v...)
}

func Notice(format string, v ...interface{}) {
	_logger.Notice(format, v...)
}

func Info(format string, v ...interface{}) {
	_logger.Info(format, v...)
}

func Debug(format string, v ...interface{}) {
	_logger.Debug(format, v...)
}

func RawLogger() *log.Logger {
	return _logger.log
}

