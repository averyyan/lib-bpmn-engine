package bpmn_engine_store

import (
	"context"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/process_instance"
	"time"
)

type IBpmnEngineStore interface {
	//FindProcesses(ctx context.Context) ([]IProcessInfo, error)

	CreateProcess(ctx context.Context, processKey IProcessInfoKey, definitions BPMN20.TDefinitions) (IProcessInfo, error)

	FindProcessInfo(ctx context.Context, processKey IProcessInfoKey) (IProcessInfo, error)

	CreateProcessInstance(ctx context.Context, processKey IProcessInfoKey, processInstanceKey IProcessInstanceKey, variableContext map[string]interface{}) (IProcessInstanceInfo, error)
	FindProcessInstanceInfo(ctx context.Context, processInstanceKey IProcessInstanceKey) (IProcessInstanceInfo, error)
	GetProcessInstanceState(ctx context.Context, processInstanceKey IProcessInstanceKey) (process_instance.State, error)
	SetProcessInstanceState(ctx context.Context, processInstanceKey IProcessInstanceKey, state process_instance.State) error
	SetProcessInstanceVariable(ctx context.Context, processInstanceKey IProcessInstanceKey, key string, value interface{}) error
	GetProcessInstanceVariable(ctx context.Context, processInstanceKey IProcessInstanceKey, key string) (interface{}, error)
	GetProcessInstanceVariableContext(ctx context.Context, processInstanceKey IProcessInstanceKey) (map[string]interface{}, error)

	FindJobs(ctx context.Context, processInstanceKey IProcessInstanceKey) ([]IJob, error)
	FindJob(ctx context.Context, processInstanceKey IProcessInstanceKey, jobKey IJobKey) (IJob, error)

	CreateJob(ctx context.Context, engineStore IBpmnEngineStore, processInstanceKey IProcessInstanceKey, elementId string, elementInstanceKey IJobKey) (IJob, error)

	GetJobState(ctx context.Context, processInstanceKey IProcessInstanceKey, jobKey IJobKey) (activity.LifecycleState, error)
	SetJobState(ctx context.Context, processInstanceKey IProcessInstanceKey, jobKey IJobKey, state activity.LifecycleState, reason string) error

	CreateActivatedJob(job IJob) IActivatedJob

	CreateCatchEvent(ctx context.Context, engineStore IBpmnEngineStore, processInstanceKey IProcessInstanceKey, catchEventKey ICatchEventKey, messageName string, variables map[string]interface{}) (ICatchEvent, error)
	SetCatchEventConsumed(ctx context.Context, processInstanceKey IProcessInstanceKey, catchEventKey ICatchEventKey, consumed bool) error
	GetCatchEventConsumed(ctx context.Context, processInstanceKey IProcessInstanceKey, catchEventKey ICatchEventKey) (bool, error)
	GetCatchEventVariables(ctx context.Context, processInstanceKey IProcessInstanceKey, catchEventKey ICatchEventKey) (map[string]interface{}, error)
	FindCatchEvents(ctx context.Context, processInstanceKey IProcessInstanceKey) ([]ICatchEvent, error)

	CreateMessageSubscription(ctx context.Context, engineStore IBpmnEngineStore, processInstanceKey IProcessInstanceKey, elementInstanceKey IMessageSubscriptionKey, ice BPMN20.TIntermediateCatchEvent) (IMessageSubscription, error)
	FindMessageSubscriptions(ctx context.Context, processInstanceKey IProcessInstanceKey) ([]IMessageSubscription, error)
	GetMessageSubscriptionState(ctx context.Context, processInstanceKey IProcessInstanceKey, messageSubscriptionKey IMessageSubscriptionKey) (activity.LifecycleState, error)
	SetMessageSubscriptionState(ctx context.Context, processInstanceKey IProcessInstanceKey, messageSubscriptionKey IMessageSubscriptionKey, state activity.LifecycleState) error

	CreateTimer(ctx context.Context, engineStore IBpmnEngineStore, processInstanceKey IProcessInstanceKey, timerKey ITimerKey, ice BPMN20.TIntermediateCatchEvent) (ITimer, error)
	FindTimers(ctx context.Context, processInstanceKey IProcessInstanceKey) ([]ITimer, error)
	GetTimerState(ctx context.Context, processInstanceKey IProcessInstanceKey, timerKey ITimerKey) (bpmn_engine.TimerState, error)
	SetTimerState(ctx context.Context, processInstanceKey IProcessInstanceKey, timerKey ITimerKey, state bpmn_engine.TimerState) error

	IsExistScheduledFlow(ctx context.Context, processInstanceKey IProcessInstanceKey, flowId IScheduledFlowKey) (bool, error)
	CreateScheduledFlow(ctx context.Context, processInstanceKey IProcessInstanceKey, flowId IScheduledFlowKey) error
	RemoveScheduledFlow(ctx context.Context, processInstanceKey IProcessInstanceKey, flowId IScheduledFlowKey) error
}

type IProcessInfo interface {
	GetProcessKey() IProcessInfoKey
	GetDefinitions() BPMN20.TDefinitions
	GetBpmnProcessId() string
	GetVersion() int32
}

type IProcessInstanceInfo interface {
	// base
	GenerateKey() int64

	FindJobs(ctx context.Context) ([]IJob, error)
	CreateJob(ctx context.Context, elementId string) (IJob, error)
	SetJobState(ctx context.Context, jobKey IJobKey, state activity.LifecycleState, reason string) error
	GetJobState(ctx context.Context, jobKey IJobKey) (activity.LifecycleState, error)

	CreateActivatedJob(job IJob) IActivatedJob

	CreateScheduledFlow(ctx context.Context, flowId IScheduledFlowKey) error
	RemoveScheduledFlow(ctx context.Context, flowId IScheduledFlowKey) error
	IsExistScheduledFlow(ctx context.Context, flowId IScheduledFlowKey) (bool, error)

	CreateMessageSubscription(ctx context.Context, ice BPMN20.TIntermediateCatchEvent) (IMessageSubscription, error)
	FindMessageSubscriptions(ctx context.Context) ([]IMessageSubscription, error)

	CreateTimer(ctx context.Context, ice BPMN20.TIntermediateCatchEvent) (ITimer, error)
	FindTimers(ctx context.Context) ([]ITimer, error)

	CreateCatchEvent(ctx context.Context, messageName string, variables map[string]interface{}) (ICatchEvent, error)
	FindCatchEvents(ctx context.Context) ([]ICatchEvent, error)

	GetInstanceKey() IProcessInstanceKey
	GetProcessInfo() IProcessInfo
	GetVariableContext(ctx context.Context) (map[string]interface{}, error)
	GetState(ctx context.Context) (process_instance.State, error)
	SetState(ctx context.Context, state process_instance.State) error
	GetVariable(ctx context.Context, key string) (interface{}, error)
	SetVariable(ctx context.Context, key string, value interface{}) error
	GetCreatedAt() time.Time
}

type ICatchEvent interface {
	SetConsumed(ctx context.Context, consumed bool) error
	GetConsumed(ctx context.Context) (bool, error)
	GetVariables(ctx context.Context) (map[string]interface{}, error)
	GetName() string
}

type IJob interface {
	GetElementId() string

	GetVariable(ctx context.Context, key string) (interface{}, error)

	SetVariable(ctx context.Context, key string, value interface{}) error

	GetElementInstanceKey() IJobKey

	SetState(ctx context.Context, state activity.LifecycleState, reason string) error
	GetState(ctx context.Context) (activity.LifecycleState, error)
}

type IActivatedJob interface {
	//IProcessInstanceInfo

	// GetKey the key, a unique identifier for the job
	//GetKey() int64

	// GetProcessInstanceKey the job's process instance key
	//GetProcessInstanceKey() int64

	// GetBpmnProcessId Retrieve id of the job process definition
	//GetBpmnProcessId() string

	// GetProcessDefinitionVersion Retrieve version of the job process definition
	//GetProcessDefinitionVersion() int32

	// GetProcessDefinitionKey Retrieve key of the job process definition
	//GetProcessDefinitionKey() int64

	GetVariable(ctx context.Context, key string) (interface{}, error)
	SetVariable(ctx context.Context, key string, value interface{}) error

	// GetElementId Get element id of the job
	GetElementId() string

	// Fail does set the state the worker missed completing the job
	// Fail and Complete mutual exclude each other
	Fail(ctx context.Context, reason string) error

	// Complete does set the state the worker successfully completing the job
	// Fail and Complete mutual exclude each other
	Complete(ctx context.Context) error
}

type IMessageSubscription interface {
	GetElementId() string
	GetName() string

	GetProcessInstanceKey() IProcessInstanceKey

	SetState(ctx context.Context, state activity.LifecycleState) error
	GetState(ctx context.Context) (activity.LifecycleState, error)
}

type ITimer interface {
	GetElementId() string
	GetProcessInstanceKey() IProcessInstanceKey
	GetDueAt() time.Time
	GetState(ctx context.Context) (bpmn_engine.TimerState, error)
	SetState(ctx context.Context, state bpmn_engine.TimerState) error
}

type IScheduledFlow interface {
}
