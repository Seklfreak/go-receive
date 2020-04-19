package go_receive

import (
	"net/http"

	"github.com/Seklfreak/go-receive/companies"
)

type Company interface {
	Track(httpClient *http.Client, id string) (*companies.Status, error)
}
