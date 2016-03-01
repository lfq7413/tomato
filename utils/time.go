package utils

import "time"

// ISO8601 ...
const ISO8601 = "2006-01-02T15:04:05.000Z"

// StringtoTime ...
func StringtoTime(s string) (time.Time, error) {
	return time.ParseInLocation(ISO8601, s, time.UTC)
}

// TimetoString ...
func TimetoString(t time.Time) string {
	return t.Format(ISO8601)
}

// UnixmillitoTime ...
func UnixmillitoTime(m int64) time.Time {
	return time.Unix(m/1000, m%1000*1e6).UTC()
}

// TimetoUnixmilli ...
func TimetoUnixmilli(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

// StringtoUnixmilli ...
func StringtoUnixmilli(s string) (int64, error) {
	var m int64
	t, err := StringtoTime(s)
	if err != nil {
		return m, err
	}
	m = TimetoUnixmilli(t)
	return m, err
}

// UnixmillitoString ...
func UnixmillitoString(m int64) string {
	return TimetoString(UnixmillitoTime(m))
}
