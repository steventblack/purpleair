package purpleair

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// Retype the sensor field labels to help enforce typing
type DataField string

// Helper function returning list of all available sensor fields for
// accessing sensor information. Wrapping in this manner to provide
// better const control given Go doesn't support const slices.
func DataFields() []DataField {
	return []DataField{
		// Station information and status fields:
		"name", "icon", "model", "hardware", "location_type", "private", "latitude",
		"longitude", "altitude", "position_rating", "led_brightness", "firmware_version",
		"firmware_upgrade", "rssi", "uptime", "pa_latency", "memory", "last_seen",
		"last_modified", "date_created", "channel_state", "channel_flags", "channel_flags_manual",
		"channel_flags_auto", "confidence", "confidence_manual", "confidence_auto",
		// Environmental fields:
		"humidity", "humidity_a", "humidity_b", "temperature", "temperature_a",
		"temperature_b", "pressure", "pressure_a", "pressure_b",
		// Miscellaneous fields:
		"voc", "voc_a", "voc_b", "ozone1", "analog_input",
		// PM 1.0 fields:
		"pm1.0", "pm1.0_a", "pm1.0_b", "pm1.0_atm", "pm1.0_atm_a", "pm1.0_atm_b",
		"pm1.0_cf_1", "pm1.0_cf_1_a", "pm1.0_cf_1_b",
		// PM 2.5 fields:
		"pm2.5_alt", "pm2.5_alt_a", "pm2.5_alt_b", "pm2.5", "pm2.5_a", "pm2.5_b",
		"pm2.5_atm", "pm2.5_atm_a", "pm2.5_atm_b", "pm2.5_cf_1", "pm2.5_cf_1_a",
		"pm2.5_cf_1_b",
		// PM 2.5 pseudo average fields:
		"pm2.5_10minute", "pm2.5_10minute_a", "pm2.5_10minute_b",
		"pm2.5_30minute", "pm2.5_30minute_a", "pm2.5_30minute_b",
		"pm2.5_60minute", "pm2.5_60minute_a", "pm2.5_60minute_b",
		"pm2.5_6hour", "pm2.5_6hour_a", "pm2.5_6hour_b",
		"pm2.5_24hour", "pm2.5_24hour_a", "pm2.5_24hour_b",
		"pm2.5_1week", "pm2.5_1week_a", "pm2.5_1week_b",
		// PM 10.0 fields:
		"pm10.0", "pm10.0_a", "pm10.0_b", "pm10.0_atm", "pm10.0_atm_a", "pm10.0_atm_b",
		"pm10.0_cf_1", "pm10.0_cf_1_a", "pm10.0_cf_1_b",
		// Particle count fields:
		"0.3_um_count", "0.3_um_count_a", "0.3_um_count_b",
		"0.5_um_count", "0.5_um_count_a", "0.5_um_count_b",
		"1.0_um_count", "1.0_um_count_a", "1.0_um_count_b",
		"2.5_um_count", "2.5_um_count_a", "2.5_um_count_b",
		"5.0_um_count", "5.0_um_count_a", "5.0_um_count_b",
		"10.0_um_count", "10.0_um_count_a", "10.0_um_count_b",
		// ThingSpeak fields:
		"primary_id_a", "primary_key_a", "secondary_id_a", "secondary_key_a",
		"primary_id_b", "primary_key_b", "secondary_id_b", "secondary_key_b"}
}

// Retype the sensor query param to help enforce typing
type SensorParam string

// Defined list of sensor query keys. Each key expects an
// appropriately typed value.
const (
	SP_FIELDS   SensorParam = "fields"
	SP_LOCATION SensorParam = "location_type"
	SP_READKEYS SensorParam = "read_keys"
	SP_SHOWONLY SensorParam = "show_only"
	SP_MODTIME  SensorParam = "modified_since"
	SP_MAXAGE   SensorParam = "max_age"
	SP_NWLNG    SensorParam = "nwlng"
	SP_NWLAT    SensorParam = "nwlat"
	SP_SELNG    SensorParam = "selng"
	SP_SELAT    SensorParam = "selat"
)

// Map of provide sensor query params. In order to avoid misinterpretation
// of Go's default values, only explicit params pertinent for the query
// should be specified. (i.e. If a key isn't relevant for the query, then
// don't include the key in the map.)
type SensorParams map[SensorParam]interface{}
type SensorDataRow map[DataField]interface{}
type SensorDataSet map[int]SensorDataRow

// Common code for multi-sensor data collection.
// Applies to members of groups or other list of sensors
func sensorsInfo(url string, sp SensorParams) (SensorDataSet, error) {
	req, err := setupCall(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, extractError(resp)
	}

	sensorsResp := struct {
		V      string          `json:"api_version,omitempty"`
		T      int             `json:"time_stamp,omitempty"`
		D      int             `json:"data_time_stamp,omitempty"`
		G      GroupID         `json:"group_id,omitempty"` // only present for group member queries
		A      int             `json:"max_age,omitempty"`
		F      string          `json:"firmware_default_version,omitempty"`
		Fields []DataField     `json:"fields,omitempty"`
		Locs   []string        `json:"location_types,omitempty"`
		States []string        `json:"channel_states,omitempty"`
		Flags  []string        `json:"channel_flags,omitempty"`
		Data   [][]interface{} `json:"data,omitempty"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&sensorsResp)
	if err != nil {
		return nil, err
	}

	// Translate the response into a more useful map indexed by the
	// sensor_index matched with the set of key/value pairs for its data.
	// The keys for the data (and selected values for location, channel states, etc.)
	// are extracted from the relevant position in the appropriate response fields.
	var data = make(SensorDataSet)
	for _, r := range sensorsResp.Data {
		var row = make(SensorDataRow)

		for j, v := range r {
			switch k := sensorsResp.Fields[j]; k {
			case "location_type":
				row[k] = sensorsResp.Locs[int(v.(float64))]
			case "channel_states":
				row[k] = sensorsResp.States[int(v.(float64))]
			case "channel_flags":
				row[k] = sensorsResp.Flags[int(v.(float64))]
			default:
				row[k] = v
			}
		}

		if si, ok := row["sensor_index"]; ok {
			data[int(si.(float64))] = row
		} else {
			return nil, errors.New("Required element not found [sensor_index]")
		}

		log.Println()
	}

	// debug print
	for k, v := range data {
		log.Printf("sensor_index[%d]: %+v\n", k, v)
	}

	return data, nil
}

// Process the provided params and convert to a query string format.
// Note: because the PurpleAir API specifies the calls for GetMembersData
// and GetSensorsData as GET requests, passing in the params as part of the
// body is not an option. This converts the params to be properly encoded
// as part of the query string.
func addSensorParams(u *url.URL, sp SensorParams) error {
	fieldsPresent := false
	q := u.Query()

	for k, v := range sp {
		switch k {
		case SP_FIELDS:
			fieldsPresent = true
			fallthrough
		case SP_READKEYS, SP_SHOWONLY:
			q.Add(string(k), fmt.Sprintf("%s", v))
		case SP_LOCATION, SP_MODTIME, SP_MAXAGE:
			q.Add(string(k), fmt.Sprintf("%d", v))
		case SP_NWLNG, SP_NWLAT, SP_SELNG, SP_SELAT:
			q.Add(string(k), fmt.Sprintf("%f", v))
		default:
			return fmt.Errorf("Unexpected sensor param specified [%s]", k)
		}
	}

	if fieldsPresent != true {
		return errors.New("Required parameter not found [fields]")
	}

	u.RawQuery = q.Encode()

	return nil
}
