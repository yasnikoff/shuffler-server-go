package shuffler

import (
	"io/ioutil"
	"bitbucket.org/yasnikoff/shuffler/internal/md"
	"encoding/json"
	"fmt"
	"time"
	"os"
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"bitbucket.org/yasnikoff/shuffler/fs"
)

// SaveToDisk saves given document to the metadata file.txt.
func SaveDoc(d Doc, newPath string) (err error) {

	jsonBkp, err := json.Marshal(d)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if e := json.Unmarshal(jsonBkp, d); e != nil {
				panic(fmt.Sprintf("setting own jsondata to doc returned error: %s", e))
			}
			err = fmt.Errorf("can't save %s to %q: %s", d, newPath, err)
		}
	}()

	h := d.GetHeader()

	if h.Path() == "" {
		h.SetPath(newPath)
	}

	if h.ID() == id.Zero {
		h.SetID(id.New())
	}

	path := h.Path()
	info := md.GetInfo(path)

	if info.File.Exists {
		hp, err := LoadHeader(path)

		if err != nil {
			return err
		}
		if hp.ID() != h.ID() {
			return fmt.Errorf("there is another doc there: %v", hp)
		}
	}

	if !info.Dir.Exists {
		if err := os.MkdirAll(md.DirPath(path), os.ModeDir); err != nil {
			return err
		}
	}

	h.LastSavedTime.Set(time.Now())
	docBytes, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(md.FilePath(path), docBytes, 777); err != nil {
		return err
	}
	return nil
}

// Load populates Document with metadata from the path.
func LoadDoc(path string) (doc Doc, err error) {
	// TODO: change Load so it doesn't have to receive doc as a parameter?

	path, err = fs.NormalizePath(path)
	if err != nil {
		return nil, err
	}
	data, err := loadBytes(path)
	if err != nil {
		return nil, err
	}
	h, err := readHeader(data)
	if err != nil {
		return nil, err
	}
	dt := h.Type()
	d := dt.New()

	if err := json.Unmarshal(data, d); err != nil {
		return nil, err
	}

	d.SetPath(path)
	return d, nil
}

func loadBytes(path string) ([]byte, error) {
	path, err := fs.NormalizePath(path)
	if err != nil {
		return nil, err
	}
	mdFile := md.FilePath(path)

	data, err := ioutil.ReadFile(mdFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, md.ErrPathNotInitialized
		}
		return nil, err
	}

	return data, nil

}

// LoadHeader loads document's header from the path
// It can be used with any kind of document saved in the path.
// The header holds, among other things, info about document's type.
func LoadHeader(path string) (*Header, error) {
	// TODO: how not to read the entire document?
	headData, err := loadBytes(path)
	if err != nil {
		return nil, err
	}

	h, err := readHeader(headData)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func readHeader(data []byte) (*Header, error) {
	head := &Header{}
	if err := json.Unmarshal(data, head); err != nil {
		return nil, err
	}
	return head, nil
}
