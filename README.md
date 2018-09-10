# servicebroker-plan-api

## servicebroker-plan-api和版本

本程序为ServiceBroker的服务及套餐的增删改查API。
版本为v1。


### 需要的环境变量

ETCD服务入口:
* ETCDENDPOINT

ETCD用户名:
* ETCDUSER

ETCD密码:
* ETCDPASSWORD

服务监听端口:
* PLANPORT（默认为10000）

API身份认证用户名:
* SEAPIUSER

API身份认证密码:
* SEAPIPASSWORD

### 镜像运行命令

**docker run -p <宿主端口号>:<程序监听端口号> -e ETCDENDPOINT=<ETCD服务入口> -e ETCDUSER=<ETCD用户名> -e ETCDPASSWORD=<ETCD密码> -e SEAPIUSER=<API身份认证用户名> -e SEAPIPASSWORD=<API身份认证密码> -e PLANPORT=<程序监听端口号> <镜像名称>:<版本Tag>**

事例：docker run -p 8000:10000 -e ETCDENDPOINT="http://192.168.1.114:2379" -e ETCDUSER="root" -e ETCDPASSWORD="111111" -e SEAPIUSER=<API身份认证用户名> -e SEAPIPASSWORD=<API身份认证密码> -e PLANPORT="10000" mypalnapi:v1

### 生成镜像命令

工程根目录下：

*先输入make

*docker build -t mypalnapi:v1 .

*如果删除Makefile产生的多余目录及文件输入make clean

### API接口

#### GET /seapi/catalog
获取全部服务套餐列表。

curl样例：
```
curl -i -X GET http://$SEAPIUSER:$SEAPIPASSWORD@127.0.0.1:10000/seapi/catalog
```

返回结果样例：
```
{
	"services": [{
		"name": "Anaconda",
		"id": "dfc126e9-181a-4d13-a367-f84edfe617ed",
		"description": "Anaconda on Openshift",
		"tags": ["Anaconda", "openshift"],
		"bindable": true,
		"metadata": {
			"displayName": "Anaconda",
			"documentationUrl": "https://docs.anaconda.com/",
			"imageUrl": "pub/assets/Anaconda.png",
			"longDescription": "Anaconda in the cloud.",
			"providerDisplayName": "Asiainfo",
			"supportUrl": "https://www.anaconda.com/support/"
		},
		"plan_updateable": false,
		"plans": [{
			"id": "521a4a06-175a-43e6-b1bc-d9c684f76a0d",
			"name": "standalone",
			"description": "Anaconda on Openshift",
			"metadata": {
				"bullets": ["20 GB of Disk", "20 connections"],
				"displayName": "Shared and Free"
			},
			"free": true
		}, {
			"id": "0DA6918E-9D72-4A54-A75C-E0B7F9647300",
			"name": "test",
			"description": "Anaconda on Openshift",
			"metadata": {
				"bullets": ["20 GB of Disk", "20 connections"],
				"displayName": "Shared and Free"
			},
			"free": false
		}]
	}, {
		"name": "wu004",
		"id": "df17b082-a5a3-47d2-a42a-45c8dd285c70",
		"description": "service test instance",
		"tags": ["service2", "test2"],
		"bindable": false,
		"metadata": {
			"bullets": ["20 GB of Disk", "20 connections"],
			"displayName": "Shared and Free"
		},
		"plan_updateable": false,
		"plans": [{
			"id": "5840b24a-ffbe-4835-8f03-9ad793a3eeec",
			"name": "standalone",
			"description": "standalone",
			"metadata": {
				"bullets": ["22 GB of Disk", "23 connections"],
				"displayName": "Shared and Free"
			},
			"free": false
		}]
	}, {
		"name": "wu003",
		"id": "a5f614c9-04f1-4edd-8bad-6e3dc9a63aaf",
		"description": "service test instance",
		"tags": ["service1", "test1"],
		"bindable": false,
		"metadata": {
			"bullets": ["20 GB of Disk", "20 connections"],
			"displayName": "Shared and Free"
		},
		"plan_updateable": false,
		"plans": []
	}]
}
```

#### GET /seapi/services/{service_id}
获取某个服务信息。

Path参数
* `service_id`: 服务唯一标识。

curl样例：
```
curl -i -X GET http://$SEAPIUSER:$SEAPIPASSWORD@127.0.0.1:10000/seapi/services/dfc126e9-181a-4d13-a367-f84edfe617ed
```

返回结果样例：
```
{
	"name": "Anaconda",
	"id": "dfc126e9-181a-4d13-a367-f84edfe617ed",
	"description": "Anaconda on Openshift",
	"tags": ["Anaconda", "openshift"],
	"bindable": true,
	"metadata": {
		"displayName": "Anaconda",
		"documentationUrl": "https://docs.anaconda.com/",
		"imageUrl": "pub/assets/Anaconda.png",
		"longDescription": "Anaconda in the cloud.",
		"providerDisplayName": "Asiainfo",
		"supportUrl": "https://www.anaconda.com/support/"
	},
	"plan_updateable": false,
	"plans": null
}
```


#### GET /seapi/services/{service_id}/plans/{plan_id}
获取某个服务下套餐信息。

Path参数
* `service_id`: 服务唯一标识。
* `plan_id`: 套餐唯一标识。

curl样例：
```
curl -i -X GET http://$SEAPIUSER:$SEAPIPASSWORD@127.0.0.1:10000/seapi/services/dfc126e9-181a-4d13-a367-f84edfe617ed/plans/0521a4a06-175a-43e6-b1bc-d9c684f76a0d
```

返回结果样例：
```
{
	"id": "521a4a06-175a-43e6-b1bc-d9c684f76a0d",
	"name": "standalone",
	"description": "Anaconda on Openshift",
	"metadata": {
		"bullets": ["20 GB of Disk", "20 connections"],
		"displayName": "Shared and Free"
	},
	"free": true
}
```

#### GET /seapi/services/{service_id}/plans
获取某个服务下所有套餐信息列表。

Path参数
* `service_id`: 服务唯一标识。

curl样例：
```
curl -i -X GET http://$SEAPIUSER:$SEAPIPASSWORD@127.0.0.1:10000/seapi/services/dfc126e9-181a-4d13-a367-f84edfe617ed/plans
```

返回结果样例：
```
{
	"plans": [{
		"id": "521a4a06-175a-43e6-b1bc-d9c684f76a0d",
		"name": "standalone",
		"description": "Anaconda on Openshift",
		"metadata": {
			"bullets": ["20 GB of Disk", "20 connections"],
			"displayName": "Shared and Free"
		},
		"free": true
	}, {
		"id": "0DA6918E-9D72-4A54-A75C-E0B7F9647300",
		"name": "test",
		"description": "Anaconda on Openshift",
		"metadata": {
			"bullets": ["20 GB of Disk", "20 connections"],
			"displayName": "Shared and Free"
		},
		"free": false
	}]
}
```

#### POST /seapi/services/{service_name}
创建一个服务。

Path参数
* `service_name`: 服务名称。

curl样例：
```
curl -i -X POST http://$SEAPIUSER:$SEAPIPASSWORD@127.0.0.1:10000/seapi/services/test_001 -d '{
 "description":"service test instance",
 "bindable":true,
 "tags": ["service1","test1"],
 "plan_updateable":true,
 "metadata": {"bullets":["20 GB of Disk","20 connections"],"displayName":"Shared and Free"}
}'  -H "Content-Type: application/json"
```

请求体：

| 参数   | 类型| 是否必传|
|--- | ---|---|
|description | string|是|
|tags | []string|是|
|bindable | bool|是|
|metadata | json object|是|
|plan_updateable | bool|是|

返回结果样例：
```
{
	"name": "test_001",
	"id": "3e23e5f4-c516-4891-b9d0-2bbc897a13d7",
	"description": "service test instance",
	"tags": ["service1", "test1"],
	"bindable": true,
	"metadata": {
		"bullets": ["20 GB of Disk", "20 connections"],
		"displayName": "Shared and Free"
	},
	"plan_updateable": true,
	"plans": null
}
```

#### POST /seapi/services/{service_id}/plans/{plan_name}
在某个服务下创建一个套餐。

Path参数
* `service_id`: 服务唯一标识。
* `plan_name`: 套餐名称。

curl样例：
```
curl -i -X POST http://$SEAPIUSER:$SEAPIPASSWORD@127.0.0.1:10000/seapi/services/3e23e5f4-c516-4891-b9d0-2bbc897a13d7/plans/standalone -d '{
 "description":"plan pvc",
 "free":true,
 "metadata": {"bullets":["20 GB of Disk","20 connections"],"displayName":"Shared and Free" }
}'  -H "Content-Type: application/json"
```

请求体：

| 参数   | 类型| 是否必传|
|--- | ---|---|
|description | string|是|
|metadata | json object|是|
|free | bool|是|

返回结果样例：
```
{
	"id": "543a2e39-a339-436d-b708-c914e1675c55",
	"name": "standalone",
	"description": "plan standalone",
	"metadata": {
		"bullets": ["20 GB of Disk", "20 connections"],
		"displayName": "Shared and Free"
	},
	"free": true
}
```

#### PUT /seapi/services/{service_id}
更新一个服务。

Path参数
* `service_id`: 服务ID。

**注：bindable、PlanUpdatable两个字段为bool类型，为必传字段**

curl样例：
```
curl -i -X PUT http://$SEAPIUSER:$SEAPIPASSWORD@127.0.0.1:10000/seapi/services/3e23e5f4-c516-4891-b9d0-2bbc897a13d7 -d '{
 "description":"test_2",
 "bindable":false,
 "plan_updateable":false,
 "metadata": {"bullets":["5 GB of Disk","5 connections"],"displayName":"Shared and Free"}
}'  -H "Content-Type: application/json"
```

请求体：

| 参数   | 类型| 是否必传|
|--- | ---|---|
|name | string| 否|
|description | string|否|
|tags | []string|否|
|bindable | bool|是|
|metadata | json object|否|
|plan_updateable | bool|是|

返回结果样例：
```
{
	"name": "",
	"id": "",
	"description": "test_2",
	"tags": null,
	"bindable": false,
	"metadata": {
		"bullets": ["5 GB of Disk", "5 connections"],
		"displayName": "Shared and Free"
	},
	"plan_updateable": false,
	"plans": null
}
```

#### PUT /seapi/services/{service_id}/plans/{plan_id}
更新服务下的套餐

Path参数
* `service_id`: 服务ID。
* `plan_id`: 套餐ID。

**注：free字段为bool类型，为必传字段**

curl样例：
```
curl -i -X PUT http://$SEAPIUSER:$SEAPIPASSWORD@127.0.0.1:10000/seapi/services/3e23e5f4-c516-4891-b9d0-2bbc897a13d7/plans/543a2e39-a339-436d-b708-c914e1675c55 -d '{
 "description":"standalone2",
 "free":false,
 "name":"standalone2",
 "metadata": {"bullets":["15 GB of Disk","15 connections"],"displayName":"Shared and Free"}
}'  -H "Content-Type: application/json"
```
请求体：

| 参数   | 类型| 是否必传|
|--- | ---|---|
|name | string| 否|
|description | string|否|
|metadata | json object|否|
|free | bool|是|

返回结果样例：
```
{
	"id": "",
	"name": "standalone2",
	"description": "standalone2",
	"metadata": {
		"bullets": ["15 GB of Disk", "15 connections"],
		"displayName": "Shared and Free"
	},
	"free": false
}
```


#### DELETE /seapi/services/{service_id}
删除服务

Path参数
* `service_id`: 服务ID。

curl样例：
```
curl -i -X DELETE http://$SEAPIUSER:$SEAPIPASSWORD@127.0.0.1:10000/seapi/services/3e23e5f4-c516-4891-b9d0-2bbc897a13d7
```

返回结果样例：
```
{
	"key": "/servicebroker/openshift/catalog/3e23e5f4-c516-4891-b9d0-2bbc897a13d7",
	"dir": true,
	"value": "",
	"nodes": null,
	"createdIndex": 386,
	"modifiedIndex": 409
}
```

#### DELETE /seapi/services/{service_id}/plans/{plan_id}
删除服务下的套餐

Path参数
* `service_id`: 服务ID。
* `plan_id`: 套餐ID。

curl样例：
```
curl -i -X DELETE http://$SEAPIUSER:$SEAPIPASSWORD@127.0.0.1:10000/seapi/services/3e23e5f4-c516-4891-b9d0-2bbc897a13d7/plans/543a2e39-a339-436d-b708-c914e1675c55
```

返回结果样例：
```
{
	"key": "/servicebroker/openshift/catalog/3e23e5f4-c516-4891-b9d0-2bbc897a13d7/plan/543a2e39-a339-436d-b708-c914e1675c55",
	"dir": true,
	"value": "",
	"nodes": null,
	"createdIndex": 393,
	"modifiedIndex": 408
}
```

### 错误码
| Http Code   | JSON|
| --- | ---|
| 200         | JSON信息|
| 400         | 参数错误|
| 401         | API身份认证错误|
| 409         | 服务或套餐名称冲突|
| 500         | server error|
