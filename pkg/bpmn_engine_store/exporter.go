package bpmn_engine_store

import (
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine/exporter"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
)

// AddEventExporter registers an EventExporter instance
func (state *BpmnEngineState) AddEventExporter(exporter exporter.EventExporter) {
	state.exporters = append(state.exporters, exporter)
}

func (state *BpmnEngineState) exportNewProcessEvent(process IProcessInfo, xmlData []byte, resourceName string, checksum string) {
	event := exporter.ProcessEvent{
		ProcessId:    process.GetBpmnProcessId(),
		ProcessKey:   int64(process.GetProcessKey()),
		Version:      process.GetVersion(),
		XmlData:      xmlData,
		ResourceName: resourceName,
		Checksum:     checksum,
	}
	for _, exp := range state.exporters {
		exp.NewProcessEvent(&event)
	}
}

func (state *BpmnEngineState) exportEndProcessEvent(process IProcessInfo, processInstance IProcessInstanceInfo) {
	event := exporter.ProcessInstanceEvent{
		ProcessId:          process.GetBpmnProcessId(),
		ProcessKey:         int64(process.GetProcessKey()),
		Version:            process.GetVersion(),
		ProcessInstanceKey: int64(processInstance.GetInstanceKey()),
	}
	for _, exp := range state.exporters {
		exp.EndProcessEvent(&event)
	}
}

func (state *BpmnEngineState) exportProcessInstanceEvent(process IProcessInfo, processInstance IProcessInstanceInfo) {
	event := exporter.ProcessInstanceEvent{
		ProcessId:          process.GetBpmnProcessId(),
		ProcessKey:         int64(process.GetProcessKey()),
		Version:            process.GetVersion(),
		ProcessInstanceKey: int64(processInstance.GetInstanceKey()),
	}
	for _, exp := range state.exporters {
		exp.NewProcessInstanceEvent(&event)
	}
}

func (state *BpmnEngineState) exportElementEvent(process IProcessInfo, processInstance IProcessInstanceInfo, element BPMN20.BaseElement, intent exporter.Intent) {
	event := exporter.ProcessInstanceEvent{
		ProcessId:          process.GetBpmnProcessId(),
		ProcessKey:         int64(process.GetProcessKey()),
		Version:            process.GetVersion(),
		ProcessInstanceKey: int64(processInstance.GetInstanceKey()),
	}
	info := exporter.ElementInfo{
		BpmnElementType: string(element.GetType()),
		ElementId:       element.GetId(),
		Intent:          string(intent),
	}
	for _, exp := range state.exporters {
		exp.NewElementEvent(&event, &info)
	}
}

func (state *BpmnEngineState) exportSequenceFlowEvent(process IProcessInfo, processInstance IProcessInstanceInfo, flow BPMN20.TSequenceFlow) {
	event := exporter.ProcessInstanceEvent{
		ProcessId:          process.GetBpmnProcessId(),
		ProcessKey:         int64(process.GetProcessKey()),
		Version:            process.GetVersion(),
		ProcessInstanceKey: int64(processInstance.GetInstanceKey()),
	}
	info := exporter.ElementInfo{
		BpmnElementType: string(BPMN20.SequenceFlow),
		ElementId:       flow.Id,
		Intent:          string(exporter.SequenceFlowTaken),
	}
	for _, exp := range state.exporters {
		exp.NewElementEvent(&event, &info)
	}
}
