package base

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// 日志需要贯穿整个程序，因此把初始化工作放在独立的 init 函数中
func init() {
	// 1. 定义日志的格式，支持两种，一种是文本格式，一种是 json 格式，此处使用文本格式
	// formatter := &log.TextFormatter{}
	//    使用外部 prefixed 日志格式插件
	formatter := &prefixed.TextFormatter{}
	//    使用 prefixed 的强制化日志格式化功能
	formatter.ForceFormatting = true
	//    使用 prefixed 的自定义高亮颜色功能
	formatter.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "red",
		PanicLevelStyle: "red",
		DebugLevelStyle: "blue",
		PrefixStyle:     "cyan",
		TimestampStyle:  "37",
	})
	//    开启日志格式的时间戳支持
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02.15:04:05.000000 "
	log.SetFormatter(formatter)

	// 2。 定义日志级别，logrus 默认日志级别为 info，通过获取环境变量来决定是否打开 debug 日志级别
	//level := os.Getenv("log.debug")
	//if level == "true" {
	log.SetLevel(log.DebugLevel)
	//}

	// 3。 控制台高亮显示，已经被内置在 formatter 中，因此我们通过设置 formatter 参数即可完成设置
	//     强制高亮显示
	formatter.ForceColors = true
	formatter.DisableColors = false

	log.Info("测试日志条目")
	log.Debug("测试日志条目")

	// 4。 日志文件的滚动配置
	logFileSettings()

}

func logFileSettings() {
	// 配置日志输出目录
	logPath, _ := filepath.Abs("./logs")
	log.Info("log dir: %s", logPath)
	// 配置日志文件名
	logFileName := "envelop"
	// 日志文件最大保存时长
	maxAge := time.Hour * 24
	// 日志切割时间间隔
	rotationTime := time.Hour * 1
	// 在文件系统上创建目录
	os.MkdirAll(logPath, os.ModePerm)
	// 联合完成的日志目录
	baseLogPath := path.Join(logPath, logFileName)

	// 设置滚动日志输出
	writer, err := rotatelogs.New(
		strings.TrimSuffix(baseLogPath, ".log")+".%Y%m%d%H.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新的日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", err)
	}
	log.SetOutput(writer)
}
