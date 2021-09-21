package purpleair

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

// Key map for holding the read & write keys for API access
// Initialized by reading in local JSON file "keys.JSON"
// Keys are available by requesting from "contact@purpleair.com"
var km map[string]string

// Initialization of read & write keys used for API access
// If keys are not available, then unable to perform any API tests
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

// Suite of tests related to key validation
func TestKeys(t *testing.T) {
	// read key check
	kt, err := CheckAPIKey(km["read"])
	if err != nil {
		t.Log("Unable to CheckAPIKey", err)
		t.Fail()
	}
	if kt != APIKEYREAD {
		t.Log("Expected read key, got", kt)
		t.Fail()
	}

	// write key check
	kt, err = CheckAPIKey(km["write"])
	if err != nil {
		t.Log("Unable to CheckAPIKey", err)
		t.Fail()
	}
	if kt != APIKEYWRITE {
		t.Log("Expected write key, got", kt)
		t.Fail()
	}

	// bogus key check
	kt, err = CheckAPIKey("BOGUS")
	if err == nil {
		t.Log("Missing error for bogus key on CheckAPIKey")
		t.Fail()
	}
	if kt != APIKEYUNKNOWN {
		t.Log("Expected unknown key, got", kt)
		t.Fail()
	}

	// read key set
	kt, err = SetAPIKey(km["read"])
	if err != nil {
		t.Log("Unable to SetAPIKey", err)
		t.Fail()
	}
	if kt != APIKEYREAD {
		t.Log("Expected read key, got", kt)
		t.Fail()
	}

	// write key set
	kt, err = SetAPIKey(km["write"])
	if err != nil {
		t.Log("Unable to SetAPIKey", err)
		t.Fail()
	}
	if kt != APIKEYWRITE {
		t.Log("Expected write key, got", kt)
		t.Fail()
	}

	// bogus key set
	kt, err = SetAPIKey("bogus")
	if err == nil {
		t.Log("Missing error for bogus key on SetAPIKey")
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
