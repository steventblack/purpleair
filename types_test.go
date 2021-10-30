package purpleair

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnmarshalGroup(t *testing.T) {
	var g Group
	d := []byte(`{ "id": 123, "name": "fruitbat", "created": 1635575338 }`)

	err := json.Unmarshal(d, &g)
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}

	if g.ID != 123 {
		t.Logf("%s: Expected ID %d, got %d\n", t.Name(), 123, g.ID)
		t.Fail()
	}

	if g.Name != "fruitbat" {
	}

	ts := time.Unix(1635575338, 0)
	if g.Created != ts {
		t.Logf("%s: Expected created %v, got %v\n", t.Name(), ts, g.Created)
		t.Fail()
	}

}

func TestUnmarshalMember(t *testing.T) {
	var m Member
	d := []byte(`{ "id": 123, "sensor_index": 456, "created": 1635575338 }`)

	err := json.Unmarshal(d, &m)
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}

	if m.ID != 123 {
		t.Logf("%s: Expected ID %d, got %d\n", t.Name(), 123, m.ID)
		t.Fail()
	}

	if m.Index != 456 {
		t.Logf("%s: Expected Index %d, got %d\n", t.Name(), 456, m.ID)
		t.Fail()
	}

	ts := time.Unix(1635575338, 0)
	if m.Created != ts {
		t.Logf("%s: Expected created %v, got %v\n", t.Name(), ts, m.Created)
		t.Fail()
	}
}
