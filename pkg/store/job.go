package store

import (
	"context"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
	"time"
)

type job struct {
	engineStore        bpmn_engine_store.IBpmnEngineStore
	ProcessInstanceKey bpmn_engine_store.IProcessInstanceKey
	ElementInstanceKey bpmn_engine_store.IJobKey
	ElementId          string
	State              activity.LifecycleState
	CreatedAt          time.Time
	Reason             string
	JobKey             int64
}

func (job *job) SetVariable(ctx context.Context, key string, value interface{}) error {
	return job.engineStore.SetProcessInstanceVariable(ctx, job.ProcessInstanceKey, key, value)
}

func (job *job) GetVariable(ctx context.Context, key string) (interface{}, error) {
	return job.engineStore.GetProcessInstanceVariable(ctx, job.ProcessInstanceKey, key)
}

func (job *job) SetState(ctx context.Context, state activity.LifecycleState, reason string) error {
	return job.engineStore.SetJobState(ctx, job.ProcessInstanceKey, job.ElementInstanceKey, state, reason)
}

func (job *job) GetState(ctx context.Context) (activity.LifecycleState, error) {
	return job.engineStore.GetJobState(ctx, job.ProcessInstanceKey, job.ElementInstanceKey)
}

func (job *job) GetElementId() string {
	return job.ElementId
}

func (job *job) GetElementInstanceKey() bpmn_engine_store.IJobKey {
	return job.ElementInstanceKey
}
