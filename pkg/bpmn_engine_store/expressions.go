package bpmn_engine_store

import (
	"context"
	"github.com/antonmedv/expr"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
	"strings"
)

func evaluateExpression(expression string, variableContext map[string]interface{}) (interface{}, error) {
	expression = strings.TrimSpace(expression)
	expression = strings.TrimPrefix(expression, "=")
	return expr.Eval(expression, variableContext)
}

func evaluateVariableMapping(ctx context.Context, instance IProcessInstanceInfo, mappings []BPMN20.TIoMapping) error {
	for _, mapping := range mappings {
		variableContext, err := instance.GetVariableContext(ctx)
		if err != nil {
			return err
		}
		evalResult, err := evaluateExpression(mapping.Source, variableContext)
		if err != nil {
			return err
		}
		if err := instance.SetVariable(ctx, mapping.Target, evalResult); err != nil {
			return err
		}
	}
	return nil
}
