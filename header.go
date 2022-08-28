package shuffler

import (
	"fmt"
	"time"
	"bitbucket.org/yasnikoff/shuffler/fs"
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"bitbucket.org/yasnikoff/shuffler/internal/accessors"
)


type Header struct {
	Id              id.ID       `json:"id"`
	path            string        `json:"-"`
	DocType         *DocType `json:"type"`
	CreatedTime     accessors.Time`json:"created"`
	LastSavedTime   accessors.Time `json:"saved"`
	LastCheckedTime accessors.Time // LastCheckedTime is time when the doc was checked if it is up to date the last time
	doc             Doc

}



func (h *Header)init(doc Doc, docType *DocType){
	h.CreatedTime.Set(time.Now())
	h.doc=doc
	h.DocType=docType
	h.Id =id.New()
}

func (h *Header)ID() id.ID {
	return h.Id
}
func (h *Header)SetID(ID id.ID) {
	h.Id = ID
}

func (h *Header)Path() string {
	return h.path
}
func (h *Header)SetPath(path string) error {
	p, err := fs.NormalizePath(path)
	if err != nil {
		return err
	}
	h.path = p
	return nil
}

func (h *Header)Type() *DocType {
	return h.DocType
}
func (h *Header)setType(t *DocType) {
	h.DocType = t
}

func (h *Header)String() string {
	if h == nil {
		return fmt.Sprintf("%v", h)
	}
	t := h.Type()
	var typeName string
	if t == nil {
		typeName = "unsaved doc"
	}else {
		typeName = t.Name()
	}
	return fmt.Sprintf("%s:{%s}", typeName, h.Id)
}

func (h *Header)GetHeader()*Header{
	return h
}

func (h *Header)SaveTo(path string)error{
	if h.doc==nil{
		return fmt.Errorf("Can't save header without doc. Header: %v", h)
	}
	return SaveDoc(h.doc, path)
}
func (h *Header)Save()error{
	return SaveDoc(h.doc, h.Path())
}
/*

func (h *Header)SavedTime() time.Time {
	return h.LastSavedTime
}

func (h *Header)SetSavedTime(t time.Time) {
	h.LastSavedTime = t
}
*/

func (h *Header)IsUpToDate(epsilon time.Duration) bool {
	// TODO: tests
	// TODO: use filesystem watcher to optimize performance

	// commented-out as needless in server mode:
	// only one instance (server) is accessing filesystem
	/*if  time.Since(h.LastCheckedTime.Get())<epsilon{
		return true
	}
	h2, err := LoadHeader(h.Path())
	if err != nil {
		return false
	}
	h.LastCheckedTime.Set(time.Now())
	return h2.LastSavedTime.Get() == h.LastSavedTime.Get()*/
	return true
}
