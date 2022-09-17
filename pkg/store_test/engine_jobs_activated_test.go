package store_test

import (
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/store"
)

// just to get quick compiler warnings, when interface is not correctly implemented
var _ bpmn_engine_store.IActivatedJob = &store.ActivatedJob{}
