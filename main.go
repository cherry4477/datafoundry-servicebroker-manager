package main

import (
	"github.com/asiainfoldp/datafoundry-servicebroker-manager/log"
	"github.com/asiainfoldp/datafoundry-servicebroker-manager/planapi"
	"github.com/asiainfoldp/datafoundry-servicebroker-manager/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

var port, username, password string
var etcdclient tools.EtcdClient

func init() {
	port = os.Getenv("PLANPORT")
	if len(port) == 0 {
		port = "10000"
	}

	username = tools.Getenv("SEAPIUSER")
	password = tools.Getenv("SEAPIPASSWORD")
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
	gin.SetMode(gin.ReleaseMode)

	//获取路由实例
	router = gin.Default()

	router.Use(gin.BasicAuth(gin.Accounts{
		//"asiainfoLDP": "2016asia",
		username: password,
	}))

	router.GET("/seapi/catalog", planapi.Catalog)
	router.GET("/seapi/services/:service_id", planapi.PollingService)
	router.GET("/seapi/services/:service_id/plans/:plan_id", planapi.PollingPlan)
	router.GET("/seapi/services/:service_id/plans", planapi.PollingPlans)
	router.POST("/seapi/services/:service_id", planapi.ProvisionService)
	router.POST("/seapi/services/:service_id/plans/:plan_id", planapi.ProvisionPlan)

	router.PUT("/seapi/services/:service_id", planapi.UpdataService)
	router.PUT("/seapi/services/:service_id/plans/:plan_id", planapi.UpdataPlan)
	router.DELETE("/seapi/services/:service_id", planapi.DeprovisionService)
	router.DELETE("/seapi/services/:service_id/plans/:plan_id", planapi.DeprovisionPlan)

	return
}
