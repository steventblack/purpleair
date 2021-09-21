package purpleair

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// Global variables for retaining the API access keys.
// These are set via the SetAPIKey.
var (
	apiReadKey  string
	apiWriteKey string
)

type KeyType string   // maps the response from PurpleAir when checking the validity and permissions
type SensorIndex int  // uniquely identifies a sensor within the PurpleAir service
type SensorID string  // unique identifier of a sensor found on its label
type GroupID int      // unique identifier of a collection of sensors within the PurpleAir service
type MemberID int     // unique identifier of a sensor within a specific group defined in the PurpleAir service
type Location int     // enables typechecking on defined location values
type Privacy int      // setting for a sensor indicating public or private
type ChannelState int // States for the sensor data channel availability
type ChannelFlag int  // Flags for the sensor data channels

type Group struct {
	ID         GroupID `json:"id"`
	Name       string  `json:"name"`
	CreatedUTC int     `json:"created"`
}

type Member struct {
	ID         MemberID    `json:"id"`
	Index      SensorIndex `json:"sensor_index"`
	CreatedUTC int         `json:"created"`
}

// Private sensors must specify the owner's email and location in order to be accessed.
// Repeated failures to provide correct values for a private sensor may result in access key suspension.
type PrivateInfo struct {
	Email string
	Loc   Location
}

// GroupMember provides an abstract interface for referring to a sensor either by the SensorIndex,
// which is the reference generated by the PurpleAir service, or the SensorID, which is the
// reference available on the sensor device. Several calls provide flexibility on accepting
// either a SensorIndex or SensorID, so providing an abstraction helps reduce redundancy.
type GroupMember interface {
	AddMember(g GroupID, pi ...PrivateInfo) (MemberID, error)
}

// SensorParams are options that can be passed in for customizing the sensor information
// They are all optional. For calls referencing a single sensor, only the Fields parameter is considered.
type SensorParams struct {
	Fields   string   `json:"fields,omitempty"`         // which sensor data fields to return
	Loc      Location `json:"location_type,omitempty"`  // location: inside/outside
	ReadKeys string   `json:"read_keys,omitempty"`      // key required for access to private devices
	Show     string   `json:"show_only,omitempty"`      // return data only for sensorIndexes listed
	Mod      int      `json:"modified_since,omitempty"` // return data only if updated since timestamp
	MaxAge   int      `json:"max_age,omitempty"`        // return data only if updated within specifed seconds
	LngNW    float64  `json:"nwlng,omitempty"`          // bounding box: provide NW and SE coordinates
	LatNW    float64  `json:"nwlat,omitempty"`          // and return sensor info only for devices within the box
	LngSE    float64  `json:"selng,omitempty"`
	LatSE    float64  `json:"selat,omitempty"`
}

// Collection of averaged statistics for the sensor channel
type SensorStats struct {
	PM_2_5        float64 `json:"pm2.5"`
	PM_2_5_10Min  float64 `json:"pm2.5_10minute"`
	PM_2_5_30Min  float64 `json:"pm2.5_30minute"`
	PM_2_5_60Min  float64 `json:"pm2.5_60minute"`
	PM_2_5_6Hour  float64 `json:"pm2.5_6hour"`
	PM_2_5_24Hour float64 `json:"pm2.5_24hour"`
	PM_2_5_1Week  float64 `json:"pm2.5_1week"`
	Timestamp     int     `json:"time_stamp"`
}

// SensorInfo is the data response to a sensor query.
// Not all fields may be available depending on the query fields specified or hardware capabilities.
type SensorInfo struct {
	Index           SensorIndex  `json:"sensor_index,omitempty"`
	Icon            int          `json"icon,omitempty"`
	Name            string       `json:"name,omitempty"`
	Private         Privacy      `json:"private,omitempty"`
	Loc             Location     `json:"location_type,omitempty"`
	Lat             float64      `json:"latitude,omitempty"`
	Lng             float64      `json:"longitude,omitempty"`
	Alt             int          `json:"altitude,omitempty"`
	Pos             int          `json:"position_rating,omitempty"`
	Model           string       `json:"model,omitempty"`
	Hardware        string       `json:"hardware,omitempty"`
	FirmVersion     string       `json:"firmware_version,omitempty"`
	FirmUpgrade     string       `json:"firmware_upgrade,omitempty"`
	RSSI            int          `json:"rssi,omitempty"`
	Uptime          int          `json:"uptime,omitempty"`
	Latency         int          `json:"pa_latency,omitempty"`
	Memory          int          `json:"memory,omitempty"`
	LED             int          `json:"led_brightness,omitempty"`
	ChnlState       ChannelState `json:"channel_state,omitempty"`
	ChnlFlags       ChannelFlag  `json:"channel_flags,omitempty"`
	ChnlManual      ChannelFlag  `json:"channel_flags_manual,omitempty"`
	ChnlAuto        ChannelFlag  `json:"channel_flags_auto,omitempty"`
	Cfdnc           int          `json:"confidence,omitempty"`
	CfdncManual     int          `json:"confidence_manual,omitempty"`
	CfdncAuto       int          `json:"confidence_auto,omitempty"`
	Mod             int          `json:"last_modifed,omitempty"`
	Created         int          `json:"date_created,omitempty"`
	PM_1_0          float64      `json:"pm1.0,omitempty"`
	PM_1_0_A        float64      `json:"pm1.0_a,omitempty"`
	PM_1_0_B        float64      `json:"pm1.0_b,omitempty"`
	PM_1_0_Atm      float64      `json:"pm1.0_atm,omitempty"`
	PM_1_0_Atm_A    float64      `json:"pm1.0_atm_a,omitempty"`
	PM_1_0_Atm_B    float64      `json:"pm1.0_atm_b,omitempty"`
	PM_1_0_Cf_1     float64      `json:"pm1.0_cf_1,omitempty"`
	PM_1_0_Cf_1_A   float64      `json:"pm1.0_cf_1_a,omitempty"`
	PM_1_0_Cf_1_B   float64      `json:"pm1.0_cf_1_b,omitempty"`
	PM_2_5_Alt      float64      `json:"pm2.5_alt,omitempty"`
	PM_2_5_Alt_A    float64      `json:"pm2.5_alt_a,omitempty"`
	PM_2_5_Alt_B    float64      `json:"pm2.5_alt_b,omitempty"`
	PM_2_5          float64      `json:"pm2.5,omitempty"`
	PM_2_5_A        float64      `json:"pm2.5_a,omitempty"`
	PM_2_5_B        float64      `json:"pm2.5_b,omitempty"`
	PM_2_5_Atm      float64      `json:"pm2.5_atm,omitempty"`
	PM_2_5_Atm_A    float64      `json:"pm2.5_atm_a,omitempty"`
	PM_2_5_Atm_B    float64      `json:"pm2.5_atm_b,omitempty"`
	PM_2_5_Cf_1     float64      `json:"pm2.5_cf_1,omitempty"`
	PM_2_5_Cf_1_A   float64      `json:"pm2.5_cf_1_a,omitempty"`
	PM_2_5_Cf_1_B   float64      `json:"pm2.5_cf_1_b,omitempty"`
	PM_2_5_10Min    float64      `json:"pm2.5_10minute,omitempty"`
	PM_2_5_10Min_A  float64      `json:"pm2.5_10minute_a,omitempty"`
	PM_2_5_10Min_B  float64      `json:"pm2.5_10minute_b,omitempty"`
	PM_2_5_30Min    float64      `json:"pm2.5_30minute,omitempty"`
	PM_2_5_30Min_A  float64      `json:"pm2.5_30minute_a,omitempty"`
	PM_2_5_30Min_B  float64      `json:"pm2.5_30minute_b,omitempty"`
	PM_2_5_60Min    float64      `json:"pm2.5_60minute,omitempty"`
	PM_2_5_60Min_A  float64      `json:"pm2.5_60minute_a,omitempty"`
	PM_2_5_60Min_B  float64      `json:"pm2.5_60minute_b,omitempty"`
	PM_2_5_6Hour    float64      `json:"pm2.5_6hour,omitempty"`
	PM_2_5_6Hour_A  float64      `json:"pm2.5_6hour_a,omitempty"`
	PM_2_5_6Hour_B  float64      `json:"pm2.5_6hour_b,omitempty"`
	PM_2_5_24Hour   float64      `json:"pm2.5_24hour,omitempty"`
	PM_2_5_24Hour_A float64      `json:"pm2.5_24hour_a,omitempty"`
	PM_2_5_24Hour_B float64      `json:"pm2.5_24hour_b,omitempty"`
	PM_2_5_1Week    float64      `json:"pm2.5_1week,omitempty"`
	PM_2_5_1Week_A  float64      `json:"pm2.5_1week_a,omitempty"`
	PM_2_5_1Week_B  float64      `json:"pm2.5_1week_b,omitempty"`
	PM_10_0         float64      `json:"pm10.0,omitempty"`
	PM_10_0_A       float64      `json:"pm10.0_a,omitempty"`
	PM_10_0_B       float64      `json:"pm10.0_b,omitempty"`
	PM_10_0_Atm     float64      `json:"pm10.0_atm,omitempty"`
	PM_10_0_Atm_A   float64      `json:"pm10.0_atm_a,omitempty"`
	PM_10_0_Atm_B   float64      `json:"pm10.0_atm_b,omitempty"`
	PM_10_0_Cf_1    float64      `json:"pm10.0_cf_1,omitempty"`
	PM_10_0_Cf_1_A  float64      `json:"pm10.0_cf_1_a,omitempty"`
	PM_10_0_Cf_1_B  float64      `json:"pm10.0_cf_1_b,omitempty"`
	PC_0_3um        int          `json:"0.3_um_count,omitempty"`
	PC_0_3um_A      int          `json:"0.3_um_count_a,omitempty"`
	PC_0_3um_B      int          `json:"0.3_um_count_b,omitempty"`
	PC_0_5um        int          `json:"0.5_um_count,omitempty"`
	PC_0_5um_A      int          `json:"0.5_um_count_a,omitempty"`
	PC_0_5um_B      int          `json:"0.5_um_count_b,omitempty"`
	PC_1_0um        int          `json:"1.0_um_count,omitempty"`
	PC_1_0um_A      int          `json:"1.0_um_count_a,omitempty"`
	PC_1_0um_B      int          `json:"1.0_um_count_b,omitempty"`
	PC_2_5um        int          `json:"2.5_um_count,omitempty"`
	PC_2_5um_A      int          `json:"2.5_um_count_a,omitempty"`
	PC_2_5um_B      int          `json:"2.5_um_count_b,omitempty"`
	PC_5_0um        int          `json:"5.0_um_count,omitempty"`
	PC_5_0um_A      int          `json:"5.0_um_count_a,omitempty"`
	PC_5_0um_B      int          `json:"5.0_um_count_b,omitempty"`
	PC_10_0um       int          `json:"10.0_um_count,omitempty"`
	PC_10_0um_A     int          `json:"10.0_um_count_a,omitempty"`
	PC_10_0um_B     int          `json:"10.0_um_count_b,omitempty"`
	Stats           SensorStats  `json:"stats,omitempty"`
	Stats_A         SensorStats  `json:"stats_a,omitempty"`
	Stats_B         SensorStats  `json:"stats_b,omitempty"`
	Humidity        int          `json:"humidity,omitempty"`
	Humidity_A      int          `json:"humidity_a,omitempty"`
	Humidity_B      int          `json:"humidity_b,omitempty"`
	Temp            int          `json:"temperature,omitempty"`
	Temp_A          int          `json:"temperature_a,omitempty"`
	Temp_B          int          `json:"temperature_b,omitempty"`
	Pressure        float64      `json:"pressure,omitempty"`
	Pressure_A      float64      `json:"pressure_a,omitempty"`
	Pressure_B      float64      `json:"pressure_b,omitempty"`
	VOC             float64      `json:"voc,omitempty"`
	VOC_A           float64      `json:"voc_a,omitempty"`
	VOC_B           float64      `json:"voc_b,omitempty"`
	Ozone           float64      `json:"ozone1,omitempty"`
	AnalogIn        float64      `json:"analog_input,omitempty"`
	PrimaryID_A     int          `json:"primary_id_a,omitempty"`
	PrimaryKey_A    string       `json:"primary_key_a,omitempty"`
	SecondaryID_A   int          `json:"secondary_id_a,omitempty"`
	SecondaryKey_A  string       `json:"secondary_key_a,omitempty"`
	PrimaryID_B     int          `json:"primary_id_b,omitempty"`
	PrimaryKey_B    string       `json:"primary_key_b,omitempty"`
	SecondaryID_B   int          `json:"secondary_id_b,omitempty"`
	SecondaryKey_B  string       `json:"secondary_key_b,omitempty"`
}

// KeyTypes as returned from PurpleAir.
const (
	APIKEYUNKNOWN       KeyType = "UNKNOWN"
	APIKEYREAD          KeyType = "READ"
	APIKEYWRITE         KeyType = "WRITE"
	APIKEYREADDISABLED  KeyType = "READ_DISABLED"
	APIKEYWRITEDISABLED KeyType = "WRITE_DISABLED"
)

// Defined location values
const (
	OUTSIDE Location = 0
	INSIDE  Location = 1
)

// Defined privacy values
const (
	PUBLIC  Privacy = 0
	PRIVATE Privacy = 1
)

// Defined channel states
const (
	PM_NONE ChannelState = 0 // no PM sensors detected
	PM_A    ChannelState = 1 // PM sensor only on channel A
	PM_B    ChannelState = 2 // PM sensor only on channel B
	PM_ALL  ChannelState = 3 // PM sensors on both channel A & B
)

// Defined channel flags
const (
	NORMAL         ChannelFlag = 0 // no sensors marked as downgraded
	DOWNGRADED_A   ChannelFlag = 1 // channel A sensors downgraded
	DOWNGRADED_B   ChannelFlag = 2 // channel B sensors downgraded
	DOWNGRADED_ALL ChannelFlag = 3 // both channel A & B sensors downgraded
)

// PurpleAir API paths
const (
	URLKEYS    string = "https://api.purpleair.com/v1/keys"
	URLGROUPS  string = "https://api.purpleair.com/v1/groups"
	URLMEMBERS string = "https://api.purpleair.com/v1/groups/%d/members"
	URLSENSORS string = "https://api.purpleair.com/v1/sensors"
)

// APIKEYHEADER is the HTTP Request header used to pass in the access key value.
const APIKEYHEADER string = "X-API-Key"

/*
type Field string

// Fields for sensor information. Not all fields may be available on all devices.
// Sensor information can selectively return data for named fields, or all data if omitted.
const (
	Fields []Field = {"name", "icon", "model", "hardware", "location_type", "private",
		"latitude", "longitude", "altitude", "position_rating", "led_brightness",
		"firmware_version", "firmware_upgrade", "rssi", "uptime", "pa_latency",
		"memory", "last_seen", "last_modified", "date_created", "channel_state",
		"channel_flags", "channel_flags_manual", "channel_flags_auto", "confidence",
		"confidence_manual", "confidence_auto"}
)
*/

// SetAPIKey checks the validity and permissions of the provided access key string.
// If the key is valid, it will be retained by the module for further calls.
// Only one key for each permission (read or write) will be retained, and additional
// calls with other valid keys will result in replacement.
// The KeyType will be returned on success, or an error on failure.
func SetAPIKey(k string) (KeyType, error) {
	keyType, err := CheckAPIKey(k)
	if err != nil {
		log.Printf("Unable to CheckAPIKey: %v\n", err)
		return APIKEYUNKNOWN, err
	}

	if keyType == APIKEYREAD {
		log.Printf("Successfully set API read key\n")
		apiReadKey = k
	} else if keyType == APIKEYWRITE {
		log.Printf("Successfully set API write key\n")
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

	keyType := APIKEYUNKNOWN

	req, err := http.NewRequest(http.MethodGet, URLKEYS, nil)
	if err != nil {
		log.Printf("Unable to create HTTP request: %s\n", err)
		return keyType, err
	}
	req.Header.Add(APIKEYHEADER, k)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Unable to execute HTTP request: %s\n", err)
		return keyType, err
	}
	defer resp.Body.Close()

	// if invalid key, an error is returned
	if resp.StatusCode != http.StatusCreated {
		errorResp := struct {
			E string `json:"error"`
		}{}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&errorResp)
		if err != nil {
			log.Printf("Unable to decode HTTP body: %s\n", err)
			return keyType, err
		}

		return keyType, errors.New(errorResp.E)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&keyTypeResp)
	if err != nil {
		log.Printf("Unable to decode HTTP body: %s\n", err)
		return keyType, err
	}

	keyType = keyTypeResp.K
	log.Printf("Extracted key type: %s\n", keyType)

	return keyType, nil
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
		log.Printf("Unable to marshal json body: %s\n", err)
		return 0, err
	}

	req, err := setupCall(http.MethodPost, URLGROUPS, reqJSON)
	if err != nil {
		log.Printf("Unable to setup call: %s\n", err)
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Unable to execute HTTP request: %s\n", err)
		return 0, err
	}
	defer resp.Body.Close()

	groupResp := struct {
		G GroupID `json:"group_id"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&groupResp)
	if err != nil {
		log.Printf("Unable to decode HTTP body: %s\n", err)
		return 0, err
	}

	return groupResp.G, nil
}

// DeleteGroup deletes the persistent reference of a sensor group on the PurpleAir service.
// All sensor members must be removed prior to group deletion.
// This call requires a key with write permissions to be set prior to calling.
// An error will be returned on failure, or else nil
func DeleteGroup(g GroupID) error {
	url := fmt.Sprintf("%s/%d", URLGROUPS, g)
	req, err := setupCall(http.MethodDelete, url, nil)
	if err != nil {
		log.Printf("Unable to setup API call: %s\n", err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Unable to execute HTTP request: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		groupResp := struct {
			E string `json:"error"`
		}{}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&groupResp)
		if err != nil {
			log.Printf("Unable to decode HTTP response: %s\n", err)
			return err
		}

		return errors.New(groupResp.E)
	}

	return nil
}

// ListGroups provides all groups defined in the PurpleAir service associated with the access key.
// This call requires a key with read permissions to be set prior to calling.
// The list of groups will be returned on success, or else an error.
func ListGroups() ([]Group, error) {
	req, err := setupCall(http.MethodGet, URLGROUPS, nil)
	if err != nil {
		log.Printf("Unable to setup API call: %s\n", err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Unable to execute HTTP request: %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	groupResp := struct {
		Groups []Group `json:"groups"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&groupResp)
	if err != nil {
		log.Printf("Unable to decode HTTP body: %s\n", err)
		return nil, err
	}

	return groupResp.Groups, nil
}

// GroupDetails provides the list of member sensors defined for the specified group.
// This call requires a key with read permissions to be set prior to calling.
// The list of members will be returned on success, or else an error.
func GroupDetails(g GroupID) ([]Member, error) {
	url := fmt.Sprintf("%s/%d", URLGROUPS, g)
	req, err := setupCall(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Unable to setup API call: %s\n", err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Unable to execute HTTP request: %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	memberResp := struct {
		Members []Member `json:"members"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&memberResp)
	if err != nil {
		log.Printf("Unable to decode HTTP body: %s\n", err)
		return nil, err
	}

	return memberResp.Members, nil
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
		log.Printf("Unable to marshal json body: %s\n", err)
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
		log.Printf("Unable to marshal json body: %s\n", err)
		return 0, err
	}

	return addMember(g, reqJSON)
}

// addMember is the private function for handling the common code for member addition.
// Both the SensorID and SensorIndex versions of AddMember rely on this.
func addMember(g GroupID, reqJSON []byte) (MemberID, error) {
	url := fmt.Sprintf(URLMEMBERS, g)
	req, err := setupCall(http.MethodPost, url, reqJSON)
	if err != nil {
		log.Printf("Unable to setup call: %s\n", err)
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Unable to execute HTTP request: %s\n", err)
		return 0, err
	}
	defer resp.Body.Close()

	memberResp := struct {
		M MemberID `json:"member_id"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&memberResp)
	if err != nil {
		log.Printf("Unable to decode HTTP body: %s\n", err)
		return 0, err
	}

	return memberResp.M, nil
}

// RemoveMember removes the member specified from the group specified.
// This call requires a key with write permissions to be set prior to calling.
// On success, nil will be returned or else an error.
func RemoveMember(m MemberID, g GroupID) error {
	url := fmt.Sprintf(URLMEMBERS+"/%d", g, m)
	req, err := setupCall(http.MethodDelete, url, nil)
	if err != nil {
		log.Printf("Unable to setup API call: %s\n", err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Unable to execute HTTP request: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		groupResp := struct {
			E string `json:"error"`
		}{}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&groupResp)
		if err != nil {
			log.Printf("Unable to decode HTTP response: %s\n", err)
			return err
		}

		return errors.New(groupResp.E)
	}

	return nil
}

// MemberData returns the SensorInfo for a member of a group.
// The optional SensorFields parameter can restrict the information returned to the named fields.
// Omitting the SensorFields parameter will return all available information fields.
// This call requires a key with read permissions to be set prior to calling.
// On success, the SensorInfo will be returned, or else an error.
func MemberData(g GroupID, m MemberID, p ...SensorParams) (*SensorInfo, error) {
	reqJSON, err := processInfoParams(p)
	url := fmt.Sprintf(URLMEMBERS+"/%d", g, m)

	req, err := setupCall(http.MethodGet, url, reqJSON)
	if err != nil {
		log.Printf("Unable to setup API call: %s\n", err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Unable to execute HTTP request: %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	sensorResp := struct {
		// additional fields in the response are omitted as they aren't of any interest
		S SensorInfo `json:"sensor"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&sensorResp)
	if err != nil {
		log.Printf("Unable to decode HTTP body: %s\n", err)
		return nil, err
	}

	return &sensorResp.S, nil
}

// MembersData returns the SensorInfo for all members of a group.
// The required SensorParams must specify at least the fields option,
// and may additionally specify other selection criteria
// This call requires a key with read permissions to be set prior to calling.
// On success, a slice of SensorInfos will be returned, or else an error.
func MembersData(g GroupID, p SensorParams) ([]SensorInfo, error) {
	/*
		reqJSON, err := processInfoParams(p)
		url := fmt.Sprintf(URLMEMBERS, g)

		req, err := setupCall(http.MethodGet, url, reqJSON)
		if err != nil {
			log.Printf("Unable to setup API call: %s\n", err)
			return nil, err
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Unable to execute HTTP request: %s\n", err)
			return nil, err
		}
		defer resp.Body.Close()

		sensorResp := struct {
			// additional fields in the response are omitted as they aren't of any interest
			F []string   `json:"fields"`
			D [][]string `json:"data"`
		}{}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&sensorResp)
		if err != nil {
			log.Printf("Unable to decode HTTP body: %s\n", err)
			return nil, err
		}

		return &sensorResp.S, nil
	*/
	return nil, nil
}

// setupCall performs common tasks that are prerequisite before calling the API.
// It initializes a request object and adds the appropriate key (read or write) to the request.
// It returns a request ready for execution or an error.
func setupCall(method string, url string, reqBody []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	switch method {
	case "GET":
		if len(apiReadKey) == 0 {
			return nil, errors.New("PurpleAir read key is not set")
		}
		req.Header.Add(APIKEYHEADER, apiReadKey)
	case "POST":
		fallthrough
	case "DELETE":
		if len(apiWriteKey) == 0 {
			return nil, errors.New("PurpleAir write key is not set")
		}
		req.Header.Add(APIKEYHEADER, apiWriteKey)
	}
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

// processInfoParams converts the params (if specified) into the appropriate JSON
// for a SensorInfo request. If nothing specified, then a nil byte array is returned.
func processInfoParams(p []SensorParams) ([]byte, error) {
	switch len(p) {
	case 0:
		return nil, nil

	case 1:
		jsonBody, err := json.Marshal(p[0])
		if err != nil {
			log.Printf("Unable to marshal json body: %s\n", err)
			return nil, err
		}
		return jsonBody, nil

	default:
		return nil, fmt.Errorf("Too many SensorParams specified (%d)", len(p))
	}
}
