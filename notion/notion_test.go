package notion

import (
	"strings"
	"testing"
)

func TestHello(t *testing.T) {
	apiKey := getAPIKey()
	if !strings.HasPrefix(apiKey, "secret_") {
		t.Fatalf("Wrong Api Key found %s", apiKey)
	}
}

func Test_addRowToDB(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "init",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addRowToDB()
		})
	}
}
