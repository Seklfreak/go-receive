package go_receive

import (
	"github.com/Seklfreak/go-receive/companies"
)

func (c *Client) Track(company Company, id string) (*companies.Status, error) {
	return company.Track(c.httpClient, id)
}
