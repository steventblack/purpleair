package purpleair

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

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

	return sensorInfo(u, sp)
}

// SensorData returns the SensorInfo for the named SensorIndex.
// The SensorParams can restrict the information returned to the named fields.
// This call requires a key with read permissions to be set prior to calling.
// On success, the SensorInfo will be returned, or else an error.
// Note that if a subset of fields is specified, only that data will be returned.
func SensorData(s SensorIndex, sp SensorParams) (*SensorInfo, error) {
	u, err := url.Parse(fmt.Sprintf(urlSensors+"/%d", s))
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
	u, err := url.Parse(fmt.Sprintf(urlMembers, g))
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
	u, err := url.Parse(urlSensors)
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

// paError handles an error response back from the API and returns an error
func paError(r *http.Response) error {
	errorResp := struct {
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
		errMsg = fmt.Sprintf("[%s] %s", errorResp.E, errorResp.D)
	}

	return errors.New(errMsg)
}
