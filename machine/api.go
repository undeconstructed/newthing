package machine

import "net/http"

// Operation is a thing for running.
type Operation func(command Command) error

// Operations is named operation.
type Operations map[string]Operation

// Resource is something identifiable that has action.
type Resource struct {
	path     pathex
	Guards   Guards
	Children Resources
	Actions  Actions
}

// Resources is paths to resources.
type Resources map[string]Resource

// Guard is a thing that guards a resource.
type Guard struct {
}

// Guards is multiple guards.
type Guards []Guard

// Action is something invokable on a resource.
type Action struct {
	Handler   Handler
	Getter    Getter
	Commander Commander
}

// Actions is names to actions.
type Actions map[string]Action

// Args are pulled from a request for use in handling an action.
type Args struct {
	Vars map[string]string
}

// Response is what actions produce.
type Response struct {
	Status int
	Body   interface{}
}

// Cache gives access to the cache.
type Cache interface {
}

// Store gives access to the backend.
type Store interface {
}

// Command is how changes are applied.
type Command struct {
	id        string
	Category  string
	Key       string
	Operation string
	Message   string
}

// Handler is like the normal http handlers.
type Handler func(http.ResponseWriter, *http.Request)

// Getter is for getting data only on an action.
type Getter func(args Args, store Store) (Response, error)

// Commander is for doing real work on an action.
type Commander func(args Args) (Command, error)
