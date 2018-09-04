package main

import (
	"log"

	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
)

var (
	// searcher is coroutine safe
	searcher = riot.Engine{}
)

func main() {
	// Init
	searcher.Init(types.EngineOpts{
		Using:       4,
		NotUsingGse: true})
	defer searcher.Close()

	text := "Google Is Experimenting With Virtual Reality Advertising"
	text1 := `Google accidentally pushed Bluetooth update for Home
	speaker early`
	text2 := `Google is testing another Search results layout with 
	rounded cards, new colors, and the 4 mysterious colored dots again`

	// Add the document to the index, docId starts at 1
	searcher.IndexDoc(1, types.DocIndexData{Content: text})
	searcher.IndexDoc(2, types.DocIndexData{Content: text1}, false)
	searcher.IndexDoc(3, types.DocIndexData{Content: text2}, true)

	// Wait for the index to refresh
	searcher.Flush()
	// engine.FlushIndex()

	// The search output format is found in the types.SearchResp structure
	ret := searcher.Search(types.SearchReq{Text: "pushed"})
	log.Print(ret.Docs)
}
