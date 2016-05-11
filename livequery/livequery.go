package livequery

import "github.com/lfq7413/tomato/types"

// LiveQuery ...
type LiveQuery struct {
	classNames []string
}

// NewLiveQuery ...
func NewLiveQuery(classNames []string) *LiveQuery {
	liveQuery := &LiveQuery{}
	if classNames == nil && len(classNames) == 0 {
		liveQuery.classNames = []string{}
	} else {
		liveQuery.classNames = classNames
	}
	// TODO 添加 Publisher

	return liveQuery
}

// OnAfterSave ...
func (l *LiveQuery) OnAfterSave(className string, currentObject, originalObject types.M) {

}

// OnAfterDelete ...
func (l *LiveQuery) OnAfterDelete(className string, currentObject, originalObject types.M) {

}

// HasLiveQuery ...
func (l *LiveQuery) HasLiveQuery(className string) {

}

func (l *LiveQuery) makePublisherRequest(currentObject, originalObject types.M) types.M {
	return nil
}
