package analytics

import (
	"github.com/influxdata/influxdb/client/v2"
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/types"
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
	return types.M{}, nil
}

func (a *influxDBAdapter) trackEvent(eventName string, body types.M) (types.M, error) {
	return types.M{}, nil
}
