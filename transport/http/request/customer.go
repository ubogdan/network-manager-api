package request

import (
	"github.com/ubogdan/network-manager-api/model"
)

type Customer struct {
	Name         string `json:"name"`
	Country      string `json:"country"`
	City         string `json:"city"`
	Organization string `json:"organization"`
}

func (c *Customer) ToModel() model.Customer {
	return model.Customer{
		Name:         c.Name,
		Country:      c.Country,
		City:         c.City,
		Organization: c.Organization,
	}
}
