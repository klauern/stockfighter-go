package stockfighter

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Heartbeat struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type VenueHeartbeat struct {
	Ok    bool   `json:"ok"`
	Venue string `json:"venue,omitempty"`
	Error string `json:"error,omitempty"`
}

func GetHeartbeat() (*Heartbeat, error) {
	resp, err := http.Get(API_ENDPOINT + "heartbeat")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var heartbeat Heartbeat
	err = json.Unmarshal(body, &heartbeat)
	if err != nil {
		return nil, err
	}
	if !heartbeat.Ok || heartbeat.Error != "" {
		return nil, errors.New(heartbeat.Error)
	}
	return &heartbeat, nil
}

func GetVenueHeartbeat(venue string) (*VenueHeartbeat, error) {
	resp, err := http.Get(API_ENDPOINT + "venues/" + venue + "/heartbeat")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var heartbeat VenueHeartbeat
	err = json.Unmarshal(body, &heartbeat)
	if err != nil {
		return nil, err
	}
	if !heartbeat.Ok || heartbeat.Error != "" {
		return nil, errors.New(heartbeat.Error)
	}
	return &heartbeat, nil
}
