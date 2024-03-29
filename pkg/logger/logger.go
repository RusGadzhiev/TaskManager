package logger

import (
	"log"

	"go.uber.org/zap"
)

func NewZapLogger() *zap.SugaredLogger {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"app.log"}
	config.ErrorOutputPaths = []string{"error.log"}
	var err error
	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	return logger.Sugar()
}

// напиши свои функции с контекстом (у випа)
