package gossm

import (
	"errors"
	"fmt"
)

/*
Initialize StateMachine with transition functions, events

	onExecuted Handler is called after a transition is executed and before the next state is entered.
		oldState is the state before the transition or empty if the state machine not transitioned.
*/
func Initialize(sm *StateMachine,
	transitions map[string]TansitionFn,
	events map[string]EnterFunc,
	onExecuted func(sm StateMachine, instance Instance, oldState string) error) error {
	if sm == nil {

		return errors.New("state machine is required")
	}
	if sm.Name == "" {
		return errors.New("state machine name is required")
	}
	if len(sm.States) < 2 {
		return errors.New("state machine should have at least start and end states")
	}

	if len(sm.Transitions) < 1 {
		return errors.New("state machine should have at least one transition")
	}

	if sm.initialized {
		return errors.New("state machine already initialized")
	}
	sm.stateMap = make(map[string]*State)

	// check if states are unique
	for i := range sm.States {
		state := &sm.States[i]
		_, ok := sm.stateMap[state.Name]
		if ok {
			return errors.New(fmt.Sprint("state ", state.Name, " already exists in ", sm.Name))
		}
		sm.stateMap[state.Name] = state
		if (state.Start && sm.startState != nil) || (state.End && sm.endState != nil) {
			return errors.New("start and end states should be unique")
		}
		if state.Start {
			sm.startState = state
		}
		if state.End {
			sm.endState = state
		}

	}
	//map transition functions to states
	for stateFrom, transition := range transitions {
		state, err := sm.getStateByName(stateFrom)
		if err != nil {
			return err
		}
		if state.transitionFn != nil {
			return errors.New(fmt.Sprint("state ", state.Name, " already has a transition function"))
		}
		state.transitionFn = transition
	}

	//map events to states
	for state, onEnterFn := range events {
		state, err := sm.getStateByName(state)
		if err != nil {
			return err
		}
		if state.onEnterFn != nil {
			return errors.New(fmt.Sprint("state ", state.Name, " already has a transition function"))
		}
		state.onEnterFn = onEnterFn
	}

	// check if all states have transition functions
	for _, item := range sm.States {

		if (item.transitionFn == nil) && !item.End {
			return errors.New(fmt.Sprint("state ", item.Name, " should have a transition function in ", sm.Name))
		}
	}
	sm.onExecuted = onExecuted
	sm.initialized = true
	return nil
}

func (sm *StateMachine) getStateByName(name string) (*State, error) {
	var state, exists = sm.stateMap[name]
	if !exists {
		return nil, errors.New(fmt.Sprint("state ", name, " not found in ", sm.Name))
	}
	return state, nil
}

// Creates a new instance of the state machine that holds context. That can be executed and tracked.
func (sm *StateMachine) NewInstance(instanceId any) *Instance {
	return &Instance{
		sm:           sm,
		InstanceId:   instanceId,
		currentState: sm.startState,
		context:      make(map[string]any),
	}
}

// Load an existing instance of the state machine.
func (sm *StateMachine) LoadInstance(instanceId any, context *map[string]any, currentState string) (*Instance, error) {
	if (*context) == nil {
		(*context) = make(map[string]any)
	}
	var state, err = sm.getStateByName(currentState)
	if err != nil {
		return nil, err
	}
	return &Instance{
		sm:           sm,
		InstanceId:   instanceId,
		currentState: state,
		context:      *context,
	}, nil
}
