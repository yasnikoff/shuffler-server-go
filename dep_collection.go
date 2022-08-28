package shuffler

import (
	"encoding/json"
	"sort"
	"fmt"
	"strings"
	"path/filepath"
	"os"
	"github.com/pkg/errors"
	"sync"
)

var depTypes map[string]func() Dependency

func init() {
	depTypes = map[string]func() Dependency{
		DepTypeDepDesc: func() Dependency {
			return new(DepDesc)
		},
	}
}

type Dependency interface {
	Descriptor() *DepDesc
	Type() string
}

type Dependencies struct {
	project   *Project
	items     map[string]Dependency
	itemsLock sync.RWMutex
}

func newDeps(pr *Project) *Dependencies {
	return &Dependencies{project:pr, items:make(map[string]Dependency)}

}

func (dc *Dependencies)String() string {
	items := make([]string, len(dc.items))
	i := 0
	for path, dep := range dc.items {
		items[i] = fmt.Sprintf("%q:%s", path, dep)
		i++
	}
	return fmt.Sprintf("dependencies of %s:\n%s", dc.project, strings.Join(items, "\n\t"))
}

// Add adds dependency to the collection.
func (dc *Dependencies)Add(subPath string, dep Dependency) error {
	// TODO: dependency cycles check?

	dc.itemsLock.Lock()
	defer dc.itemsLock.Unlock()

	dc.items[subPath] = dep
	return nil
}
// Remove removes dependency at given subPath.
// It does not delete files from file system. User should call Cleanup on dependency for that.
func (dc *Dependencies)Remove(subPath string) {
	dc.itemsLock.Lock()
	defer dc.itemsLock.Unlock()

	delete(dc.items, subPath)
}

// CleanUp deletes all files got from dependent Projects.
func (dc *Dependencies)CleanUpAll() error {

	for subPath := range dc.items {
		if err := dc.CleanUp(subPath); err != nil {
			return err
		}
	}
	return nil
}

func (dc *Dependencies)CleanUp(subPath string) error {
	return os.RemoveAll(dc.Abs(subPath))
}

func (dc *Dependencies)normPath(subPath string) string {
	return filepath.ToSlash(subPath)
}

func (dc *Dependencies)Get(subPath string) (Dependency, error) {
	dc.itemsLock.RLock()
	defer dc.itemsLock.RUnlock()

	d, ok := dc.items[subPath]
	if !ok {
		return d, fmt.Errorf("there is no dependency %q in %s", subPath, dc)
	}
	return d, nil
}

// Items returns copy of Collection's internal map of descriptors
func (dc *Dependencies)Items() map[string]Dependency {
	return dc.items
}

func (dc *Dependencies)SetItems(m map[string]Dependency) error {
	dc.items = m
	// TODO: cyclic deps check
	return nil
}

func (dc *Dependencies)Count() int {
	return len(dc.items)
}

func (dc *Dependencies)Project() *Project {
	return dc.project
}

func (dc *Dependencies)SortedPaths() []string {
	result := make([]string, len(dc.items))
	i := 0
	for path := range dc.items {
		result[i] = path
		i++
	}
	sort.Strings(result)
	return result
}

func (dc *Dependencies)UpdateAll() error {
	for _, subPath := range dc.SortedPaths() {
		if err := dc.Update(subPath); err != nil {
			return err
		}
	}
	return nil
}

func (dc *Dependencies)Update(subPath string) error {
	if err := dc.CleanUp(subPath); err != nil {
		return err
	}
	return dc.Setup(subPath)
}

func (dc *Dependencies)Setup(subPath string) error {
	dep, err := dc.Get(subPath)
	if err != nil {
		return err
	}
	desc := dep.Descriptor()
	depPr, err := dc.project.repo.GetByID(desc.ID())
	if err != nil {
		return err
	}
	return depPr.SetupPath(dc.Abs(subPath), desc.Parameters)
}

func (dc *Dependencies)Abs(path string) string {
	return filepath.Join(dc.project.Path(), path)
}

type depProviderData struct {
	Type string `json:"type"`
	Data Dependency `json:"data"`
}

func (dc *Dependencies)MarshalJSON() ([]byte, error) {
	data := make(map[string]depProviderData, len(dc.items))

	for name, item := range dc.items {
		data[name] = depProviderData{Type:item.Type(), Data:item}
	}

	return json.Marshal(data)
}

type depProviderRawData struct {
	Type string `json:"type"`
	Data json.RawMessage
}

func (dc *Dependencies)UnmarshalJSON(data []byte) error {

	d := make(map[string]depProviderRawData)
	if err := json.Unmarshal(data, &d); err != nil {
		return errors.Wrap(err, "can't unmarshal dependencies data")
	}
	dc.items = make(map[string]Dependency, len(d))
	for path, data := range d {
		ddp, err := NewDep(data.Type)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data.Data, ddp); err != nil {
			return errors.Wrapf(err, "can't unmarshal data %s to a variable of type %q", data.Data, data.Type)
		}

		dc.items[path] = ddp
	}
	return nil
}

func NewDep(depType string) (Dependency, error) {
	constructor, ok := depTypes[depType]
	if !ok {
		return nil, fmt.Errorf("no such dependency type: %q", depType)
	}
	return constructor(), nil
}
