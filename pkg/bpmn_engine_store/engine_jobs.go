package bpmn_engine_store

import "time"

type ActivatedJob struct {
	//processInstanceInfo *ProcessInstanceInfo
	completeHandler func()
	failHandler     func(reason string)

	// the key, a unique identifier for the job
	Key int64
	// the job's process instance key
	ProcessInstanceKey int64
	// the bpmn process ID of the job process definition
	BpmnProcessId string
	// the version of the job process definition
	ProcessDefinitionVersion int32
	// the key of the job process definition
	ProcessDefinitionKey int64
	// the associated task element ID
	ElementId string
	// when the job was created
	CreatedAt time.Time
}
