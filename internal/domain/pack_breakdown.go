package domain

// PackBreakdown describes the number of packs to ship for a given pack size.
type PackBreakdown struct {
	Size  int `json:"size"`
	Count int `json:"count"`
}
