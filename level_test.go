package stockfighter

import (
	"fmt"
	"testing"
	"time"
)

func TestStartLevel(t *testing.T) {
	c.setAuthentication()
	level, err := NewLevel("first_steps", c)
	if err != nil {
		t.Fatal(err)
	}
	var tests = []struct {
		in  string
		out string
	}{
		{level.Error, ""},
		{fmt.Sprintf("%t", level.Ok), "true"},
	}
	for _, tt := range tests {
		if tt.in != tt.out {
			t.Errorf("Expected %s, Got %s", tt.in, tt.out)
		}
	}
}

func TestLevelControls(t *testing.T) {
	c.setAuthentication()
	level, err := NewLevel("chock_a_block", c)
	if err != nil {
		t.Errorf("StartLevel error: %s", err)
	}
	if !level.Ok {
		t.Errorf("Level Not OKAY after Starting it: %v", level)
	}
	time.Sleep(time.Second * 5)
	err = level.RestartLevel(c)
	if err != nil {
		t.Errorf("Error Restarting Level: %s\n", err)
	}
	time.Sleep(time.Second * 5)
	err = level.StopLevel(c)
	if err != nil {
		t.Errorf("Error stopping Level: %s", err)
		t.Fatal(err)
	}
	if !level.Ok {
		t.Errorf("Level Not OKAY after Stopping: %v", level)
	}
}

func TestResumeLevel(t *testing.T) {
	c := &Client{}
	c.setAuthentication()
	l, err := NewLevel("first_steps", c)
	if err != nil {
		t.Error(err)
	}
	err = l.ResumeLevel(c)
	if err != nil {
		t.Error(err)
	}
}

func TestIsLevelActive(t *testing.T) {
	c := &Client{}
	c.setAuthentication()

}
