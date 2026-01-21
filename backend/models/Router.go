package models

import (
	"github.com/gin-gonic/gin"
)

type Config struct {
	Env string
}

func NewRouter(cfg Config) *gin.Engine {
	if cfg.Env == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	return r
}
