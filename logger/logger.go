package logger

import (
	"encoding/json"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	go setLoggerConfig()
}

func setLoggerConfig() {

	rawJSON := []byte(`{
	  "level": "debug",
      "encoding": "console",
	  "outputPaths": ["stdout", "./tmp/logs"],
	  "errorOutputPaths": ["stderr", "./tmp/errors"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config

	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	c, _ := cfg.Build()
	Logger = c

	defer Logger.Sync()

}
