package response

// Customer response DTO.
type Customer struct {
	Name         string `json:"name"`
	Country      string `json:"country,omitempty"`
	City         string `json:"city,omitempty"`
	Organization string `json:"organization,omitempty"`
}
