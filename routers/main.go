package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pki-server/controllers"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/v1/pki/healthz", controllers.Healthz)

	// 设置日志格式
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf(`{"time":"%s","dest_ip":"%s","http_method":"%s","uri_path":"%s","proto":"%s","status":%d,"response_time":"%s","http_user_agent":"%s","bytes_in":%d,"errmsg":"%s"}%v`,
			param.TimeStamp.Format(time.UnixDate),
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.BodySize,
			param.ErrorMessage,
			"\r\n",
		)
	}))

	v1 := r.Group("/v1/pki")

	{
		v1.POST("/project", controllers.SignatureProject)
		v1.DELETE("/project/:name/:env", controllers.RemoveProject)
		v1.PUT("/project/:name/:env/:year", controllers.RenewalProject)
		v1.GET("/project", controllers.GetAllProject)
		v1.GET("/project/:name/:env/:filename", controllers.GetProjectFile)
	}

	return r
}
