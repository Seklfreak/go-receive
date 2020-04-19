package dpd_de

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Seklfreak/go-receive/companies"
)

var countryInScanLocationRegexp = regexp.MustCompile(`\(([A-Z]{2})\)`)

type Company struct {
}

func (c *Company) Track(httpClient *http.Client, id string) (*companies.Status, error) {
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return nil, fmt.Errorf("failure loading location \"Europe/Berlin\": %w", err)
	}

	resp, err := httpClient.Get(endpoint(id))
	if err != nil {
		return nil, fmt.Errorf("failure making request: %w", err)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failure reading request response: %w", err)
	}

	var respObject trackingResponse
	err = json.Unmarshal(respData, &respObject)
	if err != nil {
		return nil, fmt.Errorf("failure unmarshalling request response: %w", err)
	}

	var history []companies.HistoryItem
	for _, scan := range respObject.ParcellifecycleResponse.ParcelLifeCycleData.ScanInfo.Scan {
		at, err := time.ParseInLocation("2006-01-02T15:04:05", scan.Date, location)
		if err != nil {
			return nil, fmt.Errorf("failure parsing time \"%s\": %w", scan.Date, err)
		}
		at = at.UTC()

		var country string
		countryMatches := countryInScanLocationRegexp.FindStringSubmatch(scan.ScanData.Location)
		if len(countryMatches) >= 2 {
			country = countryMatches[1]
		}

		history = append(history, companies.HistoryItem{
			At:                 at,
			Location:           scan.ScanData.Location,
			LocationCountryISO: country,
			Message:            strings.Join(scan.ScanDescription.Content, ", "),
		})
	}

	status := companies.Status{
		ID:                 id,
		ReceiverCountryISO: respObject.ParcellifecycleResponse.ParcelLifeCycleData.ShipmentInfo.ReceiverCountryIsoCode,
		History:            history,
	}

	return &status, nil
}

func endpoint(id string) string {
	return "https://tracking.dpd.de/rest/plc/en_US/" + id
}

type trackingResponse struct {
	ParcellifecycleResponse struct {
		ParcelLifeCycleData struct {
			ShipmentInfo struct {
				ReceiverCountryIsoCode string `json:"receiverCountryIsoCode"`
			} `json:"shipmentInfo"`
			ScanInfo struct {
				Scan []struct {
					Date     string `json:"date"`
					ScanData struct {
						Location string `json:"location"`
					} `json:"scanData"`
					ScanDescription struct {
						Content []string `json:"content"`
					} `json:"scanDescription"`
				} `json:"scan"`
			} `json:"scanInfo"`
		} `json:"parcelLifeCycleData"`
	} `json:"parcellifecycleResponse"`
}
