package stockfighter

import (
	_ "fmt"
	"testing"
)

func TestHeartbeatApi(t *testing.T) {
	heartbeat, err := GetHeartbeat()
	if err != nil {
		t.Error(err)
	}
	if !heartbeat.Ok {
		t.Fail()
	}
	t.Logf("Heartbeat: %+v", heartbeat)
}

func TestVenueHeartbeatApi(t *testing.T) {
	heartbeat, err := GetVenueHeartbeat("TESTEX")
	if err != nil {
		t.Error(err)
	}
	if !heartbeat.Ok {
		t.Fail()
	}
	if heartbeat.Venue != "TESTEX" {
		t.Fail()
	}
	t.Logf("Venue Heartbeat: %+v", heartbeat)
}
