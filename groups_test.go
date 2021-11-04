package purpleair

import (
	"testing"
)

func init() {
	initTestInfo()

	SetAPIKey(ti.Keys["read"])
	SetAPIKey(ti.Keys["write"])
}

// Group and Membership calls often have a propagation delay within PurpleAir.
func TestGroup(t *testing.T) {
	g, err := testCreateGroup(t, "test_fruitbat")
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
		t.SkipNow()
	}

	m, err := testAddMember(t, g, SensorIndex(ti.GroupInfo["testSensorIndex"]))
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
		t.SkipNow()
	}

	err = RemoveMember(m, g)
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}

	err = DeleteGroup(g)
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}
}

func TestListGroups(t *testing.T) {
	_, err := ListGroups()
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}
}

func TestListGroupMembers(t *testing.T) {
	_, err := ListGroupMembers(GroupID(ti.GroupInfo["testGroupID"]))
	if err != nil {
		t.Log(t.Name(), err)
		t.Fail()
	}
}

func testCreateGroup(t *testing.T, n string) (GroupID, error) {
	g, err := CreateGroup(n)
	if err != nil {
		return 0, err
	}

	t.Cleanup(func() {
		DeleteGroup(g)
	})

	return g, err
}

func testAddMember(t *testing.T, g GroupID, s SensorIndex) (MemberID, error) {
	m, err := s.AddMember(g)
	if err != nil {
		return 0, err
	}

	t.Cleanup(func() {
		RemoveMember(m, g)
	})

	return m, err
}
