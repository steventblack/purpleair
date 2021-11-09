package purpleair

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// CreateGroup creates a PurpleAir group collection returning the GroupID reference.
// Sensors can then be added to the group to build up the set.
// Requires PurpleAir write permissions.
func CreateGroup(g string) (GroupID, error) {
	params := struct {
		G string `json:"name"`
	}{G: g}

	data, err := json.Marshal(params)
	if err != nil {
		return 0, err
	}

	u, err := url.Parse(urlGroups)
	if err != nil {
		return 0, err
	}

	r, err := doRequest(http.MethodPost, u, data)
	if err != nil {
		return 0, nil
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusCreated {
		return 0, paError(r)
	}

	payload := struct {
		G int `json:"group_id"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		return 0, err
	}

	return GroupID(payload.G), nil
}

// DeleteGroup removes the PurpleAir group collection.
// All members must be removed prior to group deletion or an error will result.
// Requires PurpleAir write permissions.
func DeleteGroup(g GroupID) error {
	u, err := url.Parse(fmt.Sprintf("%s/%d", urlGroups, g))
	if err != nil {
		return err
	}

	r, err := doRequest(http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusNoContent {
		return paError(r)
	}

	return nil
}

// ListGroups lists all available PurpleAir group collections associated with the account.
// Requires PurpleAir read permissions.
func ListGroups() ([]Group, error) {
	u, err := url.Parse(urlGroups)
	if err != nil {
		return nil, err
	}

	r, err := doRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, paError(r)
	}

	payload := struct {
		G []Group `json:"groups"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		return nil, err
	}

	return payload.G, nil
}

// ListGroupMembers lists all members belonging to the specified group.
// Requires PurpleAir read permissions.
func ListGroupMembers(g GroupID) ([]Member, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%d", urlGroups, g))
	if err != nil {
		return nil, err
	}

	r, err := doRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, paError(r)
	}

	payload := struct {
		M []Member `json:"members"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		return nil, err
	}

	return payload.M, nil
}

// SensorIndex implementation for the AddMember interface function.
// Adds the sensor to the specified group.
// The PrivateInfo optional argument is for private sensors which require
// additional validation for membership assignment.
// Requires PurpleAir write permissions.
func (s SensorIndex) AddMember(g GroupID, pi ...PrivateInfo) (MemberID, error) {
	params := struct {
		S SensorIndex `json:"sensor_index"`
		E string      `json:"owner_email,omitempty"`
		L Location    `json:"location_type,omitempty"`
	}{S: s}

	// If private info is supplied, include it in the request params
	if pi != nil {
		params.E = pi[0].Email
		params.L = pi[0].Loc
	}

	data, err := json.Marshal(params)
	if err != nil {
		return 0, err
	}

	return addMember(g, data)
}

// SensorID implementation for the AddMember interface function.
// Adds the sensor to the specified group.
// The PrivateInfo optional argument is for private sensors which require
// additional validation for membership assignment.
// Requires PurpleAir write permissions.
func (s SensorID) AddMember(g GroupID, pi ...PrivateInfo) (MemberID, error) {
	params := struct {
		S SensorID `json:"sensor_id"`
		E string   `json:"owner_email,omitempty"`
		L Location `json:"location_type,omitempty"`
	}{S: s}

	// If private info is supplied, include it in the request params
	if pi != nil {
		params.E = pi[0].Email
		params.L = pi[0].Loc
	}

	data, err := json.Marshal(params)
	if err != nil {
		return 0, err
	}

	return addMember(g, data)
}

// Private function of common code supporting the AddMember interface functions.
// Requires PurpleAir write permissions.
func addMember(g GroupID, data []byte) (MemberID, error) {
	u, err := url.Parse(fmt.Sprintf(urlMembers, g))
	if err != nil {
		return 0, err
	}

	r, err := doRequest(http.MethodPost, u, data)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusCreated {
		return 0, paError(r)
	}

	payload := struct {
		M MemberID `json:"member_id"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		return 0, err
	}

	return payload.M, nil
}

// Remove the member from the specified group.
// Requires PurpleAir write permissions.
func RemoveMember(m MemberID, g GroupID) error {
	u, err := url.Parse(fmt.Sprintf(urlMembers+"/%d", g, m))
	if err != nil {
		return err
	}

	r, err := doRequest(http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusNoContent {
		return paError(r)
	}

	return nil
}

// MemberData returns the SensorInfo for a member of a group.
// The SensorParams can restrict the information returned to the named fields.
// This call requires a key with read permissions to be set prior to calling.
// On success, the SensorInfo will be returned, or else an error.
// Note that if a subset of fields is specified, only that data will be returned.
func MemberData(g GroupID, m MemberID, sp SensorParams) (*SensorInfo, error) {
	u, err := url.Parse(fmt.Sprintf(urlMembers+"/%d", g, m))
	if err != nil {
		return nil, err
	}

	// check for permitted/required params
	for k, _ := range sp {
		switch k {
		case paramFields:
		default:
			return nil, fmt.Errorf("Unexpected sensor param encountered [%s]", k)
		}
	}

	return paSensor(u, sp)
}

// MembersData returns the information requested for the set (or subset)
// of sensors within the specified Group. The SensorParams must specify
// the elements requested in the "fields" parameter.
// The return value is a map of key/value pairs for each field element
// specified indexed by the sensor_index.
func MembersData(g GroupID, sp SensorParams) (SensorDataSet, error) {
	u, err := url.Parse(fmt.Sprintf(urlMembers, g))
	if err != nil {
		return nil, err
	}

	// check for permitted/required params
	requiredField := false
	for k, _ := range sp {
		switch k {
		case paramFields:
			requiredField = true
		case paramLocation, paramReadKeys, paramShowOnly, paramModTime, paramMaxAge:
		case paramNWLong, paramNWLat, paramSELong, paramSELat:
		default:
			return nil, fmt.Errorf("Unexpected sensor param encountered [%s]", k)
		}
	}

	if requiredField == false {
		return nil, fmt.Errorf("Required sensor param not found [%s]", paramFields)
	}

	return paSensors(u, sp)
}
