package main

import (
	"flag"
	"bitbucket.org/yasnikoff/shuffler/server/api"
	"fmt"
	"log"
	"net/http"
)

const DefaultPort = 7474

func main() {
	p := flag.Int("p", DefaultPort, "port on localhost to serve api")

	flag.Parse()

	apiAddr := fmt.Sprintf(":%d", *p)

	args := flag.Args()
	if len(args) % 2 == 1 {
		log.Fatal("Uneven number of arguments is given.\n Provide NAME PATH pairs for each repo to serve")
	}
	var reposCount int
	reposCount = len(args) / 2
	repos := make(map[string]string, reposCount)
	for i := 0; i < reposCount; i++ {
		name, path := args[2 * i], args[2 * i + 1]
		if _, ok := repos[name]; ok {
			log.Fatal("name %s specified more than once", name)
		}
		repos[name] = path
	}
	go func(){
		// TODO: remove hardcoded path
		http.ListenAndServe(":8282", http.FileServer(http.Dir(`./build/web`)))
	}()
	log.Printf("Listening on port %d", *p)
	log.Fatal(api.ListenAndServe(apiAddr, repos))

}
