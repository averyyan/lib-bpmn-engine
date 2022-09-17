package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"time"
)

func (store *EngineMemoryStore) GetCatchEventVariables(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, catchEventKey bpmn_engine_store.ICatchEventKey) (map[string]interface{}, error) {
	event, err := store.findCatchEvent(processInstanceKey, catchEventKey)
	if err != nil {
		return nil, err
	}
	return event.variables, err
}

func (store *EngineMemoryStore) findCatchEvents(processInstanceKey bpmn_engine_store.IProcessInstanceKey) (map[bpmn_engine_store.ICatchEventKey]*catchEvent, error) {
	catchEvents, isExist := store.catchEvents[processInstanceKey]
	if !isExist {
		return nil, errors.New(fmt.Sprintf("catchEvents: ProcessInstanceKey=%d  is not exist", processInstanceKey))
	}
	return catchEvents, nil
}

func (store *EngineMemoryStore) findCatchEvent(processInstanceKey bpmn_engine_store.IProcessInstanceKey, catchEventKey bpmn_engine_store.ICatchEventKey) (*catchEvent, error) {
	catchEvents, err := store.findCatchEvents(processInstanceKey)
	if err != nil {
		return nil, err
	}
	catchEvent, isExist := catchEvents[catchEventKey]
	if !isExist {
		return nil, errors.New(fmt.Sprintf("catchEvents: ProcessInstanceKey=%d && CatchEventKey=%d is not exist", processInstanceKey, catchEventKey))
	}
	return catchEvent, nil
}

func (store *EngineMemoryStore) GetCatchEventConsumed(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, catchEventKey bpmn_engine_store.ICatchEventKey) (bool, error) {
	event, err := store.findCatchEvent(processInstanceKey, catchEventKey)
	if err != nil {
		return false, err
	}
	return event.isConsumed, nil
}

func (store *EngineMemoryStore) SetCatchEventConsumed(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, catchEventKey bpmn_engine_store.ICatchEventKey, consumed bool) error {
	event, err := store.findCatchEvent(processInstanceKey, catchEventKey)
	if err != nil {
		return err
	}
	event.isConsumed = consumed
	return nil
}

func (store *EngineMemoryStore) CreateCatchEvent(ctx context.Context, engineStore bpmn_engine_store.IBpmnEngineStore, processInstanceKey bpmn_engine_store.IProcessInstanceKey, catchEventKey bpmn_engine_store.ICatchEventKey, messageName string, variables map[string]interface{}) (bpmn_engine_store.ICatchEvent, error) {
	if store.catchEvents[processInstanceKey] == nil {
		store.catchEvents[processInstanceKey] = map[bpmn_engine_store.ICatchEventKey]*catchEvent{}
	}
	event := &catchEvent{
		engineStore:        engineStore,
		ProcessInstanceKey: processInstanceKey,
		ElementInstanceKey: catchEventKey,
		caughtAt:           time.Now(),
		name:               messageName,
		variables:          variables,
		isConsumed:         false,
	}
	store.catchEvents[processInstanceKey][catchEventKey] = event
	return store.catchEvents[processInstanceKey][catchEventKey], nil
}
