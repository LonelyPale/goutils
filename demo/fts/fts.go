package main

import (
	"log"

	"github.com/blevesearch/bleve"
)

func main() {
	// open a new index
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("example.bleve", mapping)
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Name string
	}{
		Name: "text 123, 你好，财务管理办法，关东（かんとう）（Kantō）东北（とうほく）（Tōhoku）",
	}

	// index some data
	err = index.Index("id-001", data)
	if err != nil {
		log.Fatal(err)
	}

	// search for some text
	query := bleve.NewMatchQuery("中国 东北")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(searchResults)
}
