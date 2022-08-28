package actions

import (
	"encoding/json"
	"fmt"
)


// Collection is a container for actions.
type Collection struct {
	actions []Action
}

func NewCollection() *Collection {
	return &Collection{}
}

func (c *Collection)Add(a Action) {
	c.actions = append(c.actions, a)
}

func (c *Collection) MarshalJSON() ([]byte, error) {
	result := make([]actionData, len(c.actions))

	for i, a := range c.actions {
		data, err := json.Marshal(a)
		if err != nil {
			return nil, err
		}
		result[i] = actionData{TypeName:a.Type(), Data:data}

	}

	return json.Marshal(result)
}

func (c *Collection)UnmarshalJSON(data []byte) error {
	var result []actionData
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	c.actions = make([]Action, len(result))
	var a Action
	for i, data := range result {
		constructor, ok := types[data.TypeName]
		if !ok {
			return fmt.Errorf("Unkown ation type: %q", data.TypeName)
		}
		a = constructor()
		err := json.Unmarshal(data.Data, a)
		if err != nil {
			return err
		}
		c.actions[i] = a
	}
	return nil
}

// RunAction runs action with index ind in a separate process with working directory set to workDir.
// Returns combined output and error. About combined output see os/exec.CombinedOutput.
func (c *Collection)RunAction(ind int) ([]byte, error) {
	if ind > len(c.actions) - 1 || ind < 0 {
		return []byte{}, fmt.Errorf("actions.Collection.RunAction: Index %d is out of range", ind)
	}
	return c.actions[ind].Run()
}


