package goswiftpos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultSendTimeout = time.Second * 30
	baseURL            = "https://webstores.swiftpos.com.au:4443"
	authorisationURL   = "SwiftApi/api/Authorisation"
	salesURL           = "SwiftApi/api/Sale"
)

// Swiftpos The main struct of this package
type Swiftpos struct {
	LocationID string
	ClerkID    string
	Password   string
	APIKey     string
	Timeout    time.Duration
}

// NewClient will create a SwiftPOS client with default values
func NewClient(locationID string, clerkID string, password string) *Swiftpos {
	return &Swiftpos{
		LocationID: locationID,
		ClerkID:    clerkID,
		Password:   password,
		Timeout:    defaultSendTimeout,
	}
}

// Authorisation will get and assign the APIKey
func (v *Swiftpos) Authorisation() error {
	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return fmt.Errorf("Failed to build Swiftpos authorisation: %v", err)
	}

	u.Path = authorisationURL
	urlStr := fmt.Sprintf("%v", u)

	r, err := http.NewRequest("GET", urlStr, nil)

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")

	data := url.Values{}
	data.Add("locationId", v.LocationID)
	data.Add("userId", v.ClerkID)
	data.Add("password", v.Password)
	r.URL.RawQuery = data.Encode()

	res, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("Failed to call Swiftpos authorisation: %v", err)
	}

	if res.StatusCode == 200 {
		rawResBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Failed to read Swiftpos authorisation: %v", err)
		}

		var resp Authorisation
		err = json.Unmarshal(rawResBody, &resp)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal Swiftpos authorisation: %v", err)
		}
		v.APIKey = resp.APIKey

		fmt.Println(v.APIKey)
		return nil

	}
	return fmt.Errorf("Failed to get Swiftpos invoices: %s", res.Status)
}

// GetSales will get the sales for a site
func (v *Swiftpos) GetSales(saleID string) ([]Sales, error) {
	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to build Swiftpos sales %v", err)
	}

	u.Path = salesURL
	urlStr := fmt.Sprintf("%v", u)

	r, err := http.NewRequest("GET", urlStr, nil)

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("ApiKey", v.APIKey)

	if saleID != "" {
		data := url.Values{}
		data.Add("saleId", saleID)
		// data.Add("maxRecords", "1")
		r.URL.RawQuery = data.Encode()
	}

	res, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("Failed to call Swiftpos sales: %v", err)
	}

	if res.StatusCode == 200 {
		rawResBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Failed to read Swiftpos sales: %v", err)
		}

		//test
		// fmt.Println("rawResBody", string(rawResBody))

		var resp []Sales
		err = json.Unmarshal(rawResBody, &resp)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal SwiftPOS sales: %v", err)
		}
		return resp, nil

	}
	return nil, fmt.Errorf("Failed to get Swiftpos sales: %s", res.Status)
}

func checkRedirectFunc(req *http.Request, via []*http.Request) error {
	if req.Header.Get("Authorization") == "" {
		req.Header.Add("Authorization", via[0].Header.Get("Authorization"))
	}
	return nil
}
