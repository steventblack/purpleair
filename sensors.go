package purpleair

import (
	"fmt"
	"net/url"
)

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

	// check for permitted/required params
	for k, _ := range sp {
		switch k {
		case paramFields, paramReadKey:
		default:
			return nil, fmt.Errorf("Unexpected sensor param encountered [%s]", k)
		}
	}

	return paSensor(u, sp)
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
