package planapi

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}

type CatalogResponse struct {
	Services []Service `json:"services"`
}

type PlansResponse struct {
	Plans []Plan `json:"plans"`
}
