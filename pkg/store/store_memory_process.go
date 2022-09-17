package store

import (
	"context"
	"crypto/md5"
	"encoding/xml"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
)

func (store *EngineMemoryStore) CreateProcess(
	_ context.Context,
	processKey bpmn_engine_store.IProcessInfoKey,
	definitions BPMN20.TDefinitions,
) (bpmn_engine_store.IProcessInfo, error) {
	marshal, err := xml.Marshal(definitions)
	if err != nil {
		return nil, err
	}
	processInfo := &ProcessInfo{
		BpmnProcessId: definitions.Process.Id,
		Version:       1,
		definitions:   definitions,
		ProcessKey:    processKey,
		checksumBytes: md5.Sum(marshal),
	}
	for _, process := range store.processes {
		if process.BpmnProcessId == processInfo.BpmnProcessId {
			if process.checksumBytes == processInfo.checksumBytes {
				return process, nil
			}
			processInfo.Version += 1
		}
	}
	store.processes[processKey] = processInfo
	return store.processes[processKey], nil
}
