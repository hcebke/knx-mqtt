package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	localdpt "github.com/pakerfeldt/knx-mqtt/internal/dpt"
	"github.com/vapourismo/knx-go/knx/dpt"
)

var regexpGad = regexp.MustCompile(`^\d+\/\d+\/\d+$`)
var regexpGadOrFlat = regexp.MustCompile(`^\d+\/\d+\/\d+|\d+$`)
var regexpFlatGad = regexp.MustCompile(`^\d+$`)

var regexpTimeOfDay = regexp.MustCompile(`^(Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday)?\s*(\d{2}):(\d{2}):(\d{2})$`)
var regexpDate = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})$`)
var regexpDateTime = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})(?:\s+(Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday))?\s+(\d{2}):(\d{2}):(\d{2})(?:\s+\(Summer Time\))?(?:\s+\[(.*)\])?$`)
var regexpRgb = regexp.MustCompile(`^#([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})$`)
var regexpXyy = regexp.MustCompile(`^x:\s*(\d+)\s*y:\s*(\d+)\s*Y:\s*(\d+)\s*ColorValid:\s*(true|false),\s*BrightnessValid:\s*(true|false)$`)
var regexpRgbw = regexp.MustCompile(`^Red:\s*(\d+)\s*Green:\s*(\d+)\s*Blue:\s*(\d+)\s*White:\s*(\d+)\s*RedValid:\s*(true|false),\s*GreenValid:\s*(true|false),\s*BlueValid:\s*(true|false),\s*WhiteValid:\s*(true|false)$`)

var weekdays = map[string]uint8{
	"Monday": 1, "Tuesday": 2, "Wednesday": 3,
	"Thursday": 4, "Friday": 5, "Saturday": 6, "Sunday": 7,
}

func IsRegularGroupAddress(address string) bool {
	return regexpGad.MatchString(address)
}

func IsFlatGroupAddress(address string) bool {
	return regexpFlatGad.MatchString(address)
}

func IsRegularOrFlatGroupAddress(address string) bool {
	return regexpGadOrFlat.MatchString(address)
}

func PackString(datatype string, value string) ([]byte, error) {
	if packFunc, exists := boolPackFunctions[datatype]; exists {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("could not convert to boolean: %s", value)
		}
		return packFunc(boolValue), nil
	}

	if packFunc, exists := float32PackFunctions[datatype]; exists {
		float32Value, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, fmt.Errorf("could not convert to float32: %s", value)
		}
		return packFunc(float32(float32Value)), nil
	}

	if packFunc, exists := uint8PackFunctions[datatype]; exists {
		uint8Value, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("could not convert to uint8: %s", value)
		}
		return packFunc(uint8(uint8Value)), nil
	}

	if packFunc, exists := int8PackFunctions[datatype]; exists {
		int8Value, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("could not convert to int8: %s", value)
		}
		return packFunc(int8(int8Value)), nil
	}

	if packFunc, exists := uint16PackFunctions[datatype]; exists {
		uint16Value, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("could not convert to uint16: %s", value)
		}
		return packFunc(uint16(uint16Value)), nil
	}

	if packFunc, exists := int16PackFunctions[datatype]; exists {
		int16Value, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("could not convert to int16: %s", value)
		}
		return packFunc(int16(int16Value)), nil
	}

	if packFunc, exists := uint32PackFunctions[datatype]; exists {
		uint32Value, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("could not convert to uint32: %s", value)
		}
		return packFunc(uint32(uint32Value)), nil
	}

	if packFunc, exists := int32PackFunctions[datatype]; exists {
		int32Value, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("could not convert to uint32: %s", value)
		}
		return packFunc(int32(int32Value)), nil
	}

	if packFunc, exists := stringPackFunctions[datatype]; exists {
		return packFunc(value)
	}

	return nil, fmt.Errorf("unsupported datatype: %s", datatype)
}

// Define a map for datatype to packing functions
var boolPackFunctions = map[string]func(bool) []byte{
	"1.001": func(value bool) []byte { return dpt.DPT_1001(value).Pack() },
	"1.002": func(value bool) []byte { return dpt.DPT_1002(value).Pack() },
	"1.003": func(value bool) []byte { return dpt.DPT_1003(value).Pack() },
	"1.004": func(value bool) []byte { return dpt.DPT_1004(value).Pack() },
	"1.005": func(value bool) []byte { return dpt.DPT_1005(value).Pack() },
	"1.006": func(value bool) []byte { return dpt.DPT_1006(value).Pack() },
	"1.007": func(value bool) []byte { return dpt.DPT_1007(value).Pack() },
	"1.008": func(value bool) []byte { return dpt.DPT_1008(value).Pack() },
	"1.009": func(value bool) []byte { return dpt.DPT_1009(value).Pack() },
	"1.010": func(value bool) []byte { return dpt.DPT_1010(value).Pack() },
	"1.011": func(value bool) []byte { return dpt.DPT_1011(value).Pack() },
	"1.012": func(value bool) []byte { return dpt.DPT_1012(value).Pack() },
	"1.013": func(value bool) []byte { return dpt.DPT_1013(value).Pack() },
	"1.014": func(value bool) []byte { return dpt.DPT_1014(value).Pack() },
	"1.015": func(value bool) []byte { return dpt.DPT_1015(value).Pack() },
	"1.016": func(value bool) []byte { return dpt.DPT_1016(value).Pack() },
	"1.017": func(value bool) []byte { return dpt.DPT_1017(value).Pack() },
	"1.018": func(value bool) []byte { return dpt.DPT_1018(value).Pack() },
	"1.019": func(value bool) []byte { return dpt.DPT_1019(value).Pack() },
	"1.021": func(value bool) []byte { return dpt.DPT_1021(value).Pack() },
	"1.022": func(value bool) []byte { return dpt.DPT_1022(value).Pack() },
	"1.023": func(value bool) []byte { return dpt.DPT_1023(value).Pack() },
	"1.024": func(value bool) []byte { return dpt.DPT_1024(value).Pack() },
	"1.100": func(value bool) []byte { return dpt.DPT_1100(value).Pack() },
}

var float32PackFunctions = map[string]func(float32) []byte{
	"5.001": func(value float32) []byte { return dpt.DPT_5001(value).Pack() },
	"5.003": func(value float32) []byte { return dpt.DPT_5003(value).Pack() },

	"9.001": func(value float32) []byte { return dpt.DPT_9001(value).Pack() },
	"9.002": func(value float32) []byte { return dpt.DPT_9002(value).Pack() },
	"9.003": func(value float32) []byte { return dpt.DPT_9003(value).Pack() },
	"9.004": func(value float32) []byte { return dpt.DPT_9004(value).Pack() },
	"9.005": func(value float32) []byte { return dpt.DPT_9005(value).Pack() },
	"9.006": func(value float32) []byte { return dpt.DPT_9006(value).Pack() },
	"9.007": func(value float32) []byte { return dpt.DPT_9007(value).Pack() },
	"9.008": func(value float32) []byte { return dpt.DPT_9008(value).Pack() },
	"9.010": func(value float32) []byte { return dpt.DPT_9010(value).Pack() },
	"9.011": func(value float32) []byte { return dpt.DPT_9011(value).Pack() },
	"9.020": func(value float32) []byte { return dpt.DPT_9020(value).Pack() },
	"9.021": func(value float32) []byte { return dpt.DPT_9021(value).Pack() },
	"9.022": func(value float32) []byte { return dpt.DPT_9022(value).Pack() },
	"9.023": func(value float32) []byte { return dpt.DPT_9023(value).Pack() },
	"9.024": func(value float32) []byte { return dpt.DPT_9024(value).Pack() },
	"9.025": func(value float32) []byte { return dpt.DPT_9025(value).Pack() },
	"9.026": func(value float32) []byte { return dpt.DPT_9026(value).Pack() },
	"9.027": func(value float32) []byte { return dpt.DPT_9027(value).Pack() },
	"9.028": func(value float32) []byte { return dpt.DPT_9028(value).Pack() },
	"9.029": func(value float32) []byte { return dpt.DPT_9029(value).Pack() },

	"14.000":  func(value float32) []byte { return dpt.DPT_14000(value).Pack() },
	"14.001":  func(value float32) []byte { return dpt.DPT_14001(value).Pack() },
	"14.002":  func(value float32) []byte { return dpt.DPT_14002(value).Pack() },
	"14.003":  func(value float32) []byte { return dpt.DPT_14003(value).Pack() },
	"14.004":  func(value float32) []byte { return dpt.DPT_14004(value).Pack() },
	"14.005":  func(value float32) []byte { return dpt.DPT_14005(value).Pack() },
	"14.006":  func(value float32) []byte { return dpt.DPT_14006(value).Pack() },
	"14.007":  func(value float32) []byte { return dpt.DPT_14007(value).Pack() },
	"14.008":  func(value float32) []byte { return dpt.DPT_14008(value).Pack() },
	"14.009":  func(value float32) []byte { return dpt.DPT_14009(value).Pack() },
	"14.010":  func(value float32) []byte { return dpt.DPT_14010(value).Pack() },
	"14.011":  func(value float32) []byte { return dpt.DPT_14011(value).Pack() },
	"14.012":  func(value float32) []byte { return dpt.DPT_14012(value).Pack() },
	"14.013":  func(value float32) []byte { return dpt.DPT_14013(value).Pack() },
	"14.014":  func(value float32) []byte { return dpt.DPT_14014(value).Pack() },
	"14.015":  func(value float32) []byte { return dpt.DPT_14015(value).Pack() },
	"14.016":  func(value float32) []byte { return dpt.DPT_14016(value).Pack() },
	"14.017":  func(value float32) []byte { return dpt.DPT_14017(value).Pack() },
	"14.018":  func(value float32) []byte { return dpt.DPT_14018(value).Pack() },
	"14.019":  func(value float32) []byte { return dpt.DPT_14019(value).Pack() },
	"14.020":  func(value float32) []byte { return dpt.DPT_14020(value).Pack() },
	"14.021":  func(value float32) []byte { return dpt.DPT_14021(value).Pack() },
	"14.022":  func(value float32) []byte { return dpt.DPT_14022(value).Pack() },
	"14.023":  func(value float32) []byte { return dpt.DPT_14023(value).Pack() },
	"14.024":  func(value float32) []byte { return dpt.DPT_14024(value).Pack() },
	"14.025":  func(value float32) []byte { return dpt.DPT_14025(value).Pack() },
	"14.026":  func(value float32) []byte { return dpt.DPT_14026(value).Pack() },
	"14.027":  func(value float32) []byte { return dpt.DPT_14027(value).Pack() },
	"14.028":  func(value float32) []byte { return dpt.DPT_14028(value).Pack() },
	"14.029":  func(value float32) []byte { return dpt.DPT_14029(value).Pack() },
	"14.030":  func(value float32) []byte { return dpt.DPT_14030(value).Pack() },
	"14.031":  func(value float32) []byte { return dpt.DPT_14031(value).Pack() },
	"14.032":  func(value float32) []byte { return dpt.DPT_14032(value).Pack() },
	"14.033":  func(value float32) []byte { return dpt.DPT_14033(value).Pack() },
	"14.034":  func(value float32) []byte { return dpt.DPT_14034(value).Pack() },
	"14.035":  func(value float32) []byte { return dpt.DPT_14035(value).Pack() },
	"14.036":  func(value float32) []byte { return dpt.DPT_14036(value).Pack() },
	"14.037":  func(value float32) []byte { return dpt.DPT_14037(value).Pack() },
	"14.038":  func(value float32) []byte { return dpt.DPT_14038(value).Pack() },
	"14.039":  func(value float32) []byte { return dpt.DPT_14039(value).Pack() },
	"14.040":  func(value float32) []byte { return dpt.DPT_14040(value).Pack() },
	"14.041":  func(value float32) []byte { return dpt.DPT_14041(value).Pack() },
	"14.042":  func(value float32) []byte { return dpt.DPT_14042(value).Pack() },
	"14.043":  func(value float32) []byte { return dpt.DPT_14043(value).Pack() },
	"14.044":  func(value float32) []byte { return dpt.DPT_14044(value).Pack() },
	"14.045":  func(value float32) []byte { return dpt.DPT_14045(value).Pack() },
	"14.046":  func(value float32) []byte { return dpt.DPT_14046(value).Pack() },
	"14.047":  func(value float32) []byte { return dpt.DPT_14047(value).Pack() },
	"14.048":  func(value float32) []byte { return dpt.DPT_14048(value).Pack() },
	"14.049":  func(value float32) []byte { return dpt.DPT_14049(value).Pack() },
	"14.050":  func(value float32) []byte { return dpt.DPT_14050(value).Pack() },
	"14.051":  func(value float32) []byte { return dpt.DPT_14051(value).Pack() },
	"14.052":  func(value float32) []byte { return dpt.DPT_14052(value).Pack() },
	"14.053":  func(value float32) []byte { return dpt.DPT_14053(value).Pack() },
	"14.054":  func(value float32) []byte { return dpt.DPT_14054(value).Pack() },
	"14.055":  func(value float32) []byte { return dpt.DPT_14055(value).Pack() },
	"14.056":  func(value float32) []byte { return dpt.DPT_14056(value).Pack() },
	"14.057":  func(value float32) []byte { return dpt.DPT_14057(value).Pack() },
	"14.058":  func(value float32) []byte { return dpt.DPT_14058(value).Pack() },
	"14.059":  func(value float32) []byte { return dpt.DPT_14059(value).Pack() },
	"14.060":  func(value float32) []byte { return dpt.DPT_14060(value).Pack() },
	"14.061":  func(value float32) []byte { return dpt.DPT_14061(value).Pack() },
	"14.062":  func(value float32) []byte { return dpt.DPT_14062(value).Pack() },
	"14.063":  func(value float32) []byte { return dpt.DPT_14063(value).Pack() },
	"14.064":  func(value float32) []byte { return dpt.DPT_14064(value).Pack() },
	"14.065":  func(value float32) []byte { return dpt.DPT_14065(value).Pack() },
	"14.066":  func(value float32) []byte { return dpt.DPT_14066(value).Pack() },
	"14.067":  func(value float32) []byte { return dpt.DPT_14067(value).Pack() },
	"14.068":  func(value float32) []byte { return dpt.DPT_14068(value).Pack() },
	"14.069":  func(value float32) []byte { return dpt.DPT_14069(value).Pack() },
	"14.070":  func(value float32) []byte { return dpt.DPT_14070(value).Pack() },
	"14.071":  func(value float32) []byte { return dpt.DPT_14071(value).Pack() },
	"14.072":  func(value float32) []byte { return dpt.DPT_14072(value).Pack() },
	"14.073":  func(value float32) []byte { return dpt.DPT_14073(value).Pack() },
	"14.074":  func(value float32) []byte { return dpt.DPT_14074(value).Pack() },
	"14.075":  func(value float32) []byte { return dpt.DPT_14075(value).Pack() },
	"14.076":  func(value float32) []byte { return dpt.DPT_14076(value).Pack() },
	"14.077":  func(value float32) []byte { return dpt.DPT_14077(value).Pack() },
	"14.078":  func(value float32) []byte { return dpt.DPT_14078(value).Pack() },
	"14.079":  func(value float32) []byte { return dpt.DPT_14079(value).Pack() },
	"14.1200": func(value float32) []byte { return dpt.DPT_141200(value).Pack() },
}

var uint8PackFunctions = map[string]func(uint8) []byte{
	"5.004": func(value uint8) []byte { return dpt.DPT_5004(value).Pack() },
	"5.005": func(value uint8) []byte { return dpt.DPT_5005(value).Pack() },

	"17.001": func(value uint8) []byte { return dpt.DPT_17001(value).Pack() },
	"18.001": func(value uint8) []byte { return dpt.DPT_18001(value).Pack() },
	"20.102": func(value uint8) []byte { return dpt.DPT_20102(value).Pack() },
	"20.105": func(value uint8) []byte { return dpt.DPT_20105(value).Pack() },
}

var int8PackFunctions = map[string]func(int8) []byte{
	"6.010": func(value int8) []byte { return dpt.DPT_6010(value).Pack() },
}

var uint16PackFunctions = map[string]func(uint16) []byte{
	"7.001": func(value uint16) []byte { return dpt.DPT_7001(value).Pack() },
	"7.002": func(value uint16) []byte { return dpt.DPT_7002(value).Pack() },
	"7.003": func(value uint16) []byte { return dpt.DPT_7003(value).Pack() },
	"7.004": func(value uint16) []byte { return dpt.DPT_7004(value).Pack() },
	"7.005": func(value uint16) []byte { return dpt.DPT_7005(value).Pack() },
	"7.006": func(value uint16) []byte { return dpt.DPT_7006(value).Pack() },
	"7.007": func(value uint16) []byte { return dpt.DPT_7007(value).Pack() },
	"7.010": func(value uint16) []byte { return dpt.DPT_7010(value).Pack() },
	"7.011": func(value uint16) []byte { return dpt.DPT_7011(value).Pack() },
	"7.012": func(value uint16) []byte { return dpt.DPT_7012(value).Pack() },
	"7.013": func(value uint16) []byte { return dpt.DPT_7013(value).Pack() },
	"7.600": func(value uint16) []byte { return dpt.DPT_7600(value).Pack() },
}

var int16PackFunctions = map[string]func(int16) []byte{
	"8.001": func(value int16) []byte { return dpt.DPT_8001(value).Pack() },
	"8.002": func(value int16) []byte { return dpt.DPT_8002(value).Pack() },
	"8.003": func(value int16) []byte { return dpt.DPT_8003(value).Pack() },
	"8.004": func(value int16) []byte { return dpt.DPT_8004(value).Pack() },
	"8.005": func(value int16) []byte { return dpt.DPT_8005(value).Pack() },
	"8.006": func(value int16) []byte { return dpt.DPT_8006(value).Pack() },
	"8.007": func(value int16) []byte { return dpt.DPT_8007(value).Pack() },
	"8.010": func(value int16) []byte { return dpt.DPT_8010(value).Pack() },
	"8.011": func(value int16) []byte { return dpt.DPT_8011(value).Pack() },
}

var uint32PackFunctions = map[string]func(uint32) []byte{
	"12.001": func(value uint32) []byte { return dpt.DPT_12001(value).Pack() },
}

var int32PackFunctions = map[string]func(int32) []byte{
	"13.001": func(value int32) []byte { return dpt.DPT_13001(value).Pack() },
	"13.002": func(value int32) []byte { return dpt.DPT_13002(value).Pack() },
	"13.010": func(value int32) []byte { return dpt.DPT_13010(value).Pack() },
	"13.011": func(value int32) []byte { return dpt.DPT_13011(value).Pack() },
	"13.012": func(value int32) []byte { return dpt.DPT_13012(value).Pack() },
	"13.013": func(value int32) []byte { return dpt.DPT_13013(value).Pack() },
	"13.014": func(value int32) []byte { return dpt.DPT_13014(value).Pack() },
	"13.015": func(value int32) []byte { return dpt.DPT_13015(value).Pack() },
	"13.016": func(value int32) []byte { return dpt.DPT_13016(value).Pack() },
	"13.100": func(value int32) []byte { return dpt.DPT_13100(value).Pack() },
}

var stringPackFunctions = map[string]func(string) ([]byte, error){
	"19.001": func(value string) ([]byte, error) {
		// Try to parse as ISO format first (e.g. "2023-05-15T14:30:45")
		t, err := time.Parse(time.RFC3339, value)
		if err == nil {
			return localdpt.FromTime(t).Pack(), nil
		}

		// Try to parse using our custom format
		matches := regexpDateTime.FindStringSubmatch(value)
		if matches != nil {
			year, _ := strconv.ParseUint(matches[1], 10, 16)
			month, _ := strconv.ParseUint(matches[2], 10, 8)
			day, _ := strconv.ParseUint(matches[3], 10, 8)

			var weekday uint8
			if len(matches) > 4 && matches[4] != "" {
				if day, ok := weekdays[matches[4]]; ok {
					weekday = day
				}
			}

			hour, _ := strconv.ParseUint(matches[5], 10, 8)
			minute, _ := strconv.ParseUint(matches[6], 10, 8)
			second, _ := strconv.ParseUint(matches[7], 10, 8)

			// Check if summer time is specified
			summerTime := false
			if len(matches) > 7 && strings.Contains(value, "(Summer Time)") {
				summerTime = true
			}

			// Parse flags if present
			fault := false
			workingDay := false
			externalSync := false
			reliableSync := false

			if len(matches) > 8 && matches[8] != "" {
				flags := matches[8]
				fault = strings.Contains(flags, "Fault")
				workingDay = strings.Contains(flags, "Working Day")
				externalSync = strings.Contains(flags, "External Sync")
				reliableSync = strings.Contains(flags, "Reliable Sync")
			}

			// Convert year to KNX format (offset from 1900)
			yearOffset := uint8(0)
			if year >= 1900 && year <= 2155 {
				yearOffset = uint8(year - 1900)
			}

			datapoint := localdpt.DPT_19001{
				Year:         yearOffset,
				Month:        uint8(month),
				DayOfMonth:   uint8(day),
				DayOfWeek:    weekday,
				HourOfDay:    uint8(hour),
				Minutes:      uint8(minute),
				Seconds:      uint8(second),
				Fault:        fault,
				WorkingDay:   workingDay,
				SummerTime:   summerTime,
				ExternalSync: externalSync,
				ReliableSync: reliableSync,
			}

			return datapoint.Pack(), nil
		}

		// If all parsing attempts fail, return nil
		return nil, fmt.Errorf("cannot convert \"%s\" to DPT 19.001 (Date/Time): unrecognized format", value)
	},
	"10.001": func(value string) ([]byte, error) {
		matches := regexpTimeOfDay.FindStringSubmatch(value)
		if matches == nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 10.001 (Time): unrecognized format", value)
		}

		var weekday uint8
		if day, ok := weekdays[matches[1]]; ok {
			weekday = day
		}

		hour, _ := strconv.ParseUint(matches[2], 10, 8)
		minute, _ := strconv.ParseUint(matches[3], 10, 8)
		second, _ := strconv.ParseUint(matches[4], 10, 8)

		datapoint := dpt.DPT_10001{
			Weekday: weekday,
			Hour:    uint8(hour),
			Minutes: uint8(minute),
			Seconds: uint8(second),
		}
		return datapoint.Pack(), nil
	},
	"11.001": func(value string) ([]byte, error) {
		matches := regexpDate.FindStringSubmatch(value)
		if matches == nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 11.001 (Date): unrecognized format", value)
		}

		year, _ := strconv.ParseUint(matches[1], 10, 16)
		month, _ := strconv.ParseUint(matches[2], 10, 8)
		day, _ := strconv.ParseUint(matches[3], 10, 8)

		datapoint := dpt.DPT_11001{
			Year:  uint16(year),
			Month: uint8(month),
			Day:   uint8(day),
		}

		return datapoint.Pack(), nil
	},
	"16.000": func(value string) ([]byte, error) { return dpt.DPT_16000(value).Pack(), nil },
	"16.001": func(value string) ([]byte, error) { return dpt.DPT_16001(value).Pack(), nil },
	"28.001": func(value string) ([]byte, error) { return dpt.DPT_28001(value).Pack(), nil },
	"232.600": func(value string) ([]byte, error) {
		matches := regexpRgb.FindStringSubmatch(value)
		if matches == nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 232.600 (Color RGB): unrecognized format", value)
		}

		red, err := strconv.ParseUint(matches[1], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 232.600 (Color RGB): red component not a positive integer", value)
		}
		green, err := strconv.ParseUint(matches[2], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 232.600 (Color RGB): green component not a positive integer", value)
		}
		blue, err := strconv.ParseUint(matches[3], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 232.600 (Color RGB): blue component not a positive integer", value)
		}

		datapoint := dpt.DPT_232600{
			Red:   uint8(red),
			Green: uint8(green),
			Blue:  uint8(blue),
		}

		return datapoint.Pack(), nil
	},
	"242.600": func(value string) ([]byte, error) {
		matches := regexpXyy.FindStringSubmatch(value)
		if matches == nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 242.600 (Color xyY): unrecognized format", value)
		}

		// Parse the extracted strings to appropriate types
		x, err := strconv.ParseUint(matches[1], 10, 16)
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 242.600 (Color xyY): x component not a positive integer", value)
		}
		y, err := strconv.ParseUint(matches[2], 10, 16)
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 242.600 (Color xyY): y component not a positive integer", value)
		}
		yBrightness, err := strconv.ParseUint(matches[3], 10, 8)
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 242.600 (Color xyY): Y brightness component not a positive integer", value)
		}
		colorValid, err := strconv.ParseBool(matches[4])
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 242.600 (Color xyY): color valid component not a boolean", value)
		}
		brightnessValid, err := strconv.ParseBool(matches[5])
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 242.600 (Color xyY): brightness valid component not a boolean", value)
		}

		datapoint := dpt.DPT_242600{
			X:               uint16(x),
			Y:               uint16(y),
			YBrightness:     uint8(yBrightness),
			ColorValid:      colorValid,
			BrightnessValid: brightnessValid,
		}

		return datapoint.Pack(), nil
	},
	"251.600": func(value string) ([]byte, error) {
		matches := regexpRgbw.FindStringSubmatch(value)
		if matches == nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 251.600 (Color RGBW): unrecognized format", value)
		}

		// Parse the extracted strings to appropriate types
		red, err := strconv.ParseUint(matches[1], 10, 8)
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 251.600 (Color RGBW): red component not a positive integer", value)
		}
		green, err := strconv.ParseUint(matches[2], 10, 8)
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 251.600 (Color RGBW): green component not a positive integer", value)
		}
		blue, err := strconv.ParseUint(matches[3], 10, 8)
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 251.600 (Color RGBW): blue component not a positive integer", value)
		}
		white, err := strconv.ParseUint(matches[4], 10, 8)
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 251.600 (Color RGBW): white component not a positive integer", value)
		}
		redValid, err := strconv.ParseBool(matches[5])
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 251.600 (Color RGBW): red valid component not a boolean", value)
		}
		greenValid, err := strconv.ParseBool(matches[6])
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 251.600 (Color RGBW): green valid component not a boolean", value)
		}
		blueValid, err := strconv.ParseBool(matches[7])
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 251.600 (Color RGBW): blue valid component not a boolean", value)
		}
		whiteValid, err := strconv.ParseBool(matches[8])
		if err != nil {
			return nil, fmt.Errorf("cannot convert \"%s\" to DPT 251.600 (Color RGBW): white valid component not a boolean", value)
		}

		// Create the DPT_251600 instance
		datapoint := dpt.DPT_251600{
			Red:        uint8(red),
			Green:      uint8(green),
			Blue:       uint8(blue),
			White:      uint8(white),
			RedValid:   redValid,
			GreenValid: greenValid,
			BlueValid:  blueValid,
			WhiteValid: whiteValid,
		}

		return datapoint.Pack(), nil
	},
}

func ExtractDatapointValue(dp dpt.Datapoint, dptType string) interface{} {
	mainType := strings.SplitN(dptType, ".", 2)[0]
	rv := reflect.ValueOf(dp)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	switch mainType {
	case "1":
		if rv.Kind() == reflect.Bool {
			return rv.Bool()
		}
	case "6":
		if rv.Kind() == reflect.Int8 {
			return int8(rv.Int())
		}
	case "7":
		if rv.Kind() == reflect.Uint16 {
			return uint16(rv.Uint())
		}
	case "9", "14":
		if rv.Kind() == reflect.Float32 {
			return float32(rv.Float())
		}
	case "12":
		if rv.Kind() == reflect.Uint32 {
			return uint32(rv.Uint())
		}
	case "13":
		if rv.Kind() == reflect.Int32 {
			return int32(rv.Int())
		}
	case "17", "18", "20":
		if rv.Kind() == reflect.Uint8 {
			return uint8(rv.Uint())
		}
	case "19":
		// For DateTime, properly extract the date/time value
		if dptType == "19.001" {
			// Check if we're dealing with our custom DPT_19001 type
			if dpt19, ok := dp.(*localdpt.DPT_19001); ok {
				return dpt19.ToTime().Format(time.RFC3339)
			}
		}
		// Fallback to string representation if not our custom type
		return StringWithoutSuffix(dp)
	}

	switch dptType {
	case "5.001", "5.003", "8.003", "8.004", "8.010":
		if rv.Kind() == reflect.Float32 {
			return float32(rv.Float())
		}
	case "5.004", "5.005":
		if rv.Kind() == reflect.Uint8 {
			return uint8(rv.Uint())
		}
	case "8.001", "8.002", "8.005", "8.006", "8.007", "8.011":
		if rv.Kind() == reflect.Int16 {
			return int16(rv.Int())
		}
	}

	return StringWithoutSuffix(dp)
}
