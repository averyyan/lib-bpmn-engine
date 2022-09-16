package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
)

func (store *EngineMemoryStore) CreateScheduledFlow(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, flowId bpmn_engine_store.IScheduledFlowKey) error {
	if store.scheduledFlows[processInstanceKey] == nil {
		store.scheduledFlows[processInstanceKey] = map[bpmn_engine_store.IScheduledFlowKey]int{}
	}
	store.scheduledFlows[processInstanceKey][flowId] = 1
	return nil
}

func (store *EngineMemoryStore) RemoveScheduledFlow(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, flowId bpmn_engine_store.IScheduledFlowKey) error {
	delete(store.scheduledFlows[processInstanceKey], flowId)
	return nil
}

func (store *EngineMemoryStore) IsExistScheduledFlow(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, flowId bpmn_engine_store.IScheduledFlowKey) (bool, error) {
	flows, err := store.findScheduledFlows(processInstanceKey)
	if err != nil {
		return false, err
	}
	_, isExist := flows[flowId]

	return isExist, nil
}

func (store *EngineMemoryStore) findScheduledFlows(processInstanceKey bpmn_engine_store.IProcessInstanceKey) (map[bpmn_engine_store.IScheduledFlowKey]int, error) {
	flows, isExist := store.scheduledFlows[processInstanceKey]
	if !isExist {
		return nil, errors.New(fmt.Sprintf("scheduledFlows: ProcessInstanceKey=%d  is not exist", processInstanceKey))
	}
	return flows, nil
}
