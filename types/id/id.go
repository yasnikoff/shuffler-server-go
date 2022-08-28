package id

import (
	"github.com/docker/dctx/uuid"
	"fmt"
	"encoding/json"
)

var Zero = ID{}

type ID struct {
	uuid.UUID
}

func New() ID {
	return ID{uuid.Generate()}
}

func (i ID)String() string {
	return i.UUID.String()
}

func Parse(s string) (ID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return Zero, fmt.Errorf("can not parse uuid from string %q: %s", s, err)
	}
	return ID{id}, nil
}

func (i ID)MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}
func (i *ID)UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	u, err := uuid.Parse(str)
	if err != nil {
		return err
	}
	i.UUID = u
	return nil
}
