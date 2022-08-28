package shuffler

import "path/filepath"

func (pr *Project)PreviewPath()string{
	return filepath.Join(pr.Path(), "preview.jpg")
}
