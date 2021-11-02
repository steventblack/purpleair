package purpleair

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

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
	sp[SensorParamShowOnly] = fmt.Sprintf("%d", ti.SensorInfo.TestSensorIndex)
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
	var si []string
	for _, i := range ti.SensorInfo.TestSensorsIndex {
		si = append(si, strconv.Itoa(int(i)))
	}

	sp[SensorParamShowOnly] = strings.Join(si, ",")
	_, err = SensorsData(sp)
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}
}
