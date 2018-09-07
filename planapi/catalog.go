package planapi

type Service struct {
	Name             string           `json:"name" bson:"name"`
	Id               string           `json:"id" bson:"id"`
	Description      string           `json:"description" bson:"description"`
	Tags             []string         `json:"tags" bson:"tags"`
	Requires         []string         `json:"requires,omitempty" bson:"requires"`
	Bindable         bool             `json:"bindable" bson:"bindable"`
	Metadata         interface{}      `json:"metadata" bson:"metadata"`
	Dashboard_client *DashboardClient `json:"dashboard_client,omitempty" bson:"dashboardclient"`
	PlanUpdatable    bool             `json:"plan_updateable" bson:"planupdatable"`
	Plans            []Plan           `json:"plans" bson:"plans"`
}

type Plan struct {
	Id          string      `json:"id" bson:"id"`
	Name        string      `json:"name" bson:"name"`
	Description string      `json:"description" bson:"description"`
	Metadata    interface{} `json:"metadata,omitempty" bson:"metadata"`
	Free        bool        `json:"free" bson:"free"`
	//Bindable    bool        `json:"bindable,omitempty" bson:"bindable"`
	Schemas *Schemas `json:"schemas,omitempty" bson:"schemas"`
}

type Schemas struct {
	Service_instance *ServiceInstanceSchema `json:"service_instance,omitempty"`
	Service_binding  *ServiceBindingSchema  `json:"service_binding,omitempty"`
}

type ServiceInstanceSchema struct {
	Create *InputParametersSchema `json:"create,omitempty"`
	Update *InputParametersSchema `json:"update,omitempty"`
}

type ServiceBindingSchema struct {
	Create *InputParametersSchema `json:"create,omitempty"`
}

type InputParametersSchema struct {
	Parameters interface{} `json:"parameters,omitempty"`
}

type DashboardClient struct {
	Id           string `json:"id"`
	Secret       string `json:"secret"`
	Redirect_uri string `json:"redirect_uri"`
}
