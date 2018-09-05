package main

import (
	"github.com/asiainfoLDP/servicebroker-plan-api/log"
	"github.com/asiainfoLDP/servicebroker-plan-api/planapi"
	"github.com/asiainfoLDP/servicebroker-plan-api/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

var port string

func init() {
	port = os.Getenv("PLANPORT")
	if len(port) == 0 {
		port = "10000"
	}
}

func main() {
	router := handle()
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 0,
	}
	log.Logger.Info("Service starting ...", map[string]interface{}{"port": port, "time": tools.GetTimeNow()})
	s.ListenAndServe()
}

func handle() (router *gin.Engine) {
	//设置全局环境：1.开发环境（gin.DebugMode） 2.线上环境（gin.ReleaseMode）
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	//获取路由实例
	router = gin.Default()

	//var username, password string
	router.Use(gin.BasicAuth(gin.Accounts{
		//username: password,
		"wuchao": "111111",
	}))

	router.GET("/plan/catalog", planapi.Catalog)
	router.PUT("/v2/service_instances/:instance_id", planapi.Provision)
	router.DELETE("/v2/service_instances/:instance_id", planapi.Deprovision)
	router.PATCH("/v2/service_instances/:instance_id", planapi.Update)
	router.GET("/v2/service_instances/:instance_id/service_bindings/:binding_id", planapi.PollingPlan)

	return
}
