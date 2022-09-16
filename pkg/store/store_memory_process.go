package store

import (
	"context"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
)

func (store *EngineMemoryStore) CreateProcess(ctx context.Context, processKey bpmn_engine_store.IProcessInfoKey, definitions BPMN20.TDefinitions) (bpmn_engine_store.IProcessInfo, error) {
	processInfo := &ProcessInfo{
		definitions: definitions,
		ProcessKey:  processKey,
	}
	store.processes[processKey] = processInfo
	return store.processes[processKey], nil
}
