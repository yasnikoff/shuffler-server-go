package id

import (
	"testing"
	"encoding/json"
)

func TestJSON(t *testing.T) {
	t.Parallel()
	i:=New()
	data, err:=json.Marshal(i)
	if err!=nil{
		t.Fatalf("can't marshal ID %q: %s", i, err)
	}
	var idm ID
	if err:=json.Unmarshal(data, &idm);err!=nil{
		t.Fatalf("can't unmarshal data %s into id: %s", data, err)
	}
	if i!=idm{
		t.Errorf("unmarshalled ID doesn't equal original one. Got: %s, want: %s", idm, i)
	}

}

