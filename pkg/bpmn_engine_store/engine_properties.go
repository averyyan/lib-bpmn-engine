package bpmn_engine_store

import (
	"github.com/bwmarrin/snowflake"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine/exporter"
)

type BpmnEngineState struct {
	name      string
	handlers  map[string]func(job ActivatedJob)
	snowflake *snowflake.Node
	exporters []exporter.EventExporter
}
