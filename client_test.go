package stockfighter

import (
	"fmt"
	"os"
	"testing"
)

func TestStartLevel(t *testing.T) {
	output, err := StartLevel("first_steps", os.Getenv(API_KEY_ENV))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", output)
}
