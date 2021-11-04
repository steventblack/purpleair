package purpleair

import (
	"encoding/json"
	"testing"
	"time"
)

func init() {
	initTestInfo()
}

func TestUnmarshalSensorStats(t *testing.T) {
	var s SensorStats
	d := []byte(`{ "pm2.5": 123.1, 
		"pm2.5_10minute": 123.2,
		"pm2.5_30minute": 123.3,
		"pm2.5_60minute": 123.4,
		"pm2.5_6hour": 123.5,
		"pm2.5_24hour": 123.6,
		"pm2.5_1week": 123.7,
		"time_stamp": 1635575338 }`)

	err := json.Unmarshal(d, &s)
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}

	if s.PM_2_5 != 123.1 {
		t.Logf("%s: Expected PM_2_5 %f, got %f\n", t.Name(), 123.1, s.PM_2_5)
		t.Fail()
	}

	if s.PM_2_5_10Min != 123.2 {
		t.Logf("%s: Expected PM_2_5_10Min %f, got %f\n", t.Name(), 123.2, s.PM_2_5_10Min)
		t.Fail()
	}

	if s.PM_2_5_30Min != 123.3 {
		t.Logf("%s: Expected PM_2_5_30Min %f, got %f\n", t.Name(), 123.3, s.PM_2_5_30Min)
		t.Fail()
	}

	if s.PM_2_5_60Min != 123.4 {
		t.Logf("%s: Expected PM_2_5_60Min %f, got %f\n", t.Name(), 123.4, s.PM_2_5_60Min)
		t.Fail()
	}

	if s.PM_2_5_6Hour != 123.5 {
		t.Logf("%s: Expected PM_2_5_6Hour %f, got %f\n", t.Name(), 123.5, s.PM_2_5_6Hour)
		t.Fail()
	}

	if s.PM_2_5_24Hour != 123.6 {
		t.Logf("%s: Expected PM_2_5_24Hour %f, got %f\n", t.Name(), 123.6, s.PM_2_5_24Hour)
		t.Fail()
	}

	if s.PM_2_5_1Week != 123.7 {
		t.Logf("%s: Expected PM_2_5_1Week %f, got %f\n", t.Name(), 123.7, s.PM_2_5_1Week)
		t.Fail()
	}

	ts := time.Unix(1635575338, 0)
	if s.Timestamp != ts {
		t.Logf("%s: Expected Timestamp %v, got %v\n", t.Name(), ts, s.Timestamp)
		t.Fail()
	}

}

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
		t.Logf("%s: Expected Name %s, got %s\n", t.Name(), "fruitbat", g.Name)
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
