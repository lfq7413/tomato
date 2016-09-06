package job

import (
	"time"

	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

const jobStatusCollection = "_JobStatus"

// JobStatus ...
type JobStatus struct {
	objectID string
	status   types.M
	db       *orm.DBController
}

// NewjobStatus ...
func NewjobStatus() *JobStatus {
	p := &JobStatus{
		objectID: utils.CreateObjectID(),
		db:       orm.TomatoDBController,
	}
	return p
}

// SetRunning ...
func (j *JobStatus) SetRunning(jobName string, params types.M) types.M {
	now := time.Now().UTC()
	j.status = types.M{
		"objectId":  j.objectID,
		"jobName":   jobName,
		"params":    params,
		"status":    "running",
		"source":    "api",
		"createdAt": utils.TimetoString(now),
		// lockdown!
		"ACL": types.M{},
	}
	j.db.Create(jobStatusCollection, j.status, types.M{})
	return j.status
}

// SetMessage ...
func (j *JobStatus) SetMessage(message string) {
	j.db.Update(jobStatusCollection, types.M{"objectId": j.objectID}, types.M{"message": message}, types.M{}, false)
}

// SetSucceeded ...
func (j *JobStatus) SetSucceeded(message string) {
	j.setFinalStatus("succeeded", message)
}

// SetFailed ...
func (j *JobStatus) SetFailed(message string) {
	j.setFinalStatus("failed", message)
}

// setFinalStatus ...
func (j *JobStatus) setFinalStatus(status, message string) {
	finishedAt := time.Now().UTC()
	update := types.M{
		"status":     status,
		"finishedAt": utils.TimetoString(finishedAt),
		"message":    message,
	}
	j.db.Update(jobStatusCollection, types.M{"objectId": j.objectID}, update, types.M{}, false)
}
