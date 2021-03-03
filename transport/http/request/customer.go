package request

import (
	"github.com/ubogdan/network-manager-api/model"
)

// Customer request DTO.
type Customer struct {
	Name         string `json:"name"`
	Country      string `json:"country"`
	City         string `json:"city"`
	Organization string `json:"organization"`
}

// ToModel convert DTO to customer model.
func (c *Customer) ToModel() model.Customer {
	return model.Customer{
		Name:         c.Name,
		Country:      c.Country,
		City:         c.City,
		Organization: c.Organization,
	}
}
