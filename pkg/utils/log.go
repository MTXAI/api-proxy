package utils

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

const (
	LogFileSoftLink    = "proxy.log"
	GinLogFileSoftLink = "proxy-server.log"
)

var logger *slog.Logger

func init() {
	logger = slog.Default()
}

func InitLog(logDir string, debug bool, stdout bool) error {
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create logs dir failed: %v", err)
	}

	logFilePath := filepath.Join(logDir, LogFileSoftLink)
	writer, err := rotatelogs.New(
		logFilePath+".%Y-%m-%d",
		rotatelogs.WithLinkName(logFilePath),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		return fmt.Errorf("create logger failed: %v", err)
	}

	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				fileSeg := strings.Split(source.File, "/")
				if len(fileSeg) > 0 {
					source.File = fileSeg[len(fileSeg)-1]
				}
				sourceVal := fmt.Sprintf("%s:%d", source.File, source.Line)
				a.Value = slog.StringValue(sourceVal)
			}
			return a
		},
	}

	if debug {
		opts.Level = slog.LevelDebug
	}
	multiWriter := io.MultiWriter(writer)
	if stdout {
		multiWriter = io.MultiWriter(os.Stdout, writer)
	}

	logger = slog.New(slog.NewTextHandler(multiWriter, &opts))
	slog.SetDefault(logger)
	return nil
}

func InitGinLog(logDir string, debug bool, stdout bool) error {
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create logs dir failed: %v", err)
	}

	logFilePath := filepath.Join(logDir, GinLogFileSoftLink)
	writer, err := rotatelogs.New(
		logFilePath+".%Y-%m-%d",
		rotatelogs.WithLinkName(logFilePath),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		return fmt.Errorf("create logger failed: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	if debug {
		gin.SetMode(gin.DebugMode)
	}

	multiWriter := io.MultiWriter(writer)
	if stdout {
		multiWriter = io.MultiWriter(os.Stdout, writer)
	}

	gin.DisableConsoleColor()
	gin.DefaultWriter = multiWriter
	gin.DefaultErrorWriter = multiWriter
	return nil
}
