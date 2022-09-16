package bpmn_engine_store

import (
	"context"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
)

func (state *BpmnEngineState) handleUserTask(ctx context.Context, instance IProcessInstanceInfo, element BPMN20.TaskElement) bool {
	// TODO consider different handlers, since Service Tasks are different in their definition than user tasks
	return state.handleServiceTask(ctx, instance, element)
}
