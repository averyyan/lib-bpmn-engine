package store

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/process_instance"
	"time"
)

// ProcessInstanceInfo need save to db processInfo is from ProcessInfo
// this the core of the engine
type ProcessInstanceInfo struct {
	engineStore     bpmn_engine_store.IBpmnEngineStore
	snowflake       *snowflake.Node
	processInfo     bpmn_engine_store.IProcessInfo
	instanceKey     bpmn_engine_store.IProcessInstanceKey
	variableContext map[string]interface{}
	createdAt       time.Time
	state           process_instance.State
	caughtEvents    []catchEvent
}

func (instance *ProcessInstanceInfo) CreateCatchEvent(ctx context.Context, messageName string, variables map[string]interface{}) (bpmn_engine_store.ICatchEvent, error) {
	return instance.engineStore.CreateCatchEvent(ctx, instance.engineStore, instance.instanceKey, bpmn_engine_store.ICatchEventKey(instance.GenerateKey()), messageName, variables)
}

func (instance *ProcessInstanceInfo) CreateTimer(ctx context.Context, ice BPMN20.TIntermediateCatchEvent) (bpmn_engine_store.ITimer, error) {
	return instance.engineStore.CreateTimer(ctx, instance.engineStore, instance.instanceKey, bpmn_engine_store.ITimerKey(instance.GenerateKey()), ice)
}

func (instance *ProcessInstanceInfo) CreateMessageSubscription(ctx context.Context, ice BPMN20.TIntermediateCatchEvent) (bpmn_engine_store.IMessageSubscription, error) {
	return instance.engineStore.CreateMessageSubscription(ctx, instance.engineStore, instance.instanceKey, bpmn_engine_store.IMessageSubscriptionKey(instance.GenerateKey()), ice)
}

func (instance *ProcessInstanceInfo) RemoveScheduledFlow(ctx context.Context, flowId bpmn_engine_store.IScheduledFlowKey) error {
	return instance.engineStore.RemoveScheduledFlow(ctx, instance.GetInstanceKey(), flowId)
}

func (instance *ProcessInstanceInfo) IsExistScheduledFlow(ctx context.Context, flowId bpmn_engine_store.IScheduledFlowKey) (bool, error) {
	return instance.engineStore.IsExistScheduledFlow(ctx, instance.GetInstanceKey(), flowId)
}

func (instance *ProcessInstanceInfo) CreateActivatedJob(job bpmn_engine_store.IJob) bpmn_engine_store.IActivatedJob {
	return instance.engineStore.CreateActivatedJob(job)
}

func (instance *ProcessInstanceInfo) GenerateKey() int64 {
	return instance.snowflake.Generate().Int64()
}

func (instance *ProcessInstanceInfo) CreateScheduledFlow(ctx context.Context, flowId bpmn_engine_store.IScheduledFlowKey) error {
	return instance.engineStore.CreateScheduledFlow(ctx, instance.instanceKey, flowId)
}

func (instance *ProcessInstanceInfo) GetVariable(ctx context.Context, key string) (interface{}, error) {
	return instance.engineStore.GetProcessInstanceVariable(ctx, instance.instanceKey, key)
}

func (instance *ProcessInstanceInfo) CreateJob(ctx context.Context, elementId string) (bpmn_engine_store.IJob, error) {
	return instance.engineStore.CreateJob(ctx, instance.engineStore, instance.instanceKey, elementId, bpmn_engine_store.IJobKey(instance.GenerateKey()))
}

func (instance *ProcessInstanceInfo) GetJobState(ctx context.Context, jobKey bpmn_engine_store.IJobKey) (activity.LifecycleState, error) {
	return instance.engineStore.GetJobState(ctx, instance.instanceKey, jobKey)
}

func (instance *ProcessInstanceInfo) SetJobState(ctx context.Context, jobKey bpmn_engine_store.IJobKey, state activity.LifecycleState, reason string) error {
	return instance.engineStore.SetJobState(ctx, instance.instanceKey, jobKey, state, reason)
}

func (instance *ProcessInstanceInfo) FindJobs(ctx context.Context) ([]bpmn_engine_store.IJob, error) {
	return instance.engineStore.FindJobs(ctx, instance.instanceKey)
}

func (instance *ProcessInstanceInfo) FindTimers(ctx context.Context) ([]bpmn_engine_store.ITimer, error) {
	return instance.engineStore.FindTimers(ctx, instance.instanceKey)
}

func (instance *ProcessInstanceInfo) FindCatchEvents(ctx context.Context) ([]bpmn_engine_store.ICatchEvent, error) {
	return instance.engineStore.FindCatchEvents(
		ctx,
		instance.instanceKey,
	)
}

func (instance *ProcessInstanceInfo) FindMessageSubscriptions(ctx context.Context) ([]bpmn_engine_store.IMessageSubscription, error) {
	return instance.engineStore.FindMessageSubscriptions(
		ctx,
		instance.instanceKey,
	)
}

func (instance *ProcessInstanceInfo) GetVariableContext(ctx context.Context) (map[string]interface{}, error) {
	return instance.engineStore.GetProcessInstanceVariableContext(ctx, instance.instanceKey)
}

func (instance *ProcessInstanceInfo) SetState(ctx context.Context, state process_instance.State) error {
	return instance.engineStore.SetProcessInstanceState(ctx, instance.instanceKey, state)
}

func (instance *ProcessInstanceInfo) GetState(ctx context.Context) (process_instance.State, error) {
	return instance.engineStore.GetProcessInstanceState(ctx, instance.instanceKey)
}

func (instance *ProcessInstanceInfo) GetProcessInfo() bpmn_engine_store.IProcessInfo {
	return instance.processInfo
}

func (instance *ProcessInstanceInfo) GetInstanceKey() bpmn_engine_store.IProcessInstanceKey {
	return instance.instanceKey
}

func (instance *ProcessInstanceInfo) SetVariable(ctx context.Context, key string, value interface{}) error {
	return instance.engineStore.SetProcessInstanceVariable(ctx, instance.instanceKey, key, value)
}

func (instance *ProcessInstanceInfo) GetCreatedAt() time.Time {
	return instance.createdAt
}
