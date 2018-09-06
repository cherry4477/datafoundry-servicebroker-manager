package planapi

import (
	"context"
	"github.com/asiainfoLDP/servicebroker-plan-api/tools"
	"github.com/coreos/etcd/client"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	//"golang.org/x/net/context"
	"encoding/json"
	"github.com/asiainfoLDP/servicebroker-plan-api/log"
	"github.com/pkg/errors"
	"strconv"
	"strings"
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

//获取服务名称
func findServiceNameInCatalog(service_id string) string {
	resp, err := etcdclient.Etcdget("/servicebroker/" + log.ServcieBrokerName + "/catalog/" + service_id + "/name")
	if err != nil {
		return ""
	}
	return resp.Node.Value
}

///seapi/services/:service_id
func PollingService(c *gin.Context) {
	service_id := c.Param("service_id")

	//获取service信息
	resp, err := etcdclient.GetEtcdApi().Get(context.Background(),
		"/servicebroker/"+log.ServcieBrokerName+"/catalog/"+service_id,
		&client.GetOptions{Recursive: true})
	if err != nil {
		log.Logger.Error("Can not get service information from etcd", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "can not get service information from etcd"
		c.JSON(500, errinfo)
		return
	} else {
		log.Logger.Debug("Successful get service information from etcd. NodeInfo is " + resp.Node.Key)
	}

	myService := Service{}
	for i := 0; i < len(resp.Node.Nodes); i++ {

		lowerkey := strings.ToLower(resp.Node.Nodes[i].Key)
		switch lowerkey {
		case "name":
			myService.Name = resp.Node.Nodes[i].Value
		case "description":
			myService.Description = resp.Node.Nodes[i].Value
		case "bindable":
			myService.Bindable, _ = strconv.ParseBool(resp.Node.Nodes[i].Value)
		case "tags":
			myService.Tags = strings.Split(resp.Node.Nodes[i].Value, ",")
		case "planupdatable":
			myService.PlanUpdatable, _ = strconv.ParseBool(resp.Node.Nodes[i].Value)
		case "metadata":
			json.Unmarshal([]byte(resp.Node.Nodes[i].Value), &myService.Metadata)
		}
	}

	c.JSON(200, myService)
}

///seapi/services/:service_id/plans/:plan_id
func PollingPlan(c *gin.Context) {

	service_id := c.Param("service_id")
	plan_id := c.Param("plan_id")

	service_name := findServiceNameInCatalog(service_id)

	if len(service_name) == 0 {
		log.Logger.Debug("findServiceNameInCatalog with service_id:" + service_id + " error")
		errinfo := ErrorResponse{}
		errinfo.Error = errors.New("findServiceNameInCatalog with service_id:" + service_id + " error").Error()
		errinfo.Description = "service_id:" + service_id + " is not correct."
		c.JSON(500, errinfo)
		return
	}

	//获取plan信息
	resp, err := etcdclient.GetEtcdApi().Get(context.Background(),
		"/servicebroker/"+log.ServcieBrokerName+"/catalog/"+service_id+"/plan/"+plan_id,
		&client.GetOptions{Recursive: true})
	if err != nil {
		log.Logger.Error("Can not get plan information in the service:"+service_name+" from etcd", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "Can not get plan information in the service:" + service_name + " from etcd"
		c.JSON(500, errinfo)
		return
	} else {
		log.Logger.Debug("Successful get plan information from etcd. NodeInfo is " + resp.Node.Key)
	}

	myPlan := Plan{}
	for i := 0; i < len(resp.Node.Nodes); i++ {

		lowerkey := strings.ToLower(resp.Node.Nodes[i].Key)
		switch lowerkey {
		case "name":
			myPlan.Name = resp.Node.Nodes[i].Value
		case "description":
			myPlan.Description = resp.Node.Nodes[i].Value
		case "metadata":
			json.Unmarshal([]byte(resp.Node.Nodes[i].Value), &myPlan.Metadata)
		case "free":
			myPlan.Free, _ = strconv.ParseBool(resp.Node.Nodes[i].Value)
		}
	}

	c.JSON(200, myPlan)
}

///seapi/services/:service_id/plans
func PollingPlans(c *gin.Context) {

	plansRsp := PlansResponse{}

	service_id := c.Param("service_id")

	service_name := findServiceNameInCatalog(service_id)

	if len(service_name) == 0 {
		log.Logger.Debug("findServiceNameInCatalog with service_id:" + service_id + " error")
		errinfo := ErrorResponse{}
		errinfo.Error = errors.New("findServiceNameInCatalog with service_id:" + service_id + " error").Error()
		errinfo.Description = "service_id:" + service_id + " is not correct."
		c.JSON(500, errinfo)
		return
	}
	//获取plans信息
	resp, err := etcdclient.GetEtcdApi().Get(context.Background(),
		"/servicebroker/"+log.ServcieBrokerName+"/catalog/"+service_id+"/plan",
		&client.GetOptions{Recursive: true})
	if err != nil {
		log.Logger.Error("Can not get plans information in the service:"+service_name+" from etcd", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "Can not get plans information in the service:" + service_name + " from etcd"
		c.JSON(500, errinfo)
		return
	} else {
		log.Logger.Debug("Successful get plans information from etcd. NodeInfo is " + resp.Node.Key)
	}

	for i := 0; i < len(resp.Node.Nodes); i++ {
		log.Logger.Debug("Start to Parse Plan " + resp.Node.Nodes[i].Key)

		myPlan := Plan{}
		myPlan.Id = strings.Split(resp.Node.Nodes[i].Key, "/")[len(strings.Split(resp.Node.Nodes[i].Key, "/"))-1]
		for j := 0; j < len(resp.Node.Nodes[i].Nodes); j++ {
			lowernodekey := strings.ToLower(resp.Node.Nodes[i].Nodes[j].Key)
			switch lowernodekey {
			case "name":
				myPlan.Name = resp.Node.Nodes[i].Nodes[j].Value
			case "description":
				myPlan.Description = resp.Node.Nodes[i].Nodes[j].Value
			case "free":
				myPlanfree, _ := strconv.ParseBool(resp.Node.Nodes[i].Nodes[j].Value)
				myPlan.Free = myPlanfree
			case "metadata":
				json.Unmarshal([]byte(resp.Node.Nodes[i].Nodes[j].Value), &myPlan.Metadata)
			}
		}
		plansRsp.Plans = append(plansRsp.Plans, myPlan)
	}
	c.JSON(200, plansRsp)
	return
}

///seapi/services/:service_id
func ProvisionService(c *gin.Context) {

}

///seapi/services/:service_id/plans/:plan_id
func ProvisionPlan(c *gin.Context) {

}

func Deprovision(c *gin.Context) {
	ins := c.Param("serviceinstance")
	etcdC := etcdclient.GetEtcdApi()
	req, err := etcdC.Delete(context.Background(), ins, &client.DeleteOptions{})
	if err != nil {
		log.Logger.Error("Can not Deprovision serviceinstace from etcd", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "can not delete serviceinstance from etcd"
		c.JSON(http.StatusNotImplemented, errinfo)
		return
	}
	c.JSON(http.StatusOK, req.Node)
	return
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
	req, err := etcdC.Update(context.Background(), ins, string(rBody))
	if err != nil {
		log.Logger.Error("Can not Deprovision serviceinstace from etcd", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "can not delete serviceinstance from etcd"
		c.JSON(http.StatusNotImplemented, errinfo)
		return
	}
	c.JSON(http.StatusOK, req.Node)
	return
}
