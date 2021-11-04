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
	sp[SensorParamFields] = ti.SensorParams["fields"]
	_, err = SensorData(ti.SensorInfo.TestSensorIndex, sp)
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}

	// Test with an illegal param
	sp[SensorParamLocation] = LocOutside
	_, err = SensorData(ti.SensorInfo.TestSensorIndex, sp)
	if err == nil {
		t.Log(t.Name(), err)
		t.Fail()
	}
}

func TestSensorsData(t *testing.T) {
	// test without the required fields param
	var sp = make(SensorParams)
	sp[SensorParamShowOnly] = ti.SensorInfo.TestSensorsIndex
	_, err := SensorsData(sp)
	if err == nil {
		t.Log(t.Name(), err)
		t.Fail()
	}

	// test with the required fields param
	sp[SensorParamFields] = ti.SensorParams["fields"]
	_, err = SensorsData(sp)
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}

	// test with more than one sensor index specified
	sp[SensorParamShowOnly] = ti.SensorInfo.TestSensorsIndex
	_, err = SensorsData(sp)
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}
}
