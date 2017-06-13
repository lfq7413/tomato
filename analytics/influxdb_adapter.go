package analytics

import (
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

type influxDBAdapter struct {
	c            client.Client
	databaseName string
}

func newInfluxDBAdapter() *influxDBAdapter {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.TConfig.InfluxDBURL,
		Username: config.TConfig.InfluxDBUsername,
		Password: config.TConfig.InfluxDBPassword,
	})
	if err != nil {
		panic(err)
	}
	return &influxDBAdapter{
		c:            c,
		databaseName: config.TConfig.InfluxDBDatabaseName,
	}
}

func (a *influxDBAdapter) appOpened(body types.M) (types.M, error) {
	err := a.addEvent("AppOpened", body)
	return types.M{}, err
}

func (a *influxDBAdapter) trackEvent(eventName string, body types.M) (types.M, error) {
	err := a.addEvent(eventName, body)
	return types.M{}, err
}

func (a *influxDBAdapter) addEvent(name string, event types.M) error {
	var at time.Time
	if atM := utils.M(event["at"]); atM != nil {
		if utils.S(atM["__type"]) == "Date" || utils.S(atM["iso"]) != "" {
			t, err := utils.StringtoTime(utils.S(atM["iso"]))
			if err != nil {
				at = time.Now()
			} else {
				at = t
			}
		}
	}
	if at.IsZero() {
		at = time.Now()
	}

	fields := types.M{}
	if dimensions := utils.M(event["dimensions"]); dimensions != nil {
		for k, v := range dimensions {
			fields[k] = v
		}
	}
	if len(fields) == 0 {
		fields["_noFields"] = true
	}

	tags := map[string]string{
		name: name + "-total",
	}
	if enevtTag := utils.M(event["tags"]); enevtTag != nil {
		for k, v := range enevtTag {
			if tag := utils.S(v); tag != "" {
				tags[k] = tag
			}
		}
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  a.databaseName,
		Precision: "ns",
	})
	if err != nil {
		return err
	}

	pt, err := client.NewPoint(
		name,
		tags,
		fields,
		at,
	)
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	return a.c.Write(bp)
}
