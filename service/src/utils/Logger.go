package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

//日志自定义格式
type LogFormatter struct{}

func IsDirExist(path string) bool {
	p, err := os.Stat(path)
	if err == nil {
		return p.IsDir()
	}
	return false
}

//格式详情
func (s *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("0102-150405.000")
	var file string
	var length int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		length = entry.Caller.Line
	}
	msg := entry.Message
	if strings.HasPrefix(msg, "sql") { //获取gorm sql日志
		args := strings.Split(msg, ":")
		if len(args) > 0 {
			msg = args[1]
			i := strings.Index(msg, "s")
			msg = entry.Message[i+len(args[0])+2:]
			file = filepath.Base(args[0][3:])
			length = 0
		}
	}
	msg = fmt.Sprintf("[%s] [%s:%d] [%s] %s\n", timestamp, file, length, strings.ToUpper(entry.Level.String()), msg)
	return []byte(msg), nil
}

func NewLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}
	logpath := "./logs/server.log"
	if IsDirExist("/data/log") {
		logpath = "/data/log/server.log"
	}
	hook := lumberjack.Logger{
		Filename:   logpath, // 日志文件路径
		MaxSize:    10,      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,      // 日志文件最多保存多少个备份
		MaxAge:     7,       // 文件最多保存多少天
		Compress:   true,    // 是否压缩
	}

	// pathMap := lfshook.PathMap{
	// 	logrus.InfoLevel:  "./info.log",
	// 	logrus.ErrorLevel: "./info.log",
	// }
	Log = logrus.New()
	// Log.Hooks.Add(lfshook.NewHook(
	// 	pathMap,
	// 	&logrus.JSONFormatter{},
	// ))
	Log.SetReportCaller(true)
	// Log.SetFormatter(&logrus.TextFormatter{
	// 	//以下设置只是为了使输出更美观
	// 	DisableColors:   true,
	// 	TimestampFormat: "2006-01-02 15:03:04",
	// })
	Log.SetFormatter(new(LogFormatter))
	Log.SetOutput(&hook)
	return Log
}

func init() {
	Log = NewLogger()
	// Log as JSON instead of the default ASCII formatter.
	//Log.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	//Log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	//Log.SetLevel(log.WarnLevel)
}
