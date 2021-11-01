package purpleair

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

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

// paSensor provides the common code for single-sensor requests.
// Single-sensor calls are supported both by direct reference of the
// SensorIndex or by the MemberID of a Group.
// This function returns a SensorInfo structure with all available fields.
// Not all fields may be filled out or valid depending on the SensorParams
// specified and hardware capabilities.
func paSensor(u *url.URL, sp SensorParams) (*SensorInfo, error) {
	err := paSensorParams(u, sp)
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
		S SensorInfo `json:"sensor"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		return nil, err
	}

	return &payload.S, nil
}

// paSensors provides the common code for multi-sensor requests.
// Multi-sensor calls are supported both by a list of SensorIndex values
// or by the sensors collected in a Group.
// This function returns a SensorDataSet which contains a list of
// the specified fields and their values in a map indexed by the
// the SensorIndex value.
func paSensors(u *url.URL, sp SensorParams) (SensorDataSet, error) {
	err := paSensorParams(u, sp)
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
		G  GroupID         `json:"group_id,omitempty"`
		F  []DataField     `json:"fields,omitempty"`
		L  []string        `json:"location_types,omitempty"`
		CS []string        `json:"channel_states,omitempty"`
		CF []string        `json:"channel_flags,omitempty"`
		D  [][]interface{} `json:"data,omitempty"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		return nil, err
	}

	// Transform the data returned in the payload to a more useful form.
	var sd = make(SensorDataSet)
	for _, r := range payload.D {
		var row = make(SensorDataRow)

		// Fill out a data row with key/value pairs for each field element
		// the key name is found in the matching location of the fields (F) list
		// For selected values, translate the numerical value returned to the
		// appropriate label
		for i, v := range r {
			switch k := payload.F[i]; k {
			case "location_type":
				row[k] = payload.L[int(v.(float64))]
			case "channel_states":
				row[k] = payload.CS[int(v.(float64))]
			case "channel_flags":
				row[k] = payload.CF[int(v.(float64))]
			default:
				row[k] = v
			}
		}

		// Identify the SensorIndex for the data row and assign the row
		// to the data set referenced by the index value.
		// If no SensorIndex found, there's a big problem.
		if si, ok := row["sensor_index"]; ok {
			sd[int(si.(float64))] = row
		} else {
			return nil, errors.New("Required element not found [sensor_index]")
		}
	}

	return sd, nil
}

// paSensorParams processes the parameters passed in for sensor information
// calls and converts them into url query parameters (properly encoded).
// This call is used by all single and multi-sensor information calls
// although the permitted parameters vary by call. Each call should
// filter the parameters and pass only legal ones before passing in
// the SensorParams to the common code.
func paSensorParams(u *url.URL, sp SensorParams) error {
	q := u.Query()

	for k, v := range sp {
		switch k {
		case SensorParamFields, SensorParamReadKeys, SensorParamShowOnly:
			q.Add(string(k), fmt.Sprintf("%s", v))
		case SensorParamLocation, SensorParamModTime, SensorParamMaxAge:
			q.Add(string(k), fmt.Sprintf("%d", v))
		case SensorParamNWLong, SensorParamNWLat, SensorParamSELong, SensorParamSELat:
			q.Add(string(k), fmt.Sprintf("%f", v))
		default:
			return fmt.Errorf("Unexpected sensor param specified [%s]", k)
		}
	}

	u.RawQuery = q.Encode()

	return nil
}
