package stockfighter

import (
	"fmt"
	"testing"
)

var c *Client = &Client{}

func TestStartLevel(t *testing.T) {
	output, err := c.StartLevel("first_steps")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", output)
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
