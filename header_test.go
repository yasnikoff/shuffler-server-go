package shuffler

import (
	"testing"
	"encoding/json"
	"bitbucket.org/yasnikoff/shuffler/types/id"
)

const testDocTypeName = "testdoc"

type testdoc struct {
	Header
}

var testdoctype *DocType

func init() {
	testdoctype = RegisterType(testDocTypeName, func() Doc {
		return &testdoc{}
	})
}

func checkHeader(msg string, doc Doc, h *Header, t *testing.T) {
	//t.Parallel()
	h2 := doc.GetHeader()
	if h2.ID() != h.ID() {
		t.Errorf("%s: Unmarshalled header's ID doesn't match original one.\nGot:\t%s,\nwant:\t%s", msg, h.ID(), h2.ID())
	}
	if h2.Type().Name() != h.Type().Name() {
		t.Errorf("%s: Unmarshalled header's Type's name doesn't match original one.\nGot:\t%s,\nwant:\t%s", msg, h.Type().Name(), h2.Type().Name())
	}
	if h.doc != doc {
		t.Errorf("%s: Header.doc is %v (type: %[2]T), want %v (type: %[3]T)", msg, h.doc, doc)
	}
}

func TestJSON_zeroValue(t *testing.T) {
	//t.Parallel()
	for _, dt := range docTypes{
		d := dt.New()
		data, err := json.Marshal(d)
		if err != nil {
			t.Fatalf("cant marshal header %v: %s", d, err)
		}
		h := &Header{}
		h.doc=d // Header.doc is not marshaled by json
		if err = json.Unmarshal(data, h); err != nil {
			t.Fatalf("can't unmarshal data %s\ninto testdoc: %s", data, err)
		}
		checkHeader("Json Zero value test", d, h, t)
	}

}

func TestJSON_NonZeroValue(t *testing.T) {
	//t.Parallel()

	for _, dt := range docTypes {
		d := dt.New()
		h := d.GetHeader()
		h.SetID(id.New())
		data, err := json.Marshal(d)
		if err != nil {
			t.Fatalf("cant marshal header %v: %s", d, err)
		}
		h1 := &Header{}
		if err = json.Unmarshal(data, h1); err != nil {
			t.Fatalf("can't unmarshal data %s\ninto Header: %s", data, err)
		}
		checkHeader("json non-zero value test", d, h, t)
	}

}
/*

func checkHeaderData(di DocData, h *Header, t *testing.T) {
	//t.Parallel()
	d := di.(*HeaderData)
	if d.ID != h.Id {
		t.Errorf("Header's data ID is different form header's ID. Data.ID()==%s, Header.ID()==%s", d.ID, h.Id)
	}
	//if data.docType!=h.docType{
	//	t.Errorf("Header's data docType is different form header's docType. Data.docType==%s, Header.docType==%s", data.docType, h.docType)
	//}
	if d.Type != h.DocType {
		t.Errorf("Header's data docType is different form header's docType. Data.docType==%s, Header.docType==%s", d.Type, h.DocType)
	}
	if d.Path != h.path {
		t.Errorf("Header's data Path is different form header's path. Data.Path==%s, Header.path==%s", d.Path, h.path)
	}
	if d.Saved != h.LastSavedTime {
		t.Errorf("Header's data Saved time is different form header's lastSavedTime. Data.Saved==%s, Header.lastSavedTime==%s", d.Saved, h.LastSavedTime)
	}

}

func TestData(t *testing.T) {
	//t.Parallel()
	ID := id.New()
	dt := testdoctype
	h := &Header{Id:ID, DocType:dt, path:"some/path", LastSavedTime:time.Now()}
	d := h.Data()
	checkHeaderData(d, h, t)
}

func TestSetState()t *testing.T) {
	//t.Parallel()
	h := &Header{}
	d := &HeaderData{ID:id.New(), Type:testdoctype, Path:"some/path", Saved:time.Now()}
	h.SetState()d)
	checkHeaderData(d, h, t)
}
*/
