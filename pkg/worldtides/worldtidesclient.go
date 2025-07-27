package worldtides

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	WorldTidesAPIURL = "https://www.worldtides.info/api/v3"
	// Hardcoded coordinates for Fuerteventura, Canary Islands
	DefaultLatitude  = 28.110419112734185
	DefaultLongitude = -14.260264983464896
)

type WorldTidesClient interface {
	GetTides(date time.Time) (*WorldTidesResponse, error)
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]*WorldTidesResponse
}

func NewCache() Cache {
	return Cache{
		data: make(map[string]*WorldTidesResponse),
	}
}

func (c *Cache) read(key string) (*WorldTidesResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.data[key]
	return value, exists
}

func (c *Cache) write(key string, value *WorldTidesResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

type worldTidesClientImpl struct {
	apiKey     string
	httpClient *http.Client
	log        echo.Logger
	cache      Cache
}

func NewWorldTidesClient(apiKey string, log echo.Logger) WorldTidesClient {
	return &worldTidesClientImpl{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		log:   log,
		cache: NewCache(),
	}
}

func (c *worldTidesClientImpl) GetTides(date time.Time) (*WorldTidesResponse, error) {
	dateFormatted := date.Format("2006-01-02")
	c.log.Debugf("Getting tides for date: %s", dateFormatted)

	cacheKey := fmt.Sprintf("tides:%s", dateFormatted)

	c.log.Debugf("Cache key: %s", cacheKey)
	c.log.Debugf("Cache contents: %+v", c.cache.data)

	cached, exists := c.cache.read(cacheKey)

	if exists {
		c.log.Debugf("Cache hit for tides request for %s", dateFormatted)
		return cached, nil
	}

	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("lat", strconv.FormatFloat(DefaultLatitude, 'f', -1, 64))
	params.Set("lon", strconv.FormatFloat(DefaultLongitude, 'f', -1, 64))
	params.Set("date", dateFormatted)
	params.Set("days", "1")
	params.Set("extremes", "")
	params.Set("heights", "")
	params.Set("datum", "MLS") // Mean Sea Level -- https://www.worldtides.info/datums
	params.Set("localtime", "")

	response, err := c.makeRequest(params)

	if err != nil {
		return nil, err
	}

	c.cache.write(cacheKey, response)

	return response, nil
}

func (c *worldTidesClientImpl) makeRequest(params url.Values) (*WorldTidesResponse, error) {
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

	var worldTidesResponse WorldTidesResponse
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
