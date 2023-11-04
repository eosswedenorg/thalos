package app

// State represents thalos runtime state
type State struct {
	// Last processed block
	CurrentBlock uint32
}

type (
	// StateLoader is a function that loads a state.
	StateLoader func(*State)

	// StateSaver is a function that saves a state.
	StateSaver func(State) error
)
