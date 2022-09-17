package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/process_instance"
	"time"
)

func (store *EngineMemoryStore) CreateProcessInstance(ctx context.Context, processKey bpmn_engine_store.IProcessInfoKey, processInstanceKey bpmn_engine_store.IProcessInstanceKey, variableContext map[string]interface{}) (bpmn_engine_store.IProcessInstanceInfo, error) {
	if variableContext == nil {
		variableContext = map[string]interface{}{}
	}
	process, err := store.FindProcessInfo(ctx, processKey)
	if err != nil {
		return nil, err
	}
	snowflakeIdGenerator := bpmn_engine_store.InitializeSnowflakeIdGenerator()
	processInstanceInfo := &ProcessInstanceInfo{
		engineStore:     store,
		processInfo:     process,
		instanceKey:     processInstanceKey,
		variableContext: variableContext,
		createdAt:       time.Now(),
		snowflake:       snowflakeIdGenerator,
		state:           process_instance.READY,
	}
	store.processInstances[processInstanceInfo.instanceKey] = processInstanceInfo

	return store.processInstances[processInstanceInfo.instanceKey], nil
}

func (store *EngineMemoryStore) GetProcessInstanceVariableContext(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey) (map[string]interface{}, error) {
	instance, err := store.findProcessInstance(processInstanceKey)
	if err != nil {
		return nil, err
	}
	return instance.variableContext, nil
}

func (store *EngineMemoryStore) GetProcessInstanceState(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey) (process_instance.State, error) {
	instance, err := store.findProcessInstance(processInstanceKey)
	if err != nil {
		return "", err
	}
	return instance.state, nil
}

func (store *EngineMemoryStore) SetProcessInstanceState(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, state process_instance.State) error {
	instance, err := store.findProcessInstance(processInstanceKey)
	if err != nil {
		return err
	}
	instance.state = state
	return nil
}

func (store *EngineMemoryStore) GetProcessInstanceVariable(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, key string) (interface{}, error) {
	instance, err := store.findProcessInstance(processInstanceKey)
	if err != nil {
		return nil, err
	}
	return instance.variableContext[key], nil
}

func (store *EngineMemoryStore) SetProcessInstanceVariable(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, key string, value interface{}) error {
	instance, err := store.findProcessInstance(processInstanceKey)
	if err != nil {
		return err
	}
	instance.variableContext[key] = value
	return nil
}

func (store *EngineMemoryStore) findProcessInstance(processInstanceKey bpmn_engine_store.IProcessInstanceKey) (*ProcessInstanceInfo, error) {
	instance, isExist := store.processInstances[processInstanceKey]
	if !isExist {
		return nil, errors.New(fmt.Sprintf("ProcessInstances: ProcessInstanceKey=%d  is not exist", processInstanceKey))
	}
	return instance, nil
}
