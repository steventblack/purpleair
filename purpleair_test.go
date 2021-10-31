package purpleair

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

// Standard test information.
// Keys are available by requesting from "contact@purpleair.com".
// Testing requires both the read and write keys.
// Because there is a highly variable propagation delay in PurpleAir's group and membership
// APIs, a group should be created prior to testing and a known sensor added to the group.
type TestInfo struct {
	Keys      map[string]string `json:"keys"`
	GroupInfo map[string]int    `json:"groupinfo"`
	Fields    string            `json:"fields"`
}

var ti TestInfo

const (
	TESTGROUP     string      = "testing_group"
	TESTSENSORIDX SensorIndex = 118475
	TESTFIELDS    string      = "sensor_index,name,location_type,hardware,latitude,longitude,rssi,model"
)

// Initialization of read & write keys used for API access
// If keys are not available, then it is unable to perform any API tests
func init() {
	f, err := ioutil.ReadFile("./keys.JSON")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(f, &ti)
	if err != nil {
		log.Fatal(err)
	}
}

// Suite of tests for Group usage
// includes: creation/deletion, listing, details, and membership (add/remove)
// Ordering of the tests is important (e.g. you can't delete the group until after you've created it)
// The PurpleAir API also suffers from internal propagation latency, so there may be race conditions.
/*
func TestGroups(t *testing.T) {
	// create a group
	// the group id (g) will be used for the other group related calls
	g, err := CreateGroup(TESTGROUP)
	if err != nil {
		t.Log("Unable to CreateGroup", err)
		t.Fail()
	}

	// add a member by sensor_index
	m, err := TESTSENSORIDX.AddMember(g)
	if err != nil {
		t.Log("Unable to AddMember by SensorIndex", err)
		t.Fail()
	}

	// insert sleep to allow data to sync within PurpleAir
	// and help smooth over race conditions within their infrastructure
	time.Sleep(3 * time.Second)

	// list all groups
	// validate group created is present with matching id & name
	gl, err := ListGroups()
	if err != nil {
		t.Log("Unable ListGroups", err)
		t.Fail()
	}

	foundGroup := false
	for _, v := range gl {
		if v.ID == g {
			foundGroup = true
			if v.Name != TESTGROUP {
				t.Logf("Group name mismatch %s vs %s\n", TESTGROUP, v.Name)
				t.Fail()
			}
			break
		}
	}
	if foundGroup != true {
		t.Logf("Unable to find group %d in GroupList\n", g)
		t.Fail()
	}

	// get group membership details
	// verify the member added to the group earlier is present
	ml, err := GroupDetails(g)
	if err != nil {
		t.Log("Unable to get GroupDetails", err)
		t.Fail()
	}

	foundMember := false
	for _, v := range ml {
		if v.ID == m {
			foundMember = true
			break
		}
	}
	if foundMember != true {
		t.Logf("Unable to find member %d in group %d\n", m, g)
		t.Fail()
	}

	// fetch a group member's data
	var mp = make(SensorParams)
	_, err = MemberData(g, m, mp)
	if err != nil {
		t.Log("Unable to get member data", err)
		t.Fail()
	}

	//fp := SensorFields{Fields: TESTFIELDS}
	mp[SensorParamFields] = TESTFIELDS
	_, err = MemberData(g, m, mp)
	if err != nil {
		t.Log("Unable to get member data with fields", err)
		t.Fail()
	}

	var sp = make(SensorParams)
	sp[SensorParamFields] = TESTFIELDS
	_, err = MembersData(599, sp)
	if err != nil {
		t.Log("Unable to get all member data", err)
		t.Fail()
	}

	// remove the group member
	err = RemoveMember(m, g)
	if err != nil {
		t.Logf("Unable to remove member %d from group %d\n", m, g)
		t.Fail()
	}

	// delete the group
	err = DeleteGroup(g)
	if err != nil {
		t.Logf("Unable to DeleteGroup %d %s\n", g, err)
		t.Fail()
	}
}
*/

// Suite of tests for retriving sensor info
func TestSensorInfo(t *testing.T) {
	// test fetching all data for a sensor
	var mp = make(SensorParams)
	sd, err := SensorData(TESTSENSORIDX, mp)
	if err != nil {
		t.Log("Unable to get sensor data", err)
		t.Fail()
	}
	t.Logf("SensorData (all):\n%v+\n", sd)

	// test fetching selected data for a sensor
	//	fp := SensorFields{Fields: TESTFIELDS}
	mp[SensorParamFields] = TESTFIELDS
	sd, err = SensorData(TESTSENSORIDX, mp)
	if err != nil {
		t.Log("Unable to get sensor data with fields", err)
		t.Fail()
	}
	t.Logf("SensorData:\n%v+\n", sd)
}

func TestSensorParams(t *testing.T) {
	/*
		// testing param block
		var p = make(SensorParams)
		p[SensorParamFields] = "sensor_index,name,latitude,longitude,location_type,model"
		p[SensorParamLocation] = OUTSIDE
		p[SensorParamNWLong] = 123.456

		_, err := processParams(p)
		if err != nil {
			t.Log("Unable to process sensor params", err)
			t.Fail()
		}

		// setup a params block without the required fields
		var pf = make(SensorParams)
		p[SensorParamLocation] = OUTSIDE
		p[SensorParamNWLong] = 123.456

		_, err = processParams(pf)
		if err == nil {
			t.Log("Missing error for missing required 'fields' param")
			t.Fail()
		}

		var pb = make(SensorParams)
		p[SensorParamFields] = "sensor_index,name,latitude,longitude,location_type,model"
		p[SensorParamLocation] = OUTSIDE
		p[SensorParamNWLong] = 123.456
		p["bogus"] = "this invalid key better throw an error"

		_, err = processParams(pb)
		if err == nil {
			t.Log("Missing error for passing invalid parameter key")
			t.Fail()
		}
	*/
}
