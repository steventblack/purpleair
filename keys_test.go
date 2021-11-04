package purpleair

import (
	"testing"
)

func init() {
	initTestInfo()
}

func TestCheckAPIKey(t *testing.T) {
	kt, err := CheckAPIKey(ti.Keys["read"])
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}
	if kt != KeyRead {
		t.Logf("%s: Expected %s, got %s\n", t.Name(), KeyRead, kt)
		t.Fail()
	}

	kt, err = CheckAPIKey(ti.Keys["write"])
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}
	if kt != KeyWrite {
		t.Logf("%s: Expected %s, got %s\n", t.Name(), KeyWrite, kt)
		t.Fail()
	}

	kt, err = CheckAPIKey("BOGUS")
	if err == nil {
		t.Logf("%s: Expected error, got nil\n", t.Name())
		t.Fail()
	}
	if kt != KeyUnknown {
		t.Logf("%s: Expected %s, got %s\n", t.Name(), KeyUnknown, kt)
		t.Fail()
	}
}

func TestSetAPIKey(t *testing.T) {
	kt, err := SetAPIKey(ti.Keys["read"])
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}
	if kt != KeyRead {
		t.Logf("%s: Expected %s, got %s\n", t.Name(), KeyRead, kt)
		t.Fail()
	}

	kt, err = SetAPIKey(ti.Keys["write"])
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}
	if kt != KeyWrite {
		t.Logf("%s: Expected %s, got %s\n", t.Name(), KeyWrite, kt)
		t.Fail()
	}

	kt, err = SetAPIKey("BOGUS")
	if err == nil {
		t.Logf("%s: Expected error, got nil\n", t.Name())
		t.Fail()
	}
	if kt != KeyUnknown {
		t.Logf("%s: Expected %s, got %s\n", t.Name(), KeyUnknown, kt)
		t.Fail()
	}
}
