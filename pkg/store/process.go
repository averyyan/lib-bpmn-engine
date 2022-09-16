package store

import (
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
)

// this is a template for db model

// ProcessInfo can save to db or read from file
type ProcessInfo struct {
	BpmnProcessId string                            // The ID as defined in the BPMN file
	Version       int32                             // A version of the process, default=1, incremented, when another process with the same ID is loaded
	ProcessKey    bpmn_engine_store.IProcessInfoKey // The engines key for this given process with version
	definitions   BPMN20.TDefinitions               // parsed file content
	checksumBytes [16]byte                          // internal checksum to identify different versions
}

func (p *ProcessInfo) GetVersion() int32 {
	return p.Version
}

func (p *ProcessInfo) GetBpmnProcessId() string {
	return p.BpmnProcessId
}

func (p *ProcessInfo) GetProcessKey() bpmn_engine_store.IProcessInfoKey {
	return p.ProcessKey
}

func (p *ProcessInfo) GetDefinitions() BPMN20.TDefinitions {
	return p.definitions
}
