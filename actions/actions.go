package actions

import (
	"encoding/json"
	"fmt"
)

var types map[string]ActionConstructor

type ActionConstructor func() Action

func init(){
	types = make(map[string]ActionConstructor)
}

// Register must be called for an action before it can be used.
// Use unique typeName and a constructor that returns new Action.
// Json data will be unmarshalled to the Action instance.
func Register(typeName string, c ActionConstructor) {
	if _, ok := types[typeName]; ok {
		panic(fmt.Errorf("Action type %q already registerred.", typeName))
	}
	types[typeName] = c
}

// Action represents command to be run in separate process.
// Command should not take much time to execute.
// They are meant to be used for short setup scripts.
// Type is the name of action under which it was registered by Register().
type Action interface {
	Type() string
	Run() (output []byte, err error)
}

// actionData is used in json marshalling - unmarshalling
type actionData struct {
	TypeName string `json:"type"`
	Data     json.RawMessage `json:"data"`
}
