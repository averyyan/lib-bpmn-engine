package store

import (
	"context"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
	"time"
)

type job struct {
	ProcessInstanceInfo bpmn_engine_store.IProcessInstanceInfo
	ElementInstanceKey  bpmn_engine_store.IJobKey
	ElementId           string
	State               activity.LifecycleState
	CreatedAt           time.Time
	Reason              string
	JobKey              int64
}

func (job *job) SetVariable(ctx context.Context, key string, value interface{}) error {
	return job.ProcessInstanceInfo.SetVariable(ctx, key, value)
}

func (job *job) GetVariable(ctx context.Context, key string) (interface{}, error) {
	return job.ProcessInstanceInfo.GetVariable(ctx, key)
}

func (job *job) SetState(ctx context.Context, state activity.LifecycleState, reason string) error {
	return job.ProcessInstanceInfo.SetJobState(ctx, job.ElementInstanceKey, state, reason)
}

func (job *job) GetState(ctx context.Context) (activity.LifecycleState, error) {
	return job.ProcessInstanceInfo.GetJobState(ctx, job.ElementInstanceKey)
}

func (job *job) GetElementId() string {
	return job.ElementId
}

func (job *job) GetElementInstanceKey() bpmn_engine_store.IJobKey {
	return job.ElementInstanceKey
}
