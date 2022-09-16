package bpmn_engine_store

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
	"io/ioutil"
)

// LoadFromFile loads a given BPMN file by filename into the engine
// and returns ProcessInfo details for the deployed workflow
func (state *BpmnEngineState) LoadFromFile(ctx context.Context, filename string) (IProcessInfo, error) {
	xmlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return state.load(ctx, xmlData, filename)
}

// LoadFromBytes loads a given BPMN file by xmlData byte array into the engine
// and returns ProcessInfo details for the deployed workflow
func (state *BpmnEngineState) LoadFromBytes(ctx context.Context, xmlData []byte) (IProcessInfo, error) {
	return state.load(ctx, xmlData, "")
}

func (state *BpmnEngineState) load(ctx context.Context, xmlData []byte, resourceName string) (IProcessInfo, error) {
	md5sum := md5.Sum(xmlData)
	var definitions BPMN20.TDefinitions
	err := xml.Unmarshal(xmlData, &definitions)
	if err != nil {
		return nil, err
	}
	processInfo, err := state.GetStore().CreateProcess(ctx, IProcessInfoKey(state.generateKey()), definitions)
	if err != nil {
		return nil, err
	}
	state.exportNewProcessEvent(processInfo, xmlData, resourceName, hex.EncodeToString(md5sum[:]))
	return processInfo, nil
}
