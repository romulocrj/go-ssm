package gossm

/*
Transition function should return the next state name or empty string to stay in the same state
*/
type TansitionFn func(sm StateMachine, instance Instance) (string, error)

/*
EnterFunc is a function that is called when a state is entered
*/
type EnterFunc func(sm StateMachine, instance Instance) error

type State struct {
	Start        bool   `json:"start"`
	End          bool   `json:"end"`
	Name         string `json:"name"`
	transitionFn TansitionFn
	onEnterFn    EnterFunc
}
type TransitionConstraint struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type StateMachine struct {
	Name        string                 `json:"name"`
	Transitions []TransitionConstraint `json:"transitions"`
	States      []State                `json:"states"`
	initialized bool
	startState  *State
	endState    *State
	onExecuted  func(sm StateMachine, instance Instance, oldState string) error
	stateMap    map[string]*State
}

type Instance struct {
	InstanceId   interface{}
	context      map[string]any
	sm           *StateMachine
	currentState *State
}
