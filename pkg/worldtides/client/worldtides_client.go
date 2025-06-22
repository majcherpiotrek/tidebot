package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"tidebot/pkg/worldtides/models"

	"github.com/labstack/echo/v4"
)

const (
	WorldTidesAPIURL = "https://www.worldtides.info/api/v3"
	// Hardcoded coordinates for Fuerteventura, Canary Islands
	DefaultLatitude  = 28.110419112734185
	DefaultLongitude = -14.260264983464896
)

type WorldTidesClient interface {
	GetTidalHeightsForDay(date string) (*models.WorldTidesResponse, error)
	GetTidalExtremesForDay(date string) (*models.WorldTidesResponse, error)
	GetTidalHeightsAndExtremesForDay(date string) (*models.WorldTidesResponse, error)
}

type worldTidesClientImpl struct {
	apiKey     string
	httpClient *http.Client
	log        echo.Logger
}

func NewWorldTidesClient(apiKey string, log echo.Logger) WorldTidesClient {
	return &worldTidesClientImpl{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		log: log,
	}
}

func (c *worldTidesClientImpl) GetTidalHeightsForDay(date string) (*models.WorldTidesResponse, error) {
	c.log.Debugf("Getting tidal heights for date: %s", date)

	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("lat", strconv.FormatFloat(DefaultLatitude, 'f', -1, 64))
	params.Set("lon", strconv.FormatFloat(DefaultLongitude, 'f', -1, 64))
	params.Set("date", date)
	params.Set("heights", "")
	params.Set("datum", "CD") // Chart Datum (recommended)

	return c.makeRequest(params)
}

func (c *worldTidesClientImpl) GetTidalExtremesForDay(date string) (*models.WorldTidesResponse, error) {
	c.log.Debugf("Getting tidal extremes for date: %s", date)

	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("lat", strconv.FormatFloat(DefaultLatitude, 'f', -1, 64))
	params.Set("lon", strconv.FormatFloat(DefaultLongitude, 'f', -1, 64))
	params.Set("date", date)
	params.Set("extremes", "")
	params.Set("datum", "CD") // Chart Datum (recommended)

	return c.makeRequest(params)
}

func (c *worldTidesClientImpl) GetTidalHeightsAndExtremesForDay(date string) (*models.WorldTidesResponse, error) {
	c.log.Debugf("Getting tidal heights and extremes for date: %s", date)

	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("lat", strconv.FormatFloat(DefaultLatitude, 'f', -1, 64))
	params.Set("lon", strconv.FormatFloat(DefaultLongitude, 'f', -1, 64))
	params.Set("date", date)
	params.Set("heights", "")
	params.Set("extremes", "")
	params.Set("datum", "CD") // Chart Datum (recommended)

	return c.makeRequest(params)
}

func (c *worldTidesClientImpl) makeRequest(params url.Values) (*models.WorldTidesResponse, error) {
	requestURL := fmt.Sprintf("%s?%s", WorldTidesAPIURL, params.Encode())
	
	c.log.Debugf("Making WorldTides API request to: %s", requestURL)

	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	c.log.Debugf("WorldTides API response status: %d, body length: %d bytes", resp.StatusCode, len(body))
	c.log.Debugf("WorldTides API response body: %s", string(body))

	var worldTidesResponse models.WorldTidesResponse
	err = json.Unmarshal(body, &worldTidesResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for API errors
	if worldTidesResponse.Status != 200 {
		return nil, fmt.Errorf("WorldTides API error (status %d): %s", worldTidesResponse.Status, worldTidesResponse.Error)
	}

	c.log.Debugf("Successfully received WorldTides data - Heights: %d, Extremes: %d", 
		len(worldTidesResponse.Heights), len(worldTidesResponse.Extremes))

	return &worldTidesResponse, nil
}