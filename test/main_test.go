package test_test

import (
	"os"
	"strings"
	"testing"
)

var apiToken = os.Getenv("DC_TOKEN")

func TestMain(m *testing.M) {
	if apiToken == "" {
		panic("Please specify an envoirement variable with the name " +
			"'DC_TOKEN' and the value of a Discord API App token.\n" +
			"Also, the App's user should be present at least on one guild.")
	}

	m.Run()
}

func assertNotEqual(t *testing.T, expected, actual interface{}, msg ...string) {
	if expected == actual {
		if len(msg) > 0 {
			t.Fatal(strings.Join(msg, " "))
		}
		t.Fatal("value was equal where unexpected: ", actual)
	}
}
