package api

// CalculateRequest is the request body for calculation.
type CalculateRequest struct {
	Amount int `json:"amount" example:"251"`
}

// PackSizesRequest is the request body for replacing configured pack sizes.
type PackSizesRequest struct {
	PackSizes []int `json:"pack_sizes" example:"250,500,1000,2000,5000"`
}

// PackSizesResponse is returned by pack size read/update endpoints.
type PackSizesResponse struct {
	PackSizes []int `json:"pack_sizes"`
}

// HealthResponse is the response model for service health checks.
type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}
