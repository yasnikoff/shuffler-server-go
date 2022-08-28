package api

import (
	"sync"
	"testing"
)

func TestApps(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		//"empty":"",
		"simple":"app",
		"simple with space": "some app",
		"simple non-ascii": "приложение",
	}
	wg := sync.WaitGroup{}
	for caseName, appName := range tests {
		wg.Add(1)
		go func(caseName, appName string) {
			defer wg.Done()
			r := NewRunner(t, appName)
			defer r.Cleanup()

			r.CaseName = caseName
			r.Expected = struct {
				Name   string `json:"name"`
			}{
				Name:appName,
			}
			r.App.Repo().SetName(appName)

			r.RunGet()

		}(caseName, appName)
	}
	wg.Wait()

}



