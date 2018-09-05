package planapi

type Service struct {
	Name             string           `json:"name"`
	Id               string           `json:"id"`
	Description      string           `json:"description"`
	Tags             []string         `json:"tags,omitempty"`
	Requires         []string         `json:"requires,omitempty"`
	Bindable         bool             `json:"bindable"`
	Metadata         interface{}      `json:"metadata,omitempty"`
	Dashboard_client *DashboardClient `json:"dashboard_client,omitempty"`
	PlanUpdatable    bool             `json:"plan_updateable,omitempty"`
	Plans            []Plan           `json:"plans"`
}

type Plan struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Metadata    interface{} `json:"metadata,omitempty"`
	Free        bool        `json:"free"`
	Bindable    bool        `json:"bindable,omitempty"`
	Schemas     *Schemas    `json:"schemas,omitempty"`
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
