
### servicebroker-plan-api和版本

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

### 镜像运行命令

**docker run -p <宿主端口号>:<程序监听端口号> -e ETCDENDPOINT=<ETCD服务入口> -e ETCDUSER=<ETCD用户名> -e ETCDPASSWORD=<ETCD密码> -e PLANPORT=<程序监听端口号> <镜像名称>:<版本Tag>**

事例：docker run -p 8000:10000 -e ETCDENDPOINT="http://192.168.1.114:2379" -e ETCDUSER="root" -e ETCDPASSWORD="111111" -e PLANPORT="10000" mypalnapi:v1

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
curl -i -X GET http://asiainfoLDP:2016asia@127.0.0.1:10000/seapi/catalog
```
#####返回值
Http Code   | JSON
----------- | -------------
200         | 服务套餐列表
500         | etcd error

#### GET /seapi/services/{service_id}
获取某个服务信息。

Path参数
* `service_id`: 服务唯一标识。

curl样例：
```
curl -i -X GET http://asiainfoLDP:2016asia@127.0.0.1:10000/seapi/services/df17b082-a5a3-47d2-a42a-45c8dd285c70
```
#####返回值
Http Code   | JSON
----------- | -------------
200         | 服务信息
500         | etcd error

#### GET /seapi/services/{service_id}/plans/{plan_id}
获取某个服务下套餐信息。

Path参数
* `service_id`: 服务唯一标识。
* `plan_id`: 套餐唯一标识。

curl样例：
```
curl -i -X GET http://asiainfoLDP:2016asia@127.0.0.1:10000/seapi/services/df17b082-a5a3-47d2-a42a-45c8dd285c70/plans/0DA6918E-9D72-4A54-A75C-E0B7F9647300
```
#####返回值
Http Code   | JSON
----------- | -------------
200         | 套餐信息
500         | etcd error

#### GET /seapi/services/{service_id}/plans
获取某个服务下所有套餐信息列表。

Path参数
* `service_id`: 服务唯一标识。

curl样例：
```
curl -i -X GET http://asiainfoLDP:2016asia@127.0.0.1:10000/seapi/services/df17b082-a5a3-47d2-a42a-45c8dd285c70/plans
```
#####返回值
Http Code   | JSON
----------- | -------------
200         | 套餐列表信息
500         | etcd error

#### POST /seapi/services/{service_name}
创建一个服务。

Path参数
* `service_name`: 服务名称。

curl样例：
```
curl -i -X PUT http://asiainfoLDP:2016asia@127.0.0.1:10000/seapi/services/test_001 -d '{
 "description":"service test instance",
 "bindable":true,
 "tags": ["service","test"],
 "plan_updateable":true,
 "metadata": {"bullets":["20 GB of Disk","20 connections"],"displayName":"Shared and Free"}
}'  -H "Content-Type: application/json"
```
#####返回值
Http Code   | JSON
----------- | -------------
200         | 套餐列表信息
400         | 参数不规范错误
500         | server error

#### POST /seapi/services/{service_id}/plans/{plan_name}
在某个服务下创建一个套餐。

Path参数
* `service_id`: 服务唯一标识。
* `plan_name`: 套餐名称。

curl样例：
```
curl -i -X POST http://asiainfoLDP:2016asia@127.0.0.1:10000/seapi/services/df17b082-a5a3-47d2-a42a-45c8dd285c70/plans/standalone -d '{
 "description":"plan pvc",
 "free":true,
 "metadata": {"bullets":["20 GB of Disk","20 connections"],"displayName":"Shared and Free" }
}'  -H "Content-Type: application/json"
```
#####返回值
Http Code   | JSON
----------- | -------------
200         | 套餐列表信息
400         | 参数不规范错误
500         | server error

#### PUT /seapi/services/:service_id", planapi.UpdataService)

#### PUT /seapi/services/:service_id/plans/:plan_id", planapi.UpdataPlan)

#### DELETE /seapi/services/:service_id", planapi.DeprovisionService)

#### DELETE /seapi/services/:service_id/plans/:plan_id", planapi.DeprovisionPlan)
