package store

import (
	"context"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/process_instance"
	"time"
)

// ProcessInstanceInfo need save to db processInfo is from ProcessInfo
type ProcessInstanceInfo struct {
	engineState     bpmn_engine_store.IBpmnEngine
	processInfo     bpmn_engine_store.IProcessInfo
	instanceKey     bpmn_engine_store.IProcessInstanceKey
	variableContext map[string]interface{}
	createdAt       time.Time
	state           process_instance.State
	caughtEvents    []catchEvent
}

func (instance *ProcessInstanceInfo) GetVariable(ctx context.Context, key string) (interface{}, error) {
	return instance.engineState.GetStore().GetProcessInstanceVariable(ctx, instance.instanceKey, key)
}

func (instance *ProcessInstanceInfo) CreateJob(ctx context.Context, elementId string) (bpmn_engine_store.IJob, error) {
	return instance.engineState.GetStore().CreateJob(ctx, instance, elementId, bpmn_engine_store.IJobKey(instance.engineState.GenerateKey()))
}

func (instance *ProcessInstanceInfo) GetJobState(ctx context.Context, jobKey bpmn_engine_store.IJobKey) (activity.LifecycleState, error) {
	return instance.engineState.GetStore().GetJobState(ctx, instance.instanceKey, jobKey)
}

func (instance *ProcessInstanceInfo) SetJobState(ctx context.Context, jobKey bpmn_engine_store.IJobKey, state activity.LifecycleState, reason string) error {
	return instance.engineState.GetStore().SetJobState(ctx, instance.instanceKey, jobKey, state, reason)
}

func (instance *ProcessInstanceInfo) FindJobs(ctx context.Context) ([]bpmn_engine_store.IJob, error) {
	return instance.engineState.GetStore().FindJobs(ctx, instance.instanceKey)
}

func (instance *ProcessInstanceInfo) FindTimers(ctx context.Context) ([]bpmn_engine_store.ITimer, error) {
	return instance.engineState.GetStore().FindTimers(ctx, instance.instanceKey)
}

func (instance *ProcessInstanceInfo) FindCatchEvents(ctx context.Context) ([]bpmn_engine_store.ICatchEvent, error) {
	return instance.engineState.GetStore().FindCatchEvents(
		ctx,
		instance.instanceKey,
	)
}

func (instance *ProcessInstanceInfo) FindMessageSubscriptions(ctx context.Context) ([]bpmn_engine_store.IMessageSubscription, error) {
	return instance.engineState.GetStore().FindMessageSubscriptions(
		ctx,
		instance.instanceKey,
	)
}

func (instance *ProcessInstanceInfo) GetVariableContext(ctx context.Context) (map[string]interface{}, error) {
	return instance.engineState.GetStore().GetProcessInstanceVariableContext(ctx, instance.instanceKey)
}

func (instance *ProcessInstanceInfo) SetState(ctx context.Context, state process_instance.State) error {
	return instance.engineState.GetStore().SetProcessInstanceState(ctx, instance.instanceKey, state)
}

func (instance *ProcessInstanceInfo) GetState(ctx context.Context) (process_instance.State, error) {
	return instance.engineState.GetStore().GetProcessInstanceState(ctx, instance.instanceKey)
}

func (instance *ProcessInstanceInfo) GetProcessInfo() bpmn_engine_store.IProcessInfo {
	return instance.processInfo
}

func (instance *ProcessInstanceInfo) GetInstanceKey() bpmn_engine_store.IProcessInstanceKey {
	return instance.instanceKey
}

func (instance *ProcessInstanceInfo) SetVariable(ctx context.Context, key string, value interface{}) error {
	return instance.engineState.GetStore().SetProcessInstanceVariable(ctx, instance.instanceKey, key, value)
}

func (instance *ProcessInstanceInfo) GetCreatedAt() time.Time {
	return instance.createdAt
}
