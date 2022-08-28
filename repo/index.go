package repo

import (
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"sync"
	"encoding/json"
)

type Map struct {
	once  sync.Once
	items map[id.ID]string
	rw sync.RWMutex
}

func (m *Map)init() {
	if m.items==nil{
		m.items = make(map[id.ID]string)
	}
}

func (m *Map)Add(Id id.ID, path string) {
	m.once.Do(m.init)
	m.rw.Lock()
	defer m.rw.Unlock()
	m.items[Id] = path
}

func (m *Map)Get(Id id.ID) (path string, ok bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	path, ok = m.items[Id]
	return
}

func (m *Map)Items()map[id.ID]string{
	m.once.Do(m.init)
	return m.items
}

func (m *Map)MarshalJSON() ([]byte, error) {
	data := make(map[string]string)
	for Id, path := range m.items {
		data[Id.String()] = path
	}
	return json.Marshal(data)
}

func (m *Map)UnmarshalJSON(data []byte) error {
	m.once.Do(m.init)
	dataMap := make(map[string]string)
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return err
	}
	for IDStr, path := range dataMap {
		ID, err := id.Parse(IDStr)
		if err != nil {
			return err
		}
		m.items[ID] = path
	}
	return nil
}

func (m *Map)Len() int {
	return len(m.items)
}




