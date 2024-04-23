package gossm_test

import (
	"encoding/json"
	"log"
	"testing"

	gossm "github.com/romulocrj/go-ssm"
)

var jsonStr = `{
	"name": "StateMachine",
	"states": [
		{"start":true, "name":"Start"},
		{"end":true, "name":"End"},
		{"name":"Welcome"},
		{"name":"Wait data"},
		{"name":"Send data"}
		],
	"transitions":[
	  {"from":"Start", "to":"Welcome"},
	  {"from":"Welcome", "to":"Wait data"},
	  {"from":"Wait data", "to":"Send data"},
	  {"from":"Send data", "to":"End"}
	]
  }`

const (
	StateStart    = "Start"
	StateEnd      = "End"
	StateWelcome  = "Welcome"
	StateWaitData = "Wait data"
	StateSendData = "Send data"
)

func TestEmptyTransition(t *testing.T) {
	sm := gossm.StateMachine{}
	json.Unmarshal([]byte(jsonStr), &sm)

	var transitions = map[string]gossm.TansitionFn{
		StateStart:    func(sm gossm.StateMachine, instance gossm.Instance) (string, error) { return StateWelcome, nil },
		StateWelcome:  func(sm gossm.StateMachine, instance gossm.Instance) (string, error) { return StateWaitData, nil },
		StateWaitData: func(sm gossm.StateMachine, instance gossm.Instance) (string, error) { return "", nil },
		StateSendData: func(sm gossm.StateMachine, instance gossm.Instance) (string, error) { return StateEnd, nil },
	}
	var events = map[string]gossm.EnterFunc{
		StateStart:    onEnterFn,
		StateWelcome:  onEnterFn,
		StateWaitData: onEnterFn,
		StateSendData: onEnterFn,
	}

	err := gossm.Initialize(&sm, transitions, events, onExecuted)
	if err != nil {
		t.Error(err)
		return
	}

	context := make(map[string]any)
	instance, err := sm.LoadInstance("12341234", &context, StateWaitData)

	if err != nil {
		t.Error(err)
		return
	}
	state := instance.GetCurrentState()
	if state != StateWaitData {
		t.Error("Expected Send Wait Data state")
		return
	}

	err = instance.Execute()
	if err != nil {
		t.Error(err)
		return
	}

	state = instance.GetCurrentState()
	if state != StateWaitData {
		t.Error("Expected Send Wait Data state")
		return
	}

}
func TestLoadAndExecuteInstance(t *testing.T) {
	sm := gossm.StateMachine{}
	json.Unmarshal([]byte(jsonStr), &sm)

	var transitions = map[string]gossm.TansitionFn{
		StateStart:    func(sm gossm.StateMachine, instance gossm.Instance) (string, error) { return StateWelcome, nil },
		StateWelcome:  func(sm gossm.StateMachine, instance gossm.Instance) (string, error) { return StateWaitData, nil },
		StateWaitData: func(sm gossm.StateMachine, instance gossm.Instance) (string, error) { return StateSendData, nil },
		StateSendData: func(sm gossm.StateMachine, instance gossm.Instance) (string, error) { return StateEnd, nil },
	}
	var events = map[string]gossm.EnterFunc{
		StateStart:    onEnterFn,
		StateWelcome:  onEnterFn,
		StateWaitData: onEnterFn,
		StateSendData: onEnterFn,
	}

	err := gossm.Initialize(&sm, transitions, events, onExecuted)
	if err != nil {
		t.Error(err)
		return
	}

	context := make(map[string]any)
	instance, err := sm.LoadInstance("12341234", &context, StateSendData)

	if err != nil {
		t.Error(err)
		return
	}
	state := instance.GetCurrentState()
	if state != StateSendData {
		t.Error("Expected Send Data state")
		return
	}

	err = instance.Execute()
	if err != nil {
		t.Error(err)
		return
	}

	if !instance.IsEnd() {
		t.Error("Expected End state")
		return
	}

}

func TestCreateStateMachineAndExecuteInstance(t *testing.T) {
	sm := gossm.StateMachine{}
	json.Unmarshal([]byte(jsonStr), &sm)

	var transitions = map[string]gossm.TansitionFn{
		StateStart:    func(sm gossm.StateMachine, instance gossm.Instance) (string, error) { return StateWelcome, nil },
		StateWelcome:  func(sm gossm.StateMachine, instance gossm.Instance) (string, error) { return StateWaitData, nil },
		StateWaitData: func(sm gossm.StateMachine, instance gossm.Instance) (string, error) { return StateSendData, nil },
		StateSendData: func(sm gossm.StateMachine, instance gossm.Instance) (string, error) { return StateEnd, nil },
	}
	var events = map[string]gossm.EnterFunc{
		StateStart:    onEnterFn,
		StateWelcome:  onEnterFn,
		StateWaitData: onEnterFn,
		StateSendData: onEnterFn,
	}

	err := gossm.Initialize(&sm, transitions, events, onExecuted)
	if err != nil {
		t.Error(err)
		return
	}

	var instance = sm.NewInstance("12341234")
	state := instance.GetCurrentState()
	if state != StateStart {
		t.Error("Expected Start state")
		return
	}

	err = instance.Execute()
	if err != nil {
		t.Error(err)
		return
	}
	state = instance.GetCurrentState()
	if state != StateWelcome {
		t.Error("Expected Welcome state")
		return
	}

	err = instance.Execute()
	if err != nil {
		t.Error(err)
		return
	}
	state = instance.GetCurrentState()
	if state != StateWaitData {
		t.Error("Expected Wait data state")
		return
	}

	err = instance.Execute()
	if err != nil {
		t.Error(err)
		return
	}
	state = instance.GetCurrentState()
	if state != StateSendData {
		t.Error("Expected Send data state")
		return
	}

	err = instance.Execute()
	if err != nil {
		t.Error(err)
		return
	}
	state = instance.GetCurrentState()
	if state != StateEnd {
		t.Error("Expected End state")
		return
	}

	err = instance.Execute()
	if err == nil {
		t.Error("Expected error")
		return
	}

}

func onEnterFn(sm gossm.StateMachine, instance gossm.Instance) error {
	log.Println(sm.Name, instance.InstanceId, "onEnter", instance.GetCurrentState())
	return nil
}
func onExecuted(sm gossm.StateMachine, instance gossm.Instance, oldState string) error {
	log.Println(sm.Name, instance.InstanceId, "onExecuted", oldState, "->", instance.GetCurrentState())
	return nil
}
