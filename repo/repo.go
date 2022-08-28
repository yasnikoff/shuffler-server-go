package repo

import (
	"path/filepath"
	"fmt"
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"bitbucket.org/yasnikoff/shuffler/fs"
	"bitbucket.org/yasnikoff/shuffler/logging"
	"bitbucket.org/yasnikoff/shuffler"
	"bitbucket.org/yasnikoff/shuffler/repo/recentbuf"
	"sync"
)


// Repo is responsible for providing Projects by ID.
//
// path that contains Project folders
type Repo struct {
	shuffler.Header
	NameData         string `json:"name"`
	Index            Map `json:"index"`
	TagsIndex        tagsIndex
	logging.Log
	LibData          string `json:"lib"`
	RecentData       recentbuf.Buffer `json:"recent,omitempty"`
	savedQueriesOnce sync.Once
	SavedQueriesData map[string]string `json:"saved queries, omitempty"`
}

func New() *Repo {
	r := DocType.New().(*Repo)
	r.SavedQueriesData = make(map[string]string)
	r.RecentData = *recentbuf.New()
	return r
}

func (r *Repo)Name() string {
	return r.NameData
}
func (r *Repo)SetName(n string) {
	r.NameData = n
}

func (r *Repo)GetByID(docID id.ID) (*shuffler.Project, error) {
	if r == nil {
		return nil, fmt.Errorf("can not get document from nil repo.")
	}
	subPath, ok := r.Index.Get(docID)
	if !ok {
		return nil, fmt.Errorf("ID %q not found", docID)
	}
	d, err := r.loadProject(subPath)
	if err != nil {
		return nil, err
	}
	d.SetRepo(r)
	return d, nil
}
func (r *Repo)GetManyByID(ids []id.ID) ([]*shuffler.Project, error) {
	var result []*shuffler.Project

	for _, ID := range ids {
		pr, err := r.GetByID(ID)
		if err != nil {
			return nil, err
		}
		result = append(result, pr)
	}

	return result, nil
}

func (r *Repo)fullPath(path string) string {
	return filepath.Join(r.libPath(), path)
}

func (r *Repo)AbsPath(path string) string {
	return filepath.Join(r.libPath(), path)
}

func (r *Repo)Root() string {
	return r.libPath()
}
func (r *Repo)Lib() string {
	return r.LibData
}

// Discover tries to find default repo for the path.
// It is the first repo met when going up in folder hierarchy starting from the path.
func Discover(path string) (*Repo, error) {
	var r *Repo
	curPath, err := fs.NormalizePath(path)
	var prevPath string
	for {
		prevPath = curPath
		r, err = LoadRepo(curPath)
		if err == nil {
			return r, nil
		}
		curPath = filepath.Dir(prevPath)
		if curPath == prevPath {
			break
		}
		if curPath == "/" || curPath == "." || curPath == "" {
			break
		}
	}
	return nil, fmt.Errorf("Can not find default repository for path %q", path)
}

func (r *Repo)loadProject(subPath string) (*shuffler.Project, error) {
	return shuffler.LoadProject(r.fullPath(subPath))
}

func (r *Repo)libPath() string {

	if r.LibData != "" {
		return filepath.Join(r.Path(), r.LibData)
	}
	return r.Path()
}

func (r *Repo)AddToIndex(pr *shuffler.Project) error {

	relPath, err := r.ToLibPath(pr.Path())
	if err != nil {
		return err
	}
	// TODO: tests for only slashes in the subPaths.
	// Note: testing may be done on Windows with backslashes as well as on *nix.
	relPath = filepath.ToSlash(relPath)
	oldPath, ok := r.Index.Get(pr.ID())
	if ok && oldPath != relPath {
		if _, err := shuffler.LoadDoc(r.AbsPath(oldPath)); err == nil {
			// TODO: tests
			return fmt.Errorf("can't add path %q for id %q: same ID is also found at %q", relPath, pr.ID(), oldPath)
		}
	}
	r.Index.Add(pr.ID(), relPath)
	r.TagsIndex.Add(pr.ID(), pr.Tags.Get()...)
	return nil
}

func (r *Repo)UpdateIndexForPath(path string) (err error) {
	// TODO: remove index entries for path that now has different ID
	relPath, err := r.ToLibPath(path)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			r.Logger().Printf("error updating index for %q:\n%s\n", relPath, err)
		}
		//r.Logger().Printf("Indexing ended.")
	}()

	//r.Logger().Printf("Start indexing path %s", relPath)
	err = shuffler.Walk(path, func(path string, d shuffler.Doc, err error) error {
		if err != nil {
			r.Logger().Printf("error: %s", err)
			return err
		}

		pr, ok := d.(*shuffler.Project)
		if !ok {
			return nil
		}
		if err := r.AddToIndex(pr); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	if err = r.Save(); err != nil {
		return err
	}
	return nil
}

func (r *Repo)UpdateIndex() error {
	return r.UpdateIndexForPath(r.libPath())
}

