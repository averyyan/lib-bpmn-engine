package store

import "github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"

func (store *EngineMemoryStore) CreateActivatedJob(engineState bpmn_engine_store.IBpmnEngine, job bpmn_engine_store.IJob) bpmn_engine_store.IActivatedJob {
	return &ActivatedJob{
		job: job,
	}
}
