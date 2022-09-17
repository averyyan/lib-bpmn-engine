package store_test

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"testing"
)

func Test_exclusive_gateway_with_expressions_selects_one_and_not_the_other(t *testing.T) {
	// setup
	cp := CallPath{}

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/exclusive-gateway-with-condition.bpmn")
	bpmnEngine.AddTaskHandler("task-a", cp.CallPathHandler)
	bpmnEngine.AddTaskHandler("task-b", cp.CallPathHandler)
	variables := map[string]interface{}{
		"price": -50,
	}

	// when
	bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), variables)

	// then
	then.AssertThat(t, cp.CallPath, is.EqualTo("task-b"))
}

func Test_exclusive_gateway_with_expressions_selects_default(t *testing.T) {
	// setup
	cp := CallPath{}

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/exclusive-gateway-with-condition-and-default.bpmn")
	bpmnEngine.AddTaskHandler("task-a", cp.CallPathHandler)
	bpmnEngine.AddTaskHandler("task-b", cp.CallPathHandler)
	variables := map[string]interface{}{
		"price": -1,
	}

	// when
	bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), variables)

	// then
	then.AssertThat(t, cp.CallPath, is.EqualTo("task-b"))
}

func Test_boolean_expression_evaluates(t *testing.T) {
	variables := map[string]interface{}{
		"aValue": 3,
	}

	result, err := bpmn_engine_store.EvaluateExpression("aValue > 1", variables)

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, result, is.True())
}

func Test_boolean_expression_with_equalsign_evaluates(t *testing.T) {
	variables := map[string]interface{}{
		"aValue": 3,
	}

	result, err := bpmn_engine_store.EvaluateExpression("= aValue > 1", variables)

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, result, is.True())
}

func Test_mathematical_expression_evaluates(t *testing.T) {
	variables := map[string]interface{}{
		"foo": 3,
		"bar": 7,
		"sum": 10,
	}

	result, err := bpmn_engine_store.EvaluateExpression("sum >= foo + bar", variables)

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, result, is.True())
}

func Test_evaluation_error_percolates_up(t *testing.T) {
	// setup

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/exclusive-gateway-with-condition.bpmn")

	// when
	// don't provide variables, for execution
	_, err := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)

	// then
	then.AssertThat(t, err, is.Not(is.Nil()))
}
