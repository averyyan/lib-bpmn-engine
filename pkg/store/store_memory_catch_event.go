package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
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

func (store *EngineMemoryStore) CreateCatchEvent(ctx context.Context, engineState bpmn_engine_store.IBpmnEngine) (bpmn_engine_store.ICatchEvent, error) {
	panic(any("aaa"))
}
