package analytics

import (
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/types"
)

var adapter analyticsAdapter

func init() {
	if config.TConfig.AnalyticsAdapter == "InfluxDB" {
		adapter = newInfluxDBAdapter()
	} else {
		adapter = &nullAnalyticsAdapter{}
	}
}

// AppOpened 统计应用打开记录
func AppOpened(body types.M) types.M {
	response, err := adapter.appOpened(body)
	if err != nil {
		return types.M{}
	}
	return response
}

// TrackEvent 统计自定义事件
func TrackEvent(eventName string, body types.M) types.M {
	response, err := adapter.trackEvent(eventName, body)
	if err != nil {
		return types.M{}
	}
	return response
}

type analyticsAdapter interface {
	appOpened(body types.M) (types.M, error)
	trackEvent(eventName string, body types.M) (types.M, error)
}
