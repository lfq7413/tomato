package livequery

import "github.com/lfq7413/tomato/types"

// LiveQuery ...
type LiveQuery struct {
	classNames         []string
	liveQueryPublisher *cloudCodePublisher
}

// NewLiveQuery 初始化 LiveQuery
func NewLiveQuery(classNames []string, pubType, pubURL string) *LiveQuery {
	liveQuery := &LiveQuery{}
	if classNames == nil && len(classNames) == 0 {
		liveQuery.classNames = []string{}
	} else {
		liveQuery.classNames = classNames
	}
	liveQuery.liveQueryPublisher = newCloudCodePublisher(pubType, pubURL)

	return liveQuery
}

// OnAfterSave 保存对象之后调用
func (l *LiveQuery) OnAfterSave(className string, currentObject, originalObject types.M) {
	if l.HasLiveQuery(className) == false {
		return
	}
	req := l.makePublisherRequest(currentObject, originalObject)
	l.liveQueryPublisher.onCloudCodeAfterSave(req)
}

// OnAfterDelete 删除对象之后调用
func (l *LiveQuery) OnAfterDelete(className string, currentObject, originalObject types.M) {
	if l.HasLiveQuery(className) == false {
		return
	}
	req := l.makePublisherRequest(currentObject, originalObject)
	l.liveQueryPublisher.onCloudCodeAfterDelete(req)
}

// HasLiveQuery 是否有对应的 className
func (l *LiveQuery) HasLiveQuery(className string) bool {
	for _, n := range l.classNames {
		if n == className {
			return true
		}
	}
	return false
}

func (l *LiveQuery) makePublisherRequest(currentObject, originalObject types.M) types.M {
	req := types.M{
		"object": currentObject,
	}
	if currentObject != nil {
		req["original"] = originalObject
	}
	return req
}
