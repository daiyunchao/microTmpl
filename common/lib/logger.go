package lib

import (
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var Logger *zap.SugaredLogger

func InitLogger() {
	encoder := GetEncoder()
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level <= zapcore.FatalLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level > zapcore.FatalLevel
	})
	infoWriter := GetLogWriter("./logs/info")
	warnWriter := GetLogWriter("./logs/error")
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, infoWriter, infoLevel),
		zapcore.NewCore(encoder, warnWriter, warnLevel),
	)
	logger := zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
}

func GetEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 修改时间编码器

	// 在日志文件中使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// NewConsoleEncoder 打印更符合人们观察的方式
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func GetLogWriter(filepath string) zapcore.WriteSyncer {
	hook, err := rotatelogs.New(
		filepath+"_%Y%m%d.log",
		rotatelogs.WithLinkName(filepath),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Minute*1),
	)
	if err != nil {
		panic(err)
	}
	return zapcore.AddSync(hook)
}
