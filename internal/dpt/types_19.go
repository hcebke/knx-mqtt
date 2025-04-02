package dpt

import (
	"fmt"
	"time"

	"github.com/vapourismo/knx-go/knx/dpt"
)

// DPT_19001 represents DPT 19.001 / DateTime.
// It contains date and time information according to the KNX standard.
type DPT_19001 struct {
	Year       uint8 // 0 = 1900, 255 = 2155
	Month      uint8 // 1 = January, 12 = December
	DayOfMonth uint8 // 1-31
	DayOfWeek  uint8 // 0 = any day, 1 = Monday, 7 = Sunday
	HourOfDay  uint8 // 0-24
	Minutes    uint8 // 0-59
	Seconds    uint8 // 0-59

	// Flags
	Fault       bool // 0 = Normal (No fault), 1 = Fault
	WorkingDay  bool // 0 = Bank day (No working day), 1 = Working day
	NoWD        bool // 0 = WD field valid, 1 = WD field not valid
	NoYear      bool // 0 = Year field valid, 1 = Year field not valid
	NoDate      bool // 0 = Month and Day of Month fields valid, 1 = Month and Day of Month fields not valid
	NoDayOfWeek bool // 0 = Day of week field valid, 1 = Day of week field not valid
	NoTime      bool // 0 = Hour of day, Minutes and Seconds fields valid, 1 = Hour of day, Minutes and Seconds fields not valid
	SummerTime  bool // 0 = Time = UT+X, 1 = Time = UT+X+1

	// Clock Quality
	ExternalSync bool // 0 = clock without ext. sync signal, 1 = clock with ext. sync signal
	ReliableSync bool // 0 = unreliable synchronisation source (mains, local quartz), 1 = reliable synchronisation source (radio, Internet)
}

func (d DPT_19001) Pack() []byte {
	// Create a 9-byte array with a leading 0-byte (padding)
	var buf = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}

	// Octet 8 (MSB) - Year
	buf[1] = d.Year

	// Octet 7 - Month
	buf[2] = d.Month & 0x0F

	// Octet 6 - Day of Month
	buf[3] = d.DayOfMonth & 0x1F

	// Octet 5 - Day of Week and Hour
	buf[4] = (d.DayOfWeek & 0x07) << 5
	buf[4] |= d.HourOfDay & 0x1F

	// Octet 4 - Minutes
	buf[5] = d.Minutes & 0x3F

	// Octet 3 - Seconds
	buf[6] = d.Seconds & 0x3F

	// Octet 2 - Flags
	if d.Fault {
		buf[7] |= 0x80
	}
	if d.WorkingDay {
		buf[7] |= 0x40
	}
	if d.NoWD {
		buf[7] |= 0x20
	}
	if d.NoYear {
		buf[7] |= 0x10
	}
	if d.NoDate {
		buf[7] |= 0x08
	}
	if d.NoDayOfWeek {
		buf[7] |= 0x04
	}
	if d.NoTime {
		buf[7] |= 0x02
	}
	if d.SummerTime {
		buf[7] |= 0x01
	}

	// Octet 1 (LSB) - Clock Quality and Reserved Bits
	if d.ExternalSync {
		buf[8] |= 0x80
	}
	if d.ReliableSync {
		buf[8] |= 0x40
	}

	return buf
}

func (d *DPT_19001) Unpack(data []byte) error {

	if len(data) != 9 {
		return dpt.ErrInvalidLength
	}

	// First byte (index 0) is a padding byte, so we skip it

	// Octet 8 (MSB) - Year
	d.Year = data[1]

	// Octet 7 - Month
	d.Month = data[2] & 0x0F

	// Octet 6 - Day of Month
	d.DayOfMonth = data[3] & 0x1F

	// Octet 5 - Day of Week and Hour
	d.DayOfWeek = (data[4] >> 5) & 0x07
	d.HourOfDay = data[4] & 0x1F

	// Octet 4 - Minutes
	d.Minutes = data[5] & 0x3F

	// Octet 3 - Seconds
	d.Seconds = data[6] & 0x3F

	// Octet 2 - Flags
	d.Fault = (data[7] & 0x80) != 0
	d.WorkingDay = (data[7] & 0x40) != 0
	d.NoWD = (data[7] & 0x20) != 0
	d.NoYear = (data[7] & 0x10) != 0
	d.NoDate = (data[7] & 0x08) != 0
	d.NoDayOfWeek = (data[7] & 0x04) != 0
	d.NoTime = (data[7] & 0x02) != 0
	d.SummerTime = (data[7] & 0x01) != 0

	// Octet 1 (LSB) - Clock Quality and Reserved Bits
	d.ExternalSync = (data[8] & 0x80) != 0
	d.ReliableSync = (data[8] & 0x40) != 0

	if !d.IsValid() {
		return fmt.Errorf("payload is out of range")
	}

	return nil
}

func (d DPT_19001) Unit() string {
	return ""
}

func (d DPT_19001) IsValid() bool {
	// Check if the date and time fields are valid
	if !d.NoDate {
		if d.Month < 1 || d.Month > 12 {
			return false
		}

		if d.DayOfMonth < 1 || d.DayOfMonth > 31 {
			return false
		}

		// Check if the day is valid for the month
		if d.Month == 2 {
			// February
			year := 1900 + int(d.Year)
			isLeapYear := (year%4 == 0 && year%100 != 0) || (year%400 == 0)
			if isLeapYear && d.DayOfMonth > 29 {
				return false
			}
			if !isLeapYear && d.DayOfMonth > 28 {
				return false
			}
		} else if d.Month == 4 || d.Month == 6 || d.Month == 9 || d.Month == 11 {
			// April, June, September, November
			if d.DayOfMonth > 30 {
				return false
			}
		}
	}

	if !d.NoDayOfWeek && d.DayOfWeek > 7 {
		return false
	}

	if !d.NoTime {
		if d.HourOfDay > 24 {
			return false
		}

		if d.Minutes > 59 {
			return false
		}

		if d.Seconds > 59 {
			return false
		}
	}

	return true
}

func (d DPT_19001) String() string {
	if d.NoYear && d.NoDate && d.NoTime {
		return "Invalid DateTime: No valid fields"
	}

	var result string

	// Format date if available
	if !d.NoDate {
		year := 1900 + int(d.Year)
		if d.NoYear {
			result = fmt.Sprintf("%02d-%02d", d.Month, d.DayOfMonth)
		} else {
			result = fmt.Sprintf("%04d-%02d-%02d", year, d.Month, d.DayOfMonth)
		}
	}

	// Add day of week if available
	if !d.NoDayOfWeek && d.DayOfWeek > 0 && d.DayOfWeek <= 7 {
		weekday := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		if result != "" {
			result += " "
		}
		result += weekday[d.DayOfWeek-1]
	}

	// Add time if available
	if !d.NoTime {
		if result != "" {
			result += " "
		}
		result += fmt.Sprintf("%02d:%02d:%02d", d.HourOfDay, d.Minutes, d.Seconds)

		// Add summer time indicator
		if d.SummerTime {
			result += " (Summer Time)"
		}
	}

	// Add flags
	var flags []string
	if d.Fault {
		flags = append(flags, "Fault")
	}
	if d.WorkingDay {
		flags = append(flags, "Working Day")
	}
	if d.ExternalSync {
		flags = append(flags, "External Sync")
	}
	if d.ReliableSync {
		flags = append(flags, "Reliable Sync")
	}

	if len(flags) > 0 {
		result += " [" + fmt.Sprintf("%s", flags) + "]"
	}

	return result
}

// ToTime converts the DPT_19001 to a Go time.Time object
// If any of the date or time fields are marked as invalid, it returns nil
func (d DPT_19001) ToTime() *time.Time {
	if d.NoYear || d.NoDate || d.NoTime {
		return nil
	}

	year := 1900 + int(d.Year)
	t := time.Date(year, time.Month(d.Month), int(d.DayOfMonth),
		int(d.HourOfDay), int(d.Minutes), int(d.Seconds), 0, time.Local)

	return &t
}

// FromTime creates a DPT_19001 from a Go time.Time object
func FromTime(t time.Time) DPT_19001 {
	year := t.Year()
	if year < 1900 {
		year = 1900
	} else if year > 2155 {
		year = 2155
	}

	weekday := t.Weekday()
	// Convert Go's weekday (0=Sunday) to KNX weekday (1=Monday, 7=Sunday)
	knxWeekday := uint8((int(weekday)+6)%7 + 1)

	return DPT_19001{
		Year:         uint8(year - 1900),
		Month:        uint8(t.Month()),
		DayOfMonth:   uint8(t.Day()),
		DayOfWeek:    knxWeekday,
		HourOfDay:    uint8(t.Hour()),
		Minutes:      uint8(t.Minute()),
		Seconds:      uint8(t.Second()),
		NoYear:       false,
		NoDate:       false,
		NoDayOfWeek:  false,
		NoTime:       false,
		SummerTime:   t.IsDST(),
		ExternalSync: true, // Assuming system time is externally synced
		ReliableSync: true, // Assuming system time is reliable
	}
}
