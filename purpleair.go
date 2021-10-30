package purpleair

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Global variables for retaining the API access keys.
// These are set via the SetAPIKey.
var (
	apiReadKey  string
	apiWriteKey string
)

// SetAPIKey checks the validity and permissions of the provided access key string.
// If the key is valid, it will be retained by the module for further calls.
// Only one key for each permission (read or write) will be retained, and additional
// calls with other valid keys will result in replacement.
// The KeyType will be returned on success, or an error on failure.
func SetAPIKey(k string) (KeyType, error) {
	keyType, err := CheckAPIKey(k)
	if err != nil {
		return KeyUnknown, err
	}

	if keyType == KeyRead {
		apiReadKey = k
	} else if keyType == KeyWrite {
		apiWriteKey = k
	}

	return keyType, nil
}

// CheckAPIKey checks the validity and permissions of the provided access key string.
// It does *not* retain the key for further calls. Use SetAPIKey if retention is desired.
// The KeyType will be returned on success, or an error on failure.
func CheckAPIKey(k string) (KeyType, error) {
	type checkKeyType struct {
		V string  `json:"api_version"`
		T int     `json:"time_stamp"`
		K KeyType `json:"api_key_type"`
	}
	var keyTypeResp checkKeyType

	keyType := KeyUnknown

	req, err := http.NewRequest(http.MethodGet, URLKEYS, nil)
	if err != nil {
		return keyType, err
	}
	req.Header.Add(keyHeader, k)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return keyType, err
	}
	defer resp.Body.Close()

	// if invalid key, an error is returned
	if resp.StatusCode != http.StatusCreated {
		return keyType, extractError(resp)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&keyTypeResp)
	if err != nil {
		return keyType, err
	}

	return keyTypeResp.K, nil
}

// CreateGroup creates a persistent reference of a defined set of sensors on the PurpleAir service.
// Sensors can be added/removed using the membership APIs.
// This call requires a key with write permissions to be set prior to calling.
// A GroupID will be returned on success, or an error on failure.
func CreateGroup(g string) (GroupID, error) {
	reqBody := struct {
		GroupName string `json:"name"`
	}{GroupName: g}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return 0, err
	}

	u, err := url.Parse(URLGROUPS)
	if err != nil {
		return 0, err
	}

	resp, err := doRequest(http.MethodPost, u, reqJSON)
	if err != nil {
		return 0, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return 0, extractError(resp)
	}

	groupResp := struct {
		V string  `json:"api_version"`
		T int     `json:"time_stamp"`
		G GroupID `json:"group_id"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&groupResp)
	if err != nil {
		return 0, err
	}

	return groupResp.G, nil
}

// DeleteGroup deletes the persistent reference of a sensor group on the PurpleAir service.
// All sensor members must be removed prior to group deletion.
// This call requires a key with write permissions to be set prior to calling.
// An error will be returned on failure, or else nil
func DeleteGroup(g GroupID) error {
	u, err := url.Parse(fmt.Sprintf("%s/%d", URLGROUPS, g))
	if err != nil {
		return err
	}

	resp, err := doRequest(http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// if unexpected response, extract & return the error
	if resp.StatusCode != http.StatusNoContent {
		return extractError(resp)
	}

	return nil
}

// ListGroups provides all groups defined in the PurpleAir service associated with the access key.
// This call requires a key with read permissions to be set prior to calling.
// The list of groups will be returned on success, or else an error.
func ListGroups() ([]Group, error) {
	u, err := url.Parse(URLGROUPS)
	if err != nil {
		return nil, err
	}

	resp, err := doRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, extractError(resp)
	}

	groupResp := struct {
		V string  `json:"api_version"`
		T int     `json:"time_stamp"`
		D int     `json:"data_time_stamp"`
		G []Group `json:"groups"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&groupResp)
	if err != nil {
		return nil, err
	}

	return groupResp.G, nil
}

// GroupDetails provides the list of member sensors defined for the specified group.
// This call requires a key with read permissions to be set prior to calling.
// The list of members will be returned on success, or else an error.
func GroupDetails(g GroupID) ([]Member, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%d", URLGROUPS, g))
	if err != nil {
		return nil, err
	}

	resp, err := doRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, extractError(resp)
	}

	memberResp := struct {
		V string   `json:"api_version"`
		T int      `json:"time_stamp"`
		D int      `json:"data_time_stamp"`
		G GroupID  `json:"group_id"`
		M []Member `json:"members"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&memberResp)
	if err != nil {
		return nil, err
	}

	return memberResp.M, nil
}

// AddMember provides the SensorIndex interface solution to adding a sensor to a group.
// The function takes an optional PrivateInfo struct, which is necessary only if the
// sensor referenced is a private sensor.
// This call requires a key with write permissions to be set prior to calling.
// The MemberID will be returned on success, or else an error.
func (s SensorIndex) AddMember(g GroupID, pi ...PrivateInfo) (MemberID, error) {
	reqBody := struct {
		S SensorIndex `json:"sensor_index"`
		E string      `json:"owner_email,omitempty"`
		L Location    `json:"location_type,omitempty"`
	}{S: s}

	// If private info supplied, update the struct to include those components
	if pi != nil {
		reqBody.E = pi[0].Email
		reqBody.L = pi[0].Loc
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return 0, err
	}

	return addMember(g, reqJSON)
}

// AddMember provides the SensorID interface solution to adding a sensor to a group.
// The function takes an optional PrivateInfo struct, which is necessary only if the
// sensor referenced is a private sensor.
// This call requires a key with write permissions to be set prior to calling.
// The MemberID will be returned on success, or else an error.
func (s SensorID) AddMember(g GroupID, pi ...PrivateInfo) (MemberID, error) {
	reqBody := struct {
		S SensorID `json:"sensor_id"`
		E string   `json:"owner_email,omitempty"`
		L Location `json:"location_type,omitempty"`
	}{S: s}

	// If private info supplied, update the struct to include those components
	if pi != nil {
		reqBody.E = pi[0].Email
		reqBody.L = pi[0].Loc
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return 0, err
	}

	return addMember(g, reqJSON)
}

// addMember is the private function for handling the common code for member addition.
// Both the SensorID and SensorIndex versions of AddMember rely on this.
func addMember(g GroupID, reqJSON []byte) (MemberID, error) {
	u, err := url.Parse(fmt.Sprintf(URLMEMBERS, g))
	if err != nil {
		return 0, err
	}

	resp, err := doRequest(http.MethodPost, u, reqJSON)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return 0, extractError(resp)
	}

	memberResp := struct {
		V string   `json:"api_version"`
		T int      `json:"time_stamp"`
		D int      `json:"data_time_stamp"`
		G GroupID  `json:"group_id"`
		M MemberID `json:"member_id"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&memberResp)
	if err != nil {
		return 0, err
	}

	return memberResp.M, nil
}

// RemoveMember removes the member specified from the group specified.
// This call requires a key with write permissions to be set prior to calling.
// On success, nil will be returned or else an error.
func RemoveMember(m MemberID, g GroupID) error {
	u, err := url.Parse(fmt.Sprintf(URLMEMBERS+"/%d", g, m))
	if err != nil {
		return err
	}

	resp, err := doRequest(http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return extractError(resp)
	}

	return nil
}

// MemberData returns the SensorInfo for a member of a group.
// The SensorParams can restrict the information returned to the named fields.
// This call requires a key with read permissions to be set prior to calling.
// On success, the SensorInfo will be returned, or else an error.
// Note that if a subset of fields is specified, only that data will be returned.
func MemberData(g GroupID, m MemberID, sp SensorParams) (*SensorInfo, error) {
	u, err := url.Parse(fmt.Sprintf(URLMEMBERS+"/%d", g, m))
	if err != nil {
		return nil, err
	}

	return sensorInfo(u, sp)
}

// SensorData returns the SensorInfo for the named SensorIndex.
// The SensorParams can restrict the information returned to the named fields.
// This call requires a key with read permissions to be set prior to calling.
// On success, the SensorInfo will be returned, or else an error.
// Note that if a subset of fields is specified, only that data will be returned.
func SensorData(s SensorIndex, sp SensorParams) (*SensorInfo, error) {
	u, err := url.Parse(fmt.Sprintf(URLSENSORS+"/%d", s))
	if err != nil {
		return nil, err
	}

	return sensorInfo(u, sp)
}

// MembersData returns the information requested for the set (or subset)
// of sensors within the specified Group. The SensorParams must specify
// the elements requested in the "fields" parameter.
// The return value is a map of key/value pairs for each field element
// specified indexed by the sensor_index.
func MembersData(g GroupID, sp SensorParams) (SensorDataSet, error) {
	u, err := url.Parse(fmt.Sprintf(URLMEMBERS, g))
	if err != nil {
		return nil, err
	}

	return sensorsInfo(u, sp)
}

// SensorsData returns the information requested for the set
// of sensors specified by the SensorParam specificiations.
// The SensorParams must specify the elements requested in the "fields" parameter.
// The return value is a map of key/value pairs for each field element
// specified indexed by the sensor_index.
func SensorsData(sp SensorParams) (SensorDataSet, error) {
	u, err := url.Parse(URLSENSORS)
	if err != nil {
		return nil, err
	}

	return sensorsInfo(u, sp)
}

// doRequest creates and executes the http request for the PurpleAir API.
// Depending on the method specified, it appends the appropriate access key required
// as well as setting the content-type. (read key for GET, write key for POST, DELETE)
// It returns the response or an error. When finished processing the response, the
// body must be closed.
func doRequest(m string, u *url.URL, b []byte) (*http.Response, error) {
	req, err := http.NewRequest(m, u.String(), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Add(contentTypeHeader, contentTypeJSON)

	switch m {
	case http.MethodGet:
		if len(apiReadKey) == 0 {
			return nil, errors.New("PurpleAir key not set [read]")
		}
		req.Header.Add(keyHeader, apiReadKey)
	case http.MethodPost, http.MethodDelete:
		if len(apiWriteKey) == 0 {
			return nil, errors.New("PurpleAir key not set [write]")
		}
		req.Header.Add(keyHeader, apiWriteKey)
	default:
		return nil, fmt.Errorf("Unexpected request method [%s]", m)
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// extractError handles an error response back from the API and returns an error
func extractError(r *http.Response) error {
	errorResp := struct {
		V string `json:"api_version"`
		T int    `json:"time_stamp"`
		E string `json:"error"`
		D string `json:"description"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&errorResp)
	if err != nil {
		return err
	}

	// If there is an error response and description, use both. Otherwise just repor the error.
	errMsg := errorResp.E
	if errorResp.D != "" {
		errMsg = fmt.Sprintf("%s: %s", errorResp.E, errorResp.D)
	}

	return errors.New(errMsg)
}
