package gossm

import (
	"errors"
	"fmt"
)

// Get data from context
func (instance *Instance) GetData(key string) any {
	return instance.context[key]
}

// Set data to context
func (instance *Instance) SetData(key string, value any) {
	instance.context[key] = value
}

// Get context
func (instance *Instance) GetContext() map[string]any {
	return instance.context
}

// instance is in end state
func (instance *Instance) IsEnd() bool {
	return instance.currentState.End
}

// instance is in start state
func (instance *Instance) IsStart() bool {
	return instance.currentState.Start
}

// Get current state
func (instance *Instance) GetCurrentState() string {
	return instance.currentState.Name
}

// Execute transition function for current state
func (instance *Instance) Execute() error {

	if !instance.sm.initialized {
		return errors.New("state machine not initialized")
	}

	if instance.IsEnd() {
		return errors.New("instance is in end state")
	}

	var transitionFn = instance.currentState.transitionFn
	if transitionFn == nil {
		return errors.New(fmt.Sprint("transition function not found for state", instance.currentState.Name))
	}

	to, err := transitionFn(*instance.sm, *instance)

	if err != nil {
		return err
	}

	if to == "" {
		if instance.sm.onExecuted != nil {
			err = instance.sm.onExecuted(*instance.sm, *instance, "")
			if err != nil {
				return err
			}
		}
		return nil
	}

	oldState := instance.currentState.Name

	_, err = instance.sm.getTransitionConstraint(instance.currentState.Name, to)

	if err != nil {
		return err
	}

	nextState, err := instance.sm.getStateByName(to)
	if err != nil {
		return err
	}
	instance.currentState = nextState

	if instance.sm.onExecuted != nil {
		err = instance.sm.onExecuted(*instance.sm, *instance, oldState)
		if err != nil {
			return err
		}
	}
	if instance.currentState.onEnterFn != nil {

		err := instance.currentState.onEnterFn(*instance.sm, *instance)

		if err != nil {
			return err
		}
	}

	return nil
}

func (sm *StateMachine) getTransitionConstraint(from string, to string) (TransitionConstraint, error) {
	for _, transition := range sm.Transitions {
		if transition.From == from && transition.To == to {
			return transition, nil
		}
	}
	str := fmt.Sprint("transition from ", from, " to ", to, " not found in ", sm.Name)
	return TransitionConstraint{}, errors.New(str)
}
