package stockfighter

import (
	"testing"
)

var c *Client = &Client{}

func TestLevelControls(t *testing.T) {
	c.setAuthentication()
	level, err := NewLevel("first_steps", c)
	if err != nil {
		t.Fatal(err)
	}
	if !level.Ok {
		t.Errorf("Level Not OKAY after Starting it: %v", level)
	}

	err = level.RestartLevel(c)
	if err != nil {
		t.Errorf("Error Restarting Level: %s", err)
		t.Fatal(err)
	}

	err = level.StopLevel(c)
	if err != nil {
		t.Errorf("Error stopping Level: %s", err)
		t.Fatal(err)
	}
	if !level.Ok {
		t.Errorf("Level Not OKAY after Stopping: %v", level)
	}
}

func TestMakeRequestHeartbeat(t *testing.T) {
	c := &Client{}
	c.setAuthentication()
	answer, err := c.MakeRequest("GET", API_ENDPOINT+"/heartbeat", nil)
	t.Logf("answer: %+v", answer)
	if err != nil {
		t.Fatal(err)
	}
	if answer == nil {
		t.Fatal("nil answer")
	}
}

func TestMakeRequestVenues(t *testing.T) {
	c := &Client{}
	c.setAuthentication()
	answer, err := c.MakeRequest("GET", API_ENDPOINT+"/venues/TESTEX/heartbeat", nil)
	if err != nil {
		t.Fatal(err)
	}
	if answer == nil {
		t.Fatal(answer)
	}
}

func TestIsLevelActive(t *testing.T) {
	c := &Client{}
	c.setAuthentication()

}
