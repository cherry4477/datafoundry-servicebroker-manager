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
	"errors"
	"fmt"
	"github.com/asiainfoLDP/servicebroker-plan-api/log"
	"reflect"
	"strconv"
	"strings"
)

const (
	KEY = "/servicebroker/" + log.ServcieBrokerName + "/catalog"
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

		lowerkey := strings.ToLower(resp.Node.Key)
		switch strings.ToLower(resp.Node.Nodes[i].Key) {
		case lowerkey + "/name":
			myService.Name = resp.Node.Nodes[i].Value
		case lowerkey + "/description":
			myService.Description = resp.Node.Nodes[i].Value
		case lowerkey + "/bindable":
			myService.Bindable, _ = strconv.ParseBool(resp.Node.Nodes[i].Value)
		case lowerkey + "/tags":
			myService.Tags = strings.Split(resp.Node.Nodes[i].Value, ",")
		case lowerkey + "/planupdatable":
			myService.PlanUpdatable, _ = strconv.ParseBool(resp.Node.Nodes[i].Value)
		case lowerkey + "/metadata":
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

///seapi/services/:service_name
func ProvisionService(c *gin.Context) {
	service_name := c.Param("service_id")

	service_id := tools.Getuuid()

	_, err := etcdclient.GetEtcdApi().Set(context.Background(),
		"/servicebroker/"+log.ServcieBrokerName+"/catalog/"+service_id,
		"",
		&client.SetOptions{Dir: true})
	if err != nil {
		log.Logger.Error("etcdapi.Set service:"+service_name+" error", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "etcdapi.Set service:" + service_name + " error"
		c.JSON(500, errinfo)
		return
	} else {
		log.Logger.Debug("Successful create service:" + service_name + "in etcd.")
	}
	rBody, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		log.Logger.Error("Get provision service data error", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "get provision service data error"
		c.JSON(400, errinfo)
		return
	}
	rsp := Service{}
	err = json.Unmarshal(rBody, &rsp)
	if err != nil {
		log.Logger.Error("Parsing service data error", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "parsing service data error"
		c.JSON(500, errinfo)
		return
	}

	etcdclient.Etcdset("/servicebroker/"+log.ServcieBrokerName+"/catalog/"+service_id+"/name", service_name)

	etcdclient.Etcdset("/servicebroker/"+log.ServcieBrokerName+"/catalog/"+service_id+"/description", rsp.Description)

	etcdclient.Etcdset("/servicebroker/"+log.ServcieBrokerName+"/catalog/"+service_id+"/bindable", strconv.FormatBool(rsp.Bindable))

	etcdclient.Etcdset("/servicebroker/"+log.ServcieBrokerName+"/catalog/"+service_id+"/planupdatable", strconv.FormatBool(rsp.PlanUpdatable))

	etcdclient.Etcdset("/servicebroker/"+log.ServcieBrokerName+"/catalog/"+service_id+"/tags", tools.GetTagstring(rsp.Tags))

	tmpval, _ := json.Marshal(rsp.Metadata)
	etcdclient.Etcdset("/servicebroker/"+log.ServcieBrokerName+"/catalog/"+service_id+"/metadata", string(tmpval))

	rsp.Id = service_id
	rsp.Name = service_name
	c.JSON(200, rsp)
	return
}

///seapi/services/:service_id/plans/:plan_name
func ProvisionPlan(c *gin.Context) {
	service_id := c.Param("service_id")
	plan_name := c.Param("plan_id")
	plan_id := tools.Getuuid()
	_, err := etcdclient.GetEtcdApi().Set(context.Background(),
		"/servicebroker/"+log.ServcieBrokerName+"/catalog/"+service_id+"/plan",
		"",
		&client.SetOptions{Dir: true})
	if err != nil {
		log.Logger.Error("etcdapi.Set plan:"+plan_name+" error", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "etcdapi.Set plan:" + plan_name + " error"
		c.JSON(500, errinfo)
		return
	} else {
		log.Logger.Debug("Successful create plan:" + plan_name + "in etcd.")
	}

	rBody, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		log.Logger.Error("Get provision plan data error", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "get provision plan data error"
		c.JSON(400, errinfo)
		return
	}
	rsp := Plan{}
	err = json.Unmarshal(rBody, &rsp)
	if err != nil {
		log.Logger.Error("Parsing plan data error", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "parsing plan data error"
		c.JSON(500, errinfo)
		return
	}

	etcdclient.Etcdset("/servicebroker/"+log.ServcieBrokerName+"/catalog/"+service_id+"/plan/"+plan_id+"/name", plan_name)

	etcdclient.Etcdset("/servicebroker/"+log.ServcieBrokerName+"/catalog/"+service_id+"/plan/"+plan_id+"/description", rsp.Description)

	etcdclient.Etcdset("/servicebroker/"+log.ServcieBrokerName+"/catalog/"+service_id+"/plan/"+plan_id+"/free", strconv.FormatBool(rsp.Free))

	tmpval, _ := json.Marshal(rsp.Metadata)
	etcdclient.Etcdset("/servicebroker/"+log.ServcieBrokerName+"/catalog/"+service_id+"/plan/"+plan_id+"/metadata", string(tmpval))

	rsp.Id = plan_id
	rsp.Name = plan_name
	c.JSON(200, rsp)
	return

}

func UpdataService(c *gin.Context) {
	sId := c.Param("service_id")
	key := KEY + "/" + sId
	rBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	defer c.Request.Body.Close()
	var pservice CatalogResponse
	err = json.Unmarshal(rBody, &pservice)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	etcdC := etcdclient.GetEtcdApi()
	req := &client.Response{}
	for _,v := range pservice.Services{
		mValue := getTag(&v)
		for mk,mv := range mValue{
			key += "/" + mk
			req,err = etcdC.Update(context.Background(),key,mv)
			if err != nil{
				log.Logger.Error("Can not UpdataService service from etcd", err)
				errinfo := ErrorResponse{}
				errinfo.Error = err.Error()
				errinfo.Description = "can not updata service from etcd"
				c.JSON(http.StatusNotImplemented, errinfo)
				return
			}
		}
	}
	c.JSON(http.StatusOK, req.Node)
	return
}

func UpdataPlan(c *gin.Context) {
	sId := c.Param("service_id")
	pId := c.Param("plan_id")
	key := KEY + "/" + sId + "/plan" + pId
	rBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "metrics": err})
		return
	}
	defer c.Request.Body.Close()
	var pservice PlansResponse
	err = json.Unmarshal(rBody, &pservice)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	etcdC := etcdclient.GetEtcdApi()
	req := &client.Response{}
	for _,v := range pservice.Plans{
		mValue := getTag(&v)
		for mk,mv := range mValue{
			key += "/" + mk
			req,err = etcdC.Update(context.Background(),key,mv)
			if err != nil{
				log.Logger.Error("Can not UpdataPlan service from etcd", err)
				errinfo := ErrorResponse{}
				errinfo.Error = err.Error()
				errinfo.Description = "can not updata service from etcd"
				c.JSON(http.StatusNotImplemented, errinfo)
				return
			}
		}
	}
	c.JSON(http.StatusOK, req.Node)
	return
}

func DeprovisionService(c *gin.Context) {
	sId := c.Param("service_id")
	etcdC := etcdclient.GetEtcdApi()
	key := KEY + "/" + sId
	req, err := etcdC.Delete(context.Background(), key, &client.DeleteOptions{})
	if err != nil {
		log.Logger.Error("Can not DeprovisionService service from etcd", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "can not delete service from etcd"
		c.JSON(http.StatusNotImplemented, errinfo)
		return
	}
	c.JSON(http.StatusOK, req.Node)
	return
}
func DeprovisionPlan(c *gin.Context) {
	sId := c.Param("service_id")
	pId := c.Param("plan_id")
	etcdC := etcdclient.GetEtcdApi()
	key := KEY + "/" + sId + "/plan" + pId

	req, err := etcdC.Delete(context.Background(), key, &client.DeleteOptions{})
	if err != nil {
		log.Logger.Error("Can not DeprovisionPlan plan from etcd", err)
		errinfo := ErrorResponse{}
		errinfo.Error = err.Error()
		errinfo.Description = "can not delete plan from etcd"
		c.JSON(http.StatusNotImplemented, errinfo)
		return
	}
	c.JSON(http.StatusOK, req.Node)
	return
}

func getTag(u interface{})(value map[string]string){
	t := reflect.TypeOf(u)
	v := reflect.ValueOf(u)
	value = make(map[string]string)
	for i := 0; i < t.Elem().NumField();i++{
		field := t.Elem().Field(i)
		vName := v.Elem().FieldByName(field.Name)
		val := fmt.Sprintf("%v", vName.Interface())
		tag := field.Tag.Get("json")
		value[tag] = val
	}
	return
}
