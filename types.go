package purpleair

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type SensorIndex int // uniquely identifies a sensor within the PurpleAir service
type SensorID string // unique identifier of a sensor found on its label

// SensorFields specify which fields are to be returned for single-sensor calls (MemberData, SensorData).
// If omitted, then all available fields will be returned.
type SensorFields struct {
	Fields string `json:"fields,omitempty"` // comma-delimited list of sensor data fields to return (return all if omitted)
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

const (
	// keyHeader is the HTTP Request header used to pass in the access key value.
	// The value for the keyHeader requires the read key for GET requests and the
	// write key for POST, PUT, or DELETE requests.
	keyHeader string = "X-API-Key"

	// Set the HTTP headers for the content type for JSON.
	contentTypeHeader string = "Content-Type"
	contentTypeJSON   string = "application/json"
)

// PurpleAir API paths
const (
	urlKeys    string = "https://api.purpleair.com/v1/keys"
	urlGroups  string = "https://api.purpleair.com/v1/groups"
	urlMembers string = "https://api.purpleair.com/v1/groups/%d/members"
	urlSensors string = "https://api.purpleair.com/v1/sensors"
)

// KeyTypes as returned from PurpleAir.
// A key can be checked with the CheckAPIKey function.
// Valid read & write keys are required for full API access.
// Retyped string for better type-checking.
type KeyType string

const (
	KeyUnknown       KeyType = "UNKNOWN"
	KeyRead                  = "READ"
	KeyWrite                 = "WRITE"
	KeyReadDisabled          = "READ_DISABLED"
	KeyWriteDisabled         = "WRITE_DISABLED"
)

// GroupMember provides an abstract interface for referring to a sensor either by the SensorIndex,
// which is the reference generated by the PurpleAir service, or the SensorID, which is the
// reference available on the sensor device.
type GroupMember interface {
	AddMember(g GroupID, pi ...PrivateInfo) (MemberID, error)
}

// Private sensors must specify the owner's email and location in order to be accessed.
// Repeated failures to provide correct values for a private sensor may result in access key suspension.
type PrivateInfo struct {
	Email string
	Loc   Location
}

// Unique identifier assigned by PurpleAir to a collection of sensors
// Retyped into for better type-checking.
type GroupID int
type Group struct {
	ID      GroupID   `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

// Custom code for unmarshaling Group structs.
// Converts the raw ID into GroupID and epoch timestamp to Go time.
func (g *Group) UnmarshalJSON(data []byte) error {
	type Shadow Group
	tmp := struct {
		ID      int `json:"id"`
		Created int `json:"created"`
		*Shadow
	}{
		Shadow: (*Shadow)(g),
	}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	g.ID = GroupID(tmp.ID)
	g.Created = time.Unix(int64(tmp.Created), 0)

	return nil
}

// Unique identifier assigned by PurpleAir to a sensor within a Group.
// MemberIDs are valid references only within the specified Group.
// Retyped int for better type-checking.
type MemberID int
type Member struct {
	ID      MemberID    `json:"id"`
	Index   SensorIndex `json:"sensor_index"`
	Created time.Time   `json:"created"`
}

// Custom code for unmashaling Member structs.
// Converts the raw ID into MemberID and epocy timestamp to Go time.
func (m *Member) UnmarshalJSON(data []byte) error {
	type Shadow Member
	tmp := struct {
		ID      int `json:"id"`
		Index   int `json:"sensor_index"`
		Created int `json:"created"`
		*Shadow
	}{
		Shadow: (*Shadow)(m),
	}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	m.ID = MemberID(tmp.ID)
	m.Index = SensorIndex(tmp.Index)
	m.Created = time.Unix(int64(tmp.Created), 0)

	return nil
}

// Sensor location values.
// Retyped int for better type-checking.
type Location int

const (
	LocOutside Location = 0
	LocInside           = 1
)

// Sensor information privacy values.
// Retyped int for better type-checking.
type Privacy int

const (
	SensorPublic  Privacy = 0
	SensorPrivate         = 1
)

// Sensor particulate measurement data channel availability. Many sensors
// have redundant data channels in order to improve the reliability of their readings.
// Retyped int for better type-checking.
type ChannelState int

const (
	ChannelStateNone ChannelState = 0 // no PM sensors detected
	ChannelStateA                 = 1 // PM sensor only on channel A
	ChannelStateB                 = 2 // PM sensor only on channel B
	ChannelStateAll               = 3 // PM sensors on both channels A & B
)

// Sensor data channel status. Sensors may indicate problems with the
// data quality by marking a data channel as downgraded. This may be due
// to defect or transient events (e.g. bug crawling on the sensor)
// Retyped int for better type-checking.
type ChannelFlag int

const (
	ChannelFlagNormal  ChannelFlag = 0 // no sensors downgraded
	ChannelFlagDownA               = 1 // channel A sensors downgraded
	ChannelFlagDownB               = 2 // channel B sensors downgrade
	ChannelFlagDownAll             = 3 // both channel A & B sensors downgrade
)

// Retype the sensor field labels to help enforce typing
type DataField string

// Map of provide sensor query params. In order to avoid misinterpretation
// of Go's default values, only explicit params pertinent for the query
// should be specified. (i.e. If a key isn't relevant for the query, then
// don't include the key in the map.)
type SensorParams map[string]interface{}
type SensorDataRow map[DataField]interface{}
type SensorDataSet map[int]SensorDataRow

const (
	paramFields   string = "fields"
	paramLocation string = "location_type"
	paramReadKey  string = "read_key"
	paramReadKeys string = "read_keys"
	paramShowOnly string = "show_only"
	paramModTime  string = "modified_since"
	paramMaxAge   string = "max_age"
	paramNWLong   string = "nwlng"
	paramNWLat    string = "nwlat"
	paramSELong   string = "selng"
	paramSELat    string = "selat"
)

// Helper types for binding the parameters used for querying sensor info
// to the different types for their values. All implement a "AddParam" interface
// allowing a common mechanism for safely adding query parameters to the call.
type ParamFields struct {
	Value []string
}

type ParamLocation struct {
	Value Location
}

type ParamReadKey struct {
	Value string
}

type ParamReadKeys struct {
	Value []string
}

type ParamShowOnly struct {
	Value []SensorIndex
}

type ParamModTime struct {
	Value time.Time
}

type ParamMaxAge struct {
	Value time.Time
}

type ParamBoundingBox struct {
	NWLong float64
	NWLat  float64
	SELong float64
	SELat  float64
}

// Interface for adding parameters to the SensorParams map.
// Because the SensorParams can take a variety of keys and value types, an interface
// approach allows the necessary flexibility. Implementing a AddParam func for
// each kind of parameter provides better typing control and clarity.
type AddParam interface {
	AddParam(sp SensorParams) SensorParams
}

// Interface implementations for each type of parameter that may be added
// to the SensorParams map. Values may be transformed from a native Go type
// to the necessary formats required by the PurpleAir API. (e.g.
// SensorIndexes transformed to a string of comma-delimited value, or the
// BoundingBox transformed into the four coordinate parameters.
func (p ParamFields) AddParam(sp SensorParams) SensorParams {
	sp[paramFields] = strings.Join(p.Value, ",")

	return sp
}

func (p ParamLocation) AddParam(sp SensorParams) SensorParams {
	sp[paramLocation] = p.Value

	return sp
}

func (p ParamReadKey) AddParam(sp SensorParams) SensorParams {
	sp[paramReadKey] = p.Value

	return sp
}

func (p ParamReadKeys) AddParam(sp SensorParams) SensorParams {
	sp[paramReadKeys] = strings.Join(p.Value, ",")

	return sp
}

func (p ParamShowOnly) AddParam(sp SensorParams) SensorParams {
	var s []string
	for _, i := range p.Value {
		s = append(s, strconv.Itoa(int(i)))
	}
	sp[paramShowOnly] = strings.Join(s, ",")

	return sp
}

func (p ParamModTime) AddParam(sp SensorParams) SensorParams {
	sp[paramModTime] = p.Value.Unix()

	return sp
}

func (p ParamMaxAge) AddParam(sp SensorParams) SensorParams {
	sp[paramMaxAge] = p.Value.Unix()

	return sp
}

func (p ParamBoundingBox) AddParam(sp SensorParams) SensorParams {
	sp[paramNWLong] = p.NWLong
	sp[paramNWLat] = p.NWLat
	sp[paramSELong] = p.SELong
	sp[paramSELat] = p.SELat

	return sp
}

// Collection of averaged statistics for the sensor channel.
// Available as part of the SensorInfo information.
type SensorStats struct {
	PM_2_5        float64   `json:"pm2.5"`
	PM_2_5_10Min  float64   `json:"pm2.5_10minute"`
	PM_2_5_30Min  float64   `json:"pm2.5_30minute"`
	PM_2_5_60Min  float64   `json:"pm2.5_60minute"`
	PM_2_5_6Hour  float64   `json:"pm2.5_6hour"`
	PM_2_5_24Hour float64   `json:"pm2.5_24hour"`
	PM_2_5_1Week  float64   `json:"pm2.5_1week"`
	Timestamp     time.Time `json:"time_stamp"`
}

// Custom code for unmarshaling SensorStats structs.
// Converts the epoch timestamp to Go time.
func (s *SensorStats) UnmarshalJSON(data []byte) error {
	type Shadow SensorStats
	tmp := struct {
		Timestamp int `json:"time_stamp"`
		*Shadow
	}{
		Shadow: (*Shadow)(s),
	}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	s.Timestamp = time.Unix(int64(tmp.Timestamp), 0)

	return nil
}
