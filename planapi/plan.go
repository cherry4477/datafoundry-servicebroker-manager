package planapi

import (
	"io/ioutil"
	"net/http"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/asiainfoLDP/servicebroker-plan-api/tools"
	"github.com/coreos/etcd/client"
	//"golang.org/x/net/context"
	"github.com/asiainfoLDP/servicebroker-plan-api/log"
	"strconv"
	"strings"
	"encoding/json"
	"reflect"
	"fmt"
)

const (
	KEY = "/servicebroker/"+log.ServcieBrokerName+"/catalog"
)

var etcdclient tools.EtcdClient

func Catalog(c *gin.Context) {
	catalogRsp := CatalogResponse{}

	//获取catalog信息
	resp, err := etcdclient.GetEtcdApi().Get(context.Background(),
		"/servicebroker/"+log.ServcieBrokerName+"/catalog",
		&client.GetOptions{Recursive: true})
	if err != nil {
		log.Logger.Error("Can not get catalog information from etcd", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "can not get catalog information from etcd"
		c.JSON(500, errinfo)
		return
	} else {
		log.Logger.Debug("Successful get catalog information from etcd. NodeInfo is " + resp.Node.Key)
	}

	for i := 0; i < len(resp.Node.Nodes); i++ {
		log.Logger.Debug("Start to Parse Service " + resp.Node.Nodes[i].Key)
		myService := Service{}
		myService.Id = strings.Split(resp.Node.Nodes[i].Key, "/")[len(strings.Split(resp.Node.Nodes[i].Key, "/"))-1]
		for j := 0; j < len(resp.Node.Nodes[i].Nodes); j++ {
			if !resp.Node.Nodes[i].Nodes[j].Dir {
				lowerkey := strings.ToLower(resp.Node.Nodes[i].Key)
				switch strings.ToLower(resp.Node.Nodes[i].Nodes[j].Key) {
				case lowerkey + "/name":
					myService.Name = resp.Node.Nodes[i].Nodes[j].Value
				case lowerkey + "/description":
					myService.Description = resp.Node.Nodes[i].Nodes[j].Value
				case lowerkey + "/bindable":
					myService.Bindable, _ = strconv.ParseBool(resp.Node.Nodes[i].Nodes[j].Value)
				case lowerkey + "/tags":
					myService.Tags = strings.Split(resp.Node.Nodes[i].Nodes[j].Value, ",")
				case lowerkey + "/planupdatable":
					myService.PlanUpdatable, _ = strconv.ParseBool(resp.Node.Nodes[i].Nodes[j].Value)
				case lowerkey + "/metadata":
					json.Unmarshal([]byte(resp.Node.Nodes[i].Nodes[j].Value), &myService.Metadata)
				}
			} else if strings.HasSuffix(strings.ToLower(resp.Node.Nodes[i].Nodes[j].Key), "plan") {

				myPlans := []Plan{}
				for k := 0; k < len(resp.Node.Nodes[i].Nodes[j].Nodes); k++ {
					log.Logger.Debug("Start to Parse Plan " + resp.Node.Nodes[i].Nodes[j].Nodes[k].Key)
					myPlan := Plan{}
					myPlan.Id = strings.Split(resp.Node.Nodes[i].Nodes[j].Nodes[k].Key, "/")[len(strings.Split(resp.Node.Nodes[i].Nodes[j].Nodes[k].Key, "/"))-1]
					for n := 0; n < len(resp.Node.Nodes[i].Nodes[j].Nodes[k].Nodes); n++ {
						lowernodekey := strings.ToLower(resp.Node.Nodes[i].Nodes[j].Nodes[k].Key)
						switch strings.ToLower(resp.Node.Nodes[i].Nodes[j].Nodes[k].Nodes[n].Key) {
						case lowernodekey + "/name":
							myPlan.Name = resp.Node.Nodes[i].Nodes[j].Nodes[k].Nodes[n].Value
						case lowernodekey + "/description":
							myPlan.Description = resp.Node.Nodes[i].Nodes[j].Nodes[k].Nodes[n].Value
						case lowernodekey + "/free":
							myPlanfree, _ := strconv.ParseBool(resp.Node.Nodes[i].Nodes[j].Nodes[k].Nodes[n].Value)
							myPlan.Free = myPlanfree
						case lowernodekey + "/metadata":
							json.Unmarshal([]byte(resp.Node.Nodes[i].Nodes[j].Nodes[k].Nodes[n].Value), &myPlan.Metadata)
						}
					}
					myPlans = append(myPlans, myPlan)
				}
				myService.Plans = myPlans

			}
		}
		catalogRsp.Services = append(catalogRsp.Services, myService)
	}
	c.JSON(200, catalogRsp)
	return

}

func Provision(c *gin.Context) {

}


func PollingPlan(c *gin.Context) {

}



func ProvisionService(c *gin.Context) {
	sId := c.Param("service_id")
	key := KEY + "/" + sId
	rBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest,err)
		return
	}
	defer c.Request.Body.Close()
	var pservice CatalogResponse
	err = json.Unmarshal(rBody,&pservice)
	if err != nil{
		c.JSON(http.StatusBadRequest, err)
		return
	}
	etcdC := etcdclient.GetEtcdApi()
	req := &client.Response{}
	for i,v := range pservice.Services{
		tagName ,value := getTag(&v,i)
		key += "/" + tagName
		req,err = etcdC.Update(context.Background(),key,value)
		if err != nil{
			log.Logger.Error("Can not ProvisionService service from etcd", err)
			errinfo := ErrorResponse{}
			errinfo.Error = err.Error()
			errinfo.Description = "can not updata service from etcd"
			c.JSON(http.StatusNotImplemented, errinfo)
			return
		}
	}
	c.JSON(http.StatusOK,req.Node)
	return
}

func ProvisionPlan(c *gin.Context) {
	sId := c.Param("service_id")
	pId := c.Param("plan_id")
	key := KEY + "/" + sId + "/plan" + pId
	rBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "metrics": err})
		return
	}
	defer c.Request.Body.Close()
	var pservice CatalogResponse
	err = json.Unmarshal(rBody,&pservice)
	if err != nil{
		c.JSON(http.StatusBadRequest, err)
		return
	}
	etcdC := etcdclient.GetEtcdApi()
	req := &client.Response{}
	for i,v := range pservice.Services{
		tagName,value := getTag(&v,i)
		key += "/" + tagName
		req,err = etcdC.Update(context.Background(),key,value)
		if err != nil{
			log.Logger.Error("Can not ProvisionService service from etcd", err)
			errinfo := ErrorResponse{}
			errinfo.Error = err.Error()
			errinfo.Description = "can not updata service from etcd"
			c.JSON(http.StatusNotImplemented, errinfo)
			return
		}
	}
	c.JSON(http.StatusOK,req.Node)
	return
}

func DeprovisionService(c *gin.Context) {
	sId := c.Param("service_id")
	etcdC := etcdclient.GetEtcdApi()
	key := KEY + "/" + sId
	req,err := etcdC.Delete(context.Background(),key,&client.DeleteOptions{})
	if err != nil{
		log.Logger.Error("Can not DeprovisionService service from etcd", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "can not delete service from etcd"
		c.JSON(http.StatusNotImplemented, errinfo)
		return
	}
	c.JSON(http.StatusOK,req.Node)
	return
}

func DeprovisionPlan(c *gin.Context) {
	sId := c.Param("service_id")
	pId := c.Param("plan_id")
	etcdC := etcdclient.GetEtcdApi()
	key := KEY + "/" + sId + "/plan" + pId
	req,err := etcdC.Delete(context.Background(),key,&client.DeleteOptions{})
	if err != nil{
		log.Logger.Error("Can not DeprovisionPlan plan from etcd", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "can not delete plan from etcd"
		c.JSON(http.StatusNotImplemented, errinfo)
		return
	}
	c.JSON(http.StatusOK,req.Node)
	return
}

func getTag(u interface{},index int)(tag string,value string){
	t := reflect.TypeOf(u)
	v := reflect.ValueOf(u)
	field := t.Elem().Field(index)
	vName := v.Elem().FieldByName(field.Name)
	tag = field.Tag.Get("json")
	value = fmt.Sprintf("%v", vName.Interface())
	return
}