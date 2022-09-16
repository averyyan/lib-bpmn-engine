package bpmn_engine_store

import (
	"github.com/bwmarrin/snowflake"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine/exporter"
)

type BpmnEngineState struct {
	name        string
	engineStore IBpmnEngineStore // store
	handlers    map[string]func(job IActivatedJob)
	snowflake   *snowflake.Node
	exporters   []exporter.EventExporter
}

func (state *BpmnEngineState) AddTaskHandler(taskId string, handler func(job IActivatedJob)) {
	if nil == state.handlers {
		state.handlers = make(map[string]func(job IActivatedJob))
	}
	state.handlers[taskId] = handler
}

func (state *BpmnEngineState) GetStore() IBpmnEngineStore {
	return state.engineStore
}
