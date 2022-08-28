package api_test

import (
	//"os"
	"path/filepath"

	"bitbucket.org/yasnikoff/shuffler/repo"
	"bitbucket.org/yasnikoff/shuffler/types/id"
	"github.com/pkg/errors"
	"bitbucket.org/yasnikoff/shuffler"
	"image/color"
	_ "image/jpeg"
	"image"
	"image/jpeg"
	"os"
	"golang.org/x/image/draw"
	"fmt"
)

type DescData struct {
	ID byte
}

type DepData struct {
	SubPath string
	Data    interface{}
}

type ProjectData struct {
	Name         string
	ID           byte
	AddToRecent  bool
	PreviewColor color.RGBA
	Tags         []string
	Deps         []DepData
}

func IDFromNum(n byte) (id.ID, error) {
	if n != 0 {
		idStr := fmt.Sprintf("00000000-0000-0000-0000-000000000%03d", n)
		prID, err := id.Parse(idStr)
		if err != nil {
			return id.Zero, fmt.Errorf("can't parse ID from: %s", idStr)
		} else {
			return prID, nil
		}
	} else {
		return id.New(), nil
	}
}

func (pd *ProjectData)Create(parentDir string) (*shuffler.Project, error) {
	pr := shuffler.NewProject();

	// ID
	prID, err := IDFromNum(pd.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create project from data %s", pd)
	}
	pr.SetID(prID)

	// Path
	path := filepath.Join(parentDir, pr.ID().String())

	//Name
	pr.SetName(pd.Name)


	//tags
	pr.Tags.Add(pd.Tags...)

	//deps
	for _, d := range pd.Deps {
		var dep shuffler.Dependency
		data := d.Data
		switch data := data.(type){
		case DescData:
			depId, err := IDFromNum(data.ID)
			if err != nil {
				return nil, errors.Wrap(err, "can't create project's dependency")
			}
			dep = &shuffler.DepDesc{IDData:depId}
		default:
			return nil, fmt.Errorf("dep type is not supported: %#v", d)
		}
		pr.Deps.Add(d.SubPath, dep)
	}



	// SAVE

	if err := pr.SaveTo(path); err != nil {
		return nil, errors.Wrapf(err, "can't create project %q", pd.Name)
	}

	// create preview image

	rect := image.Rect(0, 0, 512, 512)
	img := image.NewRGBA(rect)

	draw.Draw(img, rect, image.NewUniform(pd.PreviewColor), image.ZP, draw.Src)
	file, err := os.Create(pr.PreviewPath())
	if err != nil {
		return nil, errors.Wrapf(err, "can't create preview file for %q at %q", pr.ID(), pr.PreviewPath())
	}
	if err := jpeg.Encode(file, img, nil); err != nil {
		return nil, errors.Wrapf(err, "can't encode preview file for %q at %q", pr.ID(), pr.PreviewPath())

	}

	return pr, nil
}

type RepoData struct {
	Name     string
	Projects []*ProjectData
}

func (rd *RepoData)Create(root string) error {
	repoPath := filepath.Join(root, rd.Name)
	//if err :=os.MkdirAll(repoPath, os.ModeDir);err!=nil{
	//	return err
	//}
	r := repo.New()
	r.SetName(rd.Name)
	if err := r.SaveTo(repoPath); err != nil {
		return errors.Wrapf(err, "can't create repo %q", rd.Name)
	}

	for _, prData := range rd.Projects {

		if pr, err := prData.Create(repoPath); err != nil {
			return errors.Wrapf(err, "can't create repo %q", rd.Name)
		} else {
			if prData.AddToRecent {
				r.AddToRecent(pr.ID())
			}
		}
	}
	return r.UpdateIndex();
}

func (rd *RepoData)Update(root string) error {
	if err := os.RemoveAll(root); err != nil {
		return errors.Wrapf(err, "can't update test data")
	}
	return rd.Create(root)
}
