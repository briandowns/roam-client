package roam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const roamURL = "https://api.tile38.com"
const ua = "roam-client"

// Point
type Point struct {
	// Lat is the latitude coordinate
	Lat float64 `json:"lat"`
	// Lon is the longitude coordinate
	Lon float64 `json:"lon"`
	// ID is the identifier of the point
	ID string `json:"id"`
	// Meta is optional metadata that is returned with
	// all notifications. Limited to 256 characters.
	Meta string `json:"meta"`
}

// Hook
type Hook struct {
	// Name is the name of the geofence
	Name string `json:"name"`
	// The radius of the geofence in meters
	Meters string `json:"meters"`
	// Match is a glob pattern that is used to match against point IDs.
	// All points with IDs matching this value will become a geofence
	Match string `json:"match"`
	// Endpoint is a valid HTTP URL
	Endpoint string `json:"endpoint"`
	// Filter is an optional glob pattern that is used to filter notifications
	// based on the target point IDs
	Filter string `json:"filter"`
}

// ClientOpts
type ClientOpts struct {
	Timeout   time.Duration
	UserAgent string
	call      caller
}

// Client
type Client struct {
	token string
	ua    string
	hc    *http.Client
}

// New
func New(token string, co *ClientOpts) *Client {
	c := Client{
		token: token,
		hc:    &http.Client{},
	}
	if co != nil {
		if co.Timeout > 0 {
			c.hc.Timeout = co.Timeout
		}
		switch {
		case co.UserAgent != "":
			c.ua = co.UserAgent
		default:
			c.ua = ua
		}
	}
	return &c
}

// Notification
type Notification struct {
	UID     string `json:"uid"`
	Service string `json:"service"`
	Hook    string `json:"hook"`
	Time    string `json:"time"`
	ID      string `json:"id"`
	Object  Object `json:"object"`
	NearBy  NearBy `json:"nearby"`
}

// NearBy
type NearBy struct {
	ID     string  `json:"id"`
	Object Object  `json:"object"`
	Meters float64 `json:"meters"`
}

// Object
type Object struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// SetPoint
func (c *Client) SetPoint(p *Point) (bool, error) {
	var r response
	var params string
	switch {
	case len(p.Meta) > 0:
		params = fmt.Sprintf("/setpoint?lat=%f&lon=%f&id=%s&meta=%s", p.Lat, p.Lon, p.ID, p.Meta)
	default:
		params = fmt.Sprintf("/setpoint?lat=%f&lon=%f&id=%s", p.Lat, p.Lon, p.ID)
	}
	if err := c.call(http.MethodGet, roamURL+params, &r); err != nil {
		return false, err
	}
	return r.OK, nil
}

// SetHook
func (c *Client) SetHook(h *Hook) (bool, error) {
	var r response
	params := fmt.Sprintf("/sethook?name=%s&meters=%s&match=%s&endpoint=%s", h.Name, h.Meters, h.Match, h.Endpoint)
	if err := c.call(http.MethodGet, roamURL+params, &r); err != nil {
		return false, err
	}
	return r.OK, nil
}

// DeleteHook
func (c *Client) DeleteHook(name string) (bool, error) {
	var r response
	if err := c.call(http.MethodGet, roamURL+"/delhook?name="+name, &r); err != nil {
		return false, err
	}
	return r.OK, nil
}

// Hooks
func (c *Client) Hooks() ([]Hook, error) {
	var h hookResponse
	if err := c.call(http.MethodGet, roamURL+"/hooks", &h); err != nil {
		return nil, err
	}
	return h.Hooks, nil
}

// response
type response struct {
	OK bool `json:"ok"`
}

// hookResponse
type hookResponse struct {
	response
	Hooks []Hook `json:"hooks"`
}

// buildRequest setups up the request to the Roam API
func (c *Client) call(m, url string, result interface{}) error {
	req, err := http.NewRequest(m, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("User-Agent", c.ua)
	res, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	return json.NewDecoder(res.Body).Decode(&result)
}

// caller
type caller interface {
	call(m, url string, result interface{}) error
}
