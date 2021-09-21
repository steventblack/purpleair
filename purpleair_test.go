package purpleair

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

var km map[string]string

func init() {
	f, err := ioutil.ReadFile("./keys.JSON")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(f, &km)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Keys initialized\n")
}

func TestKeys(t *testing.T) {
	kt, err := CheckAPIKey(km["read"])
	if err != nil {
		t.Log("Unable to CheckAPIKey", err)
		t.Fail()
	}
	if kt != APIKEYREAD {
		t.Log("Expected read key, got", kt)
		t.Fail()
	}
	kt, err = CheckAPIKey(km["write"])
	if err != nil {
		t.Log("Unable to CheckAPIKey", err)
		t.Fail()
	}
	if kt != APIKEYWRITE {
		t.Log("Expected write key, got", kt)
		t.Fail()
	}
	kt, err = CheckAPIKey("BOGUS")
	if err == nil {
		t.Log("Missing error for bogus key on CheckAPIKey")
		t.Fail()
	}
	if kt != APIKEYUNKNOWN {
		t.Log("Expected unknown key, got", kt)
		t.Fail()
	}
}

func TestGroups(t *testing.T) {
}

func TestSensorInfo(t *testing.T) {
}
