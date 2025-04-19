package development_iface

import "github.com/gin-gonic/gin"

type LocalContext interface{}

type LocalRunner interface {
	Setup() error
	RegisterApi(r *gin.Engine) error
}
