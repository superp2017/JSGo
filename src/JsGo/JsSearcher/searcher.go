package main

import (
	"JsGo/JsConfig"
	. "JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"flag"
	"log"
	"sync"

	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
)

var (
	// searcher is coroutine safe
	searcher = riot.Engine{}
	mutex    sync.Mutex
	idBase   uint64 = 0
)

//最小ID Base
const (
	ID_THR = 10000000000
)

func InitSearcher() {
	var e error
	idBase, e = JsRedis.Redis_idbase()
	if e != nil {
		log.Fatal(e.Error())
	}

	if idBase < ID_THR {
		log.Fatal("search engine id base less then 10000000000, id database maybe wrong!")
	}

	dict, e := JsConfig.GetConfigString([]string{"Searcher", "Dict"})
	if e != nil {
		log.Fatal(e.Error())
	}

	stopToken, e := JsConfig.GetConfigString([]string{"Searcher", "StopToken"})
	if e != nil {
		log.Fatal(e.Error())
	}

	dir, e := JsConfig.GetConfigString([]string{"Searcher", "Dir"})
	if e != nil {
		log.Fatal(e.Error())
	}

	var (
		dictionaries = flag.String(
			"dictionaries",
			dict,
			"分词字典文件")
		stopTokenFile = flag.String(
			"stop_token_file",
			stopToken,
			"停用词文件")
		indexType               = flag.Int("indexType", types.DocIdsIndex, "索引类型")
		usePersistent           = flag.Bool("usePersistent", true, "是否使用持久存储")
		persistentStorageFolder = flag.String("persistentStorageFolder", dir, "持久存储数据库保存的目录")
		persistentStorageShards = flag.Int("persistentStorageShards", 16, "持久数据库存储裂分数目")

		options = types.RankOpts{
			OutputOffset: 0,
			MaxOutputs:   100,
		}

		// NumShards shards number
		NumShards = 2
	)

	searcher.Init(types.EngineOpts{
		SegmenterDict: *dictionaries,
		StopTokenFile: *stopTokenFile,
		IndexerOpts: &types.IndexerOpts{
			IndexType: *indexType,
		},
		NumShards:       NumShards,
		DefaultRankOpts: &options,
		UseStorage:      *usePersistent,
		StorageFolder:   *persistentStorageFolder,
		StorageShards:   *persistentStorageShards,
	})

	searcher.Flush()
}

func index(doc map[string]interface{}) (ret map[string]interface{}) {
	ret = make(map[string]interface{})
	mutex.Lock()
	defer mutex.Unlock()

	ti, ok := doc["Text"]

	sdoc := map[string]interface{}{}
	sdoc["ID"], ok = doc["ID"]
	if !ok {
		Error("ID key is error")
		ret["error"] = "ID key is error"
		return
	}
	sdoc["Type"], ok = doc["Type"]
	if !ok {
		Error("Type key is error")
		ret["error"] = "Type key is error"
		return
	}

	text := ""
	if ok {
		text, ok = ti.(string)
		if !ok {
			Error("Text key is error")
			ret["error"] = "Text key is error"
			return
		}
	} else {
		Error("Text key is error")
		ret["error"] = "Text key is error"
		return
	}

	idBase++
	e := JsRedis.Redis_idbase_update(idBase)
	if e != nil {
		Error(e.Error())
		ret["error"] = e.Error()
		return
	}
	searcher.IndexDoc(idBase, types.DocIndexData{Content: text})
	searcher.Flush()
	e = JsRedis.Redis_index(idBase, sdoc)
	if e != nil {
		Error(e.Error())
		ret["error"] = e.Error()
		return
	} else {
		ret["error"] = nil
		return
	}
}

func query(para map[string]string) (ret map[string]interface{}) {
	ret = make(map[string]interface{})
	for k, v := range para {
		if len(v) < 2 {
			continue
		}
		r := searcher.Search(types.SearchReq{Text: v})

		Docs, ok := r.Docs.(types.ScoredDocs)
		if !ok {
			Info("search no result!")
			ret["error"] = "search no result!"
			ret[k] = nil
			continue
		}
		ids := make([]uint64, len(Docs))
		for i, x := range Docs {
			ids[i] = x.DocId
		}
		var e error

		ret[k], e = JsRedis.Redis_query(ids)

		if e != nil {
			Error(e.Error())
			ret["error"] = e.Error()
			return
		}
	}

	ret["error"] = nil
	return
}
