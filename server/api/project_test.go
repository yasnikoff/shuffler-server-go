package api

import (
	"testing"
	"sync"
)

type prData struct {
	ID     string        `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Origin string`json:"origin"`
}

func TestProjects(t *testing.T) {
	t.Parallel()

	tests := map[string]prData{
		"empty":{},
		"with name":{Name:"project1"},
	}
	wg := sync.WaitGroup{}
	for caseName, testPrData := range tests {
		wg.Add(1)
		go func(caseName string, testPrData prData) {
			defer wg.Done()
			r := NewRunner(t, caseName)
			defer r.Cleanup()

			r.CaseName = caseName
			pr := r.NewProject()
			pr.SetName(testPrData.Name)
			pr.Save()
			testPrData.Origin = pr.Origin.String()
			testPrData.ID = pr.ID().String()
			libPath, err:= r.App.Repo().ToLibPath(pr.Path())
			if err!=nil{
				t.Errorf("Case %q: %s", r.CaseName, err)
			}
			testPrData.Path = libPath
			r.Expected = testPrData

			r.RunGet("projects", pr.ID().String())

		}(caseName, testPrData)
	}
	wg.Wait()
}

