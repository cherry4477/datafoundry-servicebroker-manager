package tools

import (
	log "github.com/asiainfoldp/datafoundry-servicebroker-manager/log"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"time"
)

type EtcdClient struct {
}

var etcdEndPoint string
var etcdUser string
var etcdPassword string
var etcdapi client.KeysAPI

func init() {
	etcdEndPoint = Getenv("ETCDENDPOINT")
	etcdUser = Getenv("ETCDUSER")
	etcdPassword = Getenv("ETCDPASSWORD")
	cfg := client.Config{
		Endpoints:               []string{etcdEndPoint},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second * 5,
		Username:                etcdUser,
		Password:                etcdPassword,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Logger.Error("Can not init ectd client", err)
	}
	etcdapi = client.NewKeysAPI(c)
}

func (ec *EtcdClient) GetEtcdApi() client.KeysAPI {
	return etcdapi
}

func (ec *EtcdClient) Etcdget(key string) (*client.Response, error) {
	n := 5

RETRY:
	resp, err := etcdapi.Get(context.Background(), key, nil)
	if err != nil {
		log.Logger.Error("Can not get "+key+" from etcd", err)
		n--
		if n > 0 {
			goto RETRY
		}

		return nil, err
	} else {
		log.Logger.Debug("Successful get " + key + " from etcd. value is " + resp.Node.Value)
		return resp, nil
	}
}

func (ec *EtcdClient) Etcdset(key string, value string) (*client.Response, error) {
	n := 5

RETRY:
	resp, err := etcdapi.Set(context.Background(), key, value, nil)
	if err != nil {
		log.Logger.Error("Can not set "+key+" from etcd", err)
		n--
		if n > 0 {
			goto RETRY
		}

		return nil, err
	} else {
		log.Logger.Debug("Successful set " + key + " from etcd. value is " + value)
		return resp, nil
	}
}
