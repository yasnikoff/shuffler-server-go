package accessors

import (
	"testing"
	"time"
	"encoding/json"
)

func TestJson(t *testing.T) {
	t1 := Time{value:time.Now()}
	t2 := &Time{}
	data, err := json.Marshal(t1)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, t2); err != nil {
		t.Fatal(err)
	}
	if t1.value != t2.value {
		t.Errorf("json.Unmarshal() called on marshalled data returned result that differs from initial. Got: %v, Want: %v", t2.value, t1.value)
	}

}
func TestJsonPointer(t *testing.T) {
	t1 := &Time{value:time.Now()}
	t2 := &Time{}
	data, err := json.Marshal(t1)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, t2); err != nil {
		t.Fatal(err)
	}
	if t1.value != t2.value {
		t.Errorf("json.Unmarshal() called on marshalled data returned result that differs from initial. Got: %v, Want: %v", t2.value, t1.value)
	}

}
