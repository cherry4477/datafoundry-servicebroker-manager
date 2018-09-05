package planapi

import (
	"io/ioutil"
	"net/http"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/asiainfoLDP/servicebroker-plan-api/tools"
)

var etcdclient tools.EtcdClient

func Catalog(c *gin.Context) {


}

func Provision(c *gin.Context) {

}


func PollingPlan(c *gin.Context) {

}

func Deprovision(c *gin.Context) {
	ins := c.Param("serviceinstance")
	etcdC := etcdclient.GetEtcdApi()
	etcdC.Delete(context.Background(),"",)
}

func Update(c *gin.Context) {
	ins := c.Param("serviceinstance")
	rBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "metrics": err})
		return
	}
	defer c.Request.Body.Close()
	etcdC := etcdclient.GetEtcdApi()
	etcdC.Update(context.Background(),)
}
