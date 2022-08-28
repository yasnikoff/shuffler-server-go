package recentbuf

import (
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"encoding/json"
	"sync"
)

const MaxSize int = 16

type listType []id.ID


type Buffer struct {
	elements []id.ID
	lock sync.RWMutex
}

func New() *Buffer {
	return &Buffer{elements:make([]id.ID, 0, MaxSize)}
}

func (b *Buffer)Add(ID id.ID) {
	elCount := len(b.elements)

	// if id is in the elements, remove it

	b.lock.Lock()
	defer b.lock.Unlock()

	for i:=0;i<elCount;i++{
		if b.elements[i]==ID{
			copy(b.elements[i:], b.elements[i+1:])
			b.elements = b.elements[:elCount-1]
			elCount--
		}
	}

	if len(b.elements)<MaxSize{
		b.elements=append(b.elements, id.Zero)
	}
	copy(b.elements[1:], b.elements)
	b.elements[0] = ID


}

func (b *Buffer) Elements() []id.ID {
	return b.elements
}

func (b *Buffer)MarshalJSON() ([]byte, error) {
	return json.Marshal(b.elements)
}

func (b *Buffer)UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &b.elements)
}
