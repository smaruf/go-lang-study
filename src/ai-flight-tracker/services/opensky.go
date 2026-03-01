package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// FlightState represents a single aircraft state vector from OpenSky.
type FlightState struct {
	ICAO24        string   `json:"icao24"`
	Callsign      string   `json:"callsign"`
	OriginCountry string   `json:"origin_country"`
	Lat           *float64 `json:"lat"`
	Lon           *float64 `json:"lon"`
	Altitude      *float64 `json:"altitude"`
	Velocity      *float64 `json:"velocity"`
	Heading       *float64 `json:"heading"`
	OnGround      bool     `json:"on_ground"`
	LastContact   *int64   `json:"last_contact"`
}

func parseState(state []interface{}) FlightState {
	get := func(i int) interface{} {
		if i < len(state) {
			return state[i]
		}
		return nil
	}
	getString := func(i int) string {
		v := get(i)
		if v == nil {
			return ""
		}
		s, _ := v.(string)
		return s
	}
	getFloat := func(i int) *float64 {
		v := get(i)
		if v == nil {
			return nil
		}
		f, ok := v.(float64)
		if !ok {
			return nil
		}
		return &f
	}
	getBool := func(i int) bool {
		v := get(i)
		if v == nil {
			return false
		}
		b, _ := v.(bool)
		return b
	}
	getInt64 := func(i int) *int64 {
		v := get(i)
		if v == nil {
			return nil
		}
		f, ok := v.(float64)
		if !ok {
			return nil
		}
		n := int64(f)
		return &n
	}

	// Trim spaces like Python does
	trimmed := strings.ReplaceAll(getString(1), " ", "")

	return FlightState{
		ICAO24:        getString(0),
		Callsign:      trimmed,
		OriginCountry: getString(2),
		LastContact:   getInt64(4),
		Lon:           getFloat(5),
		Lat:           getFloat(6),
		Altitude:      getFloat(7),
		OnGround:      getBool(8),
		Velocity:      getFloat(9),
		Heading:       getFloat(10),
	}
}

// OpenSkyClient is a client for the OpenSky Network API.
type OpenSkyClient struct {
	baseURL  string
	username string
	password string
	timeout  time.Duration
}

// NewOpenSkyClient creates a new OpenSky client.
func NewOpenSkyClient(baseURL, username, password string, timeoutSecs int) *OpenSkyClient {
	return &OpenSkyClient{
		baseURL:  baseURL,
		username: username,
		password: password,
		timeout:  time.Duration(timeoutSecs) * time.Second,
	}
}

// FetchFlights retrieves live flight states, optionally bounded by a geographic bbox.
func (c *OpenSkyClient) FetchFlights(lamin, lamax, lomin, lomax *float64) ([]FlightState, error) {
	url := fmt.Sprintf("%s/states/all", c.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	if lamin != nil {
		q.Set("lamin", fmt.Sprintf("%f", *lamin))
		q.Set("lamax", fmt.Sprintf("%f", *lamax))
		q.Set("lomin", fmt.Sprintf("%f", *lomin))
		q.Set("lomax", fmt.Sprintf("%f", *lomax))
	}
	req.URL.RawQuery = q.Encode()
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	client := &http.Client{Timeout: c.timeout}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("OpenSky fetch failed: %v", err)
		return []FlightState{}, nil
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		States [][]interface{} `json:"states"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return []FlightState{}, nil
	}
	var flights []FlightState
	for _, s := range data.States {
		fs := parseState(s)
		if fs.Lat == nil || fs.Lon == nil {
			continue
		}
		flights = append(flights, fs)
	}
	return flights, nil
}

// FetchFlightDetail retrieves a single flight by ICAO24 identifier.
func (c *OpenSkyClient) FetchFlightDetail(icao24 string) (*FlightState, error) {
	url := fmt.Sprintf("%s/states/all", c.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("icao24", icao24)
	req.URL.RawQuery = q.Encode()
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	client := &http.Client{Timeout: c.timeout}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("OpenSky detail fetch failed: %v", err)
		return nil, nil
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		States [][]interface{} `json:"states"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, nil
	}
	if len(data.States) > 0 {
		fs := parseState(data.States[0])
		return &fs, nil
	}
	return nil, nil
}
