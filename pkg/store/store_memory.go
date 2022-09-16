package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
)

type ScheduledFlow string

func NewMemoryStore() bpmn_engine_store.IBpmnEngineStore {
	return &EngineMemoryStore{
		processes:            map[bpmn_engine_store.IProcessInfoKey]*ProcessInfo{},
		processInstances:     map[bpmn_engine_store.IProcessInstanceKey]*ProcessInstanceInfo{},
		catchEvents:          map[bpmn_engine_store.IProcessInstanceKey]map[bpmn_engine_store.ICatchEventKey]*catchEvent{},
		messageSubscriptions: map[bpmn_engine_store.IProcessInstanceKey]map[bpmn_engine_store.IMessageSubscriptionKey]*MessageSubscription{},
		jobs:                 map[bpmn_engine_store.IProcessInstanceKey]map[bpmn_engine_store.IJobKey]*job{},
		timers:               map[bpmn_engine_store.IProcessInstanceKey]map[bpmn_engine_store.ITimerKey]*Timer{},
		scheduledFlows:       map[bpmn_engine_store.IProcessInstanceKey]map[bpmn_engine_store.IScheduledFlowKey]int{},
	}
}

type EngineMemoryStore struct {
	processes            map[bpmn_engine_store.IProcessInfoKey]*ProcessInfo
	processInstances     map[bpmn_engine_store.IProcessInstanceKey]*ProcessInstanceInfo
	catchEvents          map[bpmn_engine_store.IProcessInstanceKey]map[bpmn_engine_store.ICatchEventKey]*catchEvent
	messageSubscriptions map[bpmn_engine_store.IProcessInstanceKey]map[bpmn_engine_store.IMessageSubscriptionKey]*MessageSubscription
	jobs                 map[bpmn_engine_store.IProcessInstanceKey]map[bpmn_engine_store.IJobKey]*job
	timers               map[bpmn_engine_store.IProcessInstanceKey]map[bpmn_engine_store.ITimerKey]*Timer
	scheduledFlows       map[bpmn_engine_store.IProcessInstanceKey]map[bpmn_engine_store.IScheduledFlowKey]int
}

func (store *EngineMemoryStore) FindCatchEvents(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey) ([]bpmn_engine_store.ICatchEvent, error) {
	catchEvents, isExistCatchEvents := store.catchEvents[processInstanceKey]
	if !isExistCatchEvents {
		return nil, errors.New(fmt.Sprintf("catchEvents: ProcessInstanceKey=%d  is not exist", processInstanceKey))
	}
	var result []bpmn_engine_store.ICatchEvent
	for _, catchEvent := range catchEvents {
		result = append(result, catchEvent)
	}
	return result, nil
}

func (store *EngineMemoryStore) FindMessageSubscriptions(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey) ([]bpmn_engine_store.IMessageSubscription, error) {
	messageSubscriptions, isExistMessageSubscriptions := store.messageSubscriptions[processInstanceKey]
	if !isExistMessageSubscriptions {
		return nil, errors.New(fmt.Sprintf("messageSubscriptions: ProcessInstanceKey=%d  is not exist", processInstanceKey))
	}
	var result []bpmn_engine_store.IMessageSubscription
	for _, subscription := range messageSubscriptions {
		result = append(result, subscription)
	}
	return result, nil
}

func (store *EngineMemoryStore) FindProcesses(ctx context.Context) ([]bpmn_engine_store.IProcessInfo, error) {
	var processes []bpmn_engine_store.IProcessInfo
	for _, process := range store.processes {
		processes = append(processes, process)
	}
	return processes, nil
}

func (store *EngineMemoryStore) FindProcessInfo(ctx context.Context, processKey bpmn_engine_store.IProcessInfoKey) (bpmn_engine_store.IProcessInfo, error) {
	process, isExist := store.processes[processKey]
	if !isExist {
		return nil, errors.New(fmt.Sprintf("ProcessInfo ProcessInfoKey=%d is not exist", processKey))
	}
	return process, nil
}

func (store *EngineMemoryStore) FindProcessInstanceInfo(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey) (bpmn_engine_store.IProcessInstanceInfo, error) {
	instances, isExist := store.processInstances[processInstanceKey]
	if !isExist {
		return nil, errors.New(fmt.Sprintf("processInstanceInfo key=%d is not exist", processInstanceKey))
	}
	return instances, nil
}
