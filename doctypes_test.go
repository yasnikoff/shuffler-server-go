package shuffler

import (
	"testing"
	"encoding/json"
	"fmt"
	"reflect"
)

var typesToCheck []*DocType

func init() {

	typesToCheck = []*DocType{DocTypeProject}
}

func TestTypesToCheck(t *testing.T) {
	//t.Parallel()
	if len(typesToCheck) == 0 {
		t.Errorf("doc types didn't register themseves.")
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	for _, dt := range docTypes {
		d := dt.New()
		if d == nil {
			t.Errorf("docType.New() returned nil.")
		}
	}
}

func TestJSON(t *testing.T) {
	t.Parallel()
	for _, dt := range typesToCheck {
		var dtm *DocType
		if dt == nil {
			t.Fatalf("DocType from typesToCheck is nil. typesToCheck: %v", typesToCheck)
		}
		name := dt.Name()
		err := json.Unmarshal([]byte(fmt.Sprintf(`"%s"`, name)), &dtm)
		if err != nil {
			t.Error(err)
			continue
		}
		if dtm == nil {
			t.Errorf("failed to unmarshal data %q. Unmarshalled doctype is nil", name)
			continue
		}
		if dtm.Name() != name {
			t.Errorf("unmarshalled doctype has wrong name. Got %q, want %q", dtm.Name(), name)
		}
		d := dtm.New()
		if d == nil {
			t.Errorf("docType.New() called on unmarshalled docType returned nil")
		}
	}
}

func TestDocTypeJSONRoundTrip(t *testing.T) {
	t.Parallel()
	for _, dt := range docTypes {
		data, err := json.Marshal(dt)
		if err != nil {
			t.Errorf("cant marshal doctype %v: %s", dt, err)
			continue
		}
		var mdt DocType
		if err = json.Unmarshal(data, &mdt); err != nil {
			t.Errorf("can't unmarshal data %s into docType: %s", data, err)
			continue
		}
		if dt.Name() != mdt.Name() {
			t.Errorf("unmarshalled docType.Name()==%s, want %s", dt.Name(), mdt.Name())
		}
		newDoc := dt.New()
		newDocM := mdt.New()
		if reflect.TypeOf(newDoc) != reflect.TypeOf(newDocM) {
			t.Errorf("unmarshalled docType.New() has type %T, want %T", newDocM, newDoc)
		}
	}
}
