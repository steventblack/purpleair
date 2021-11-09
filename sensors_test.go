package purpleair

import (
	"testing"
)

func init() {
	initTestInfo()
}

func TestSensorData(t *testing.T) {
	// Test with no params specified
	var sp = make(SensorParams)
	_, err := SensorData(ti.SensorInfo.TestSensorIndex, sp)
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}

	// Test with a fields param
	f := ParamFields{Value: ti.SensorParams["fields"]}
	sp = f.AddParam(sp)
	_, err = SensorData(ti.SensorInfo.TestSensorIndex, sp)
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}

	// Test with an illegal param
	loc := ParamLocation{Value: LocOutside}
	sp = loc.AddParam(sp)
	_, err = SensorData(ti.SensorInfo.TestSensorIndex, sp)
	if err == nil {
		t.Log(t.Name(), err)
		t.Fail()
	}
}

func TestSensorsData(t *testing.T) {
	// test without the required fields param
	var sp = make(SensorParams)
	si := ParamShowOnly{Value: ti.SensorInfo.TestSensorsIndex}
	sp = si.AddParam(sp)
	_, err := SensorsData(sp)
	if err == nil {
		t.Log(t.Name(), err)
		t.Fail()
	}

	// test with the required fields param
	f := ParamFields{Value: ti.SensorParams["fields"]}
	sp = f.AddParam(sp)
	_, err = SensorsData(sp)
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}
}
