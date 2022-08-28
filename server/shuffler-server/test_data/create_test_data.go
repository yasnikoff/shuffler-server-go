package main

import (
	t "bitbucket.org/yasnikoff/shuffler/server/api/test"
	"log"
	"image/color"
	_ "image/jpeg"
)

func main() {
	root := "root"
	for _, rd := range testData {
		if err := rd.Update(root); err != nil {
			log.Fatal(err)
		}
	}
}

var testData []*t.RepoData = []*t.RepoData{
	{
		Name:"repo 1",
		Projects:[]*t.ProjectData{
			{
				Name:"project 1",
				ID: 1,
				AddToRecent:true,
				PreviewColor:color.RGBA{100, 100, 127, 0},
				Tags:[]string{"tag 1", "tag 2"},
				Deps:[]t.DepData{
					{
						SubPath:"links/dep1",
						Data:t.DescData{ID:2},
					},
					{
						SubPath:"links/dep2",
						Data:t.DescData{ID:4},
					},
				},
			},
			{
				Name:"project 2",
				ID: 2,
				AddToRecent:true,
				PreviewColor:color.RGBA{127, 127, 127, 0},
				Tags:[]string{"tag 21", "tag 22"},

			},
			{
				Name:"",
				AddToRecent:true,
				PreviewColor:color.RGBA{0, 0, 0, 0},

			},
			{
				Name:"project 4",
				ID:4,
				AddToRecent:false,
				PreviewColor:color.RGBA{127, 50, 50, 0},

			},
		},
	},
	{
		Name:"repo 2",
		Projects:[]*t.ProjectData{
			{
				Name:"project 21",
				AddToRecent:true,
				Tags:[]string{"tag1", "tag 2"},
			},
			{
				Name:"project 22",
				AddToRecent:true,
			},
			{
				Name:"project 23",
				AddToRecent:false,
			},
		},
	},
}

