package api

import (
	"errors"
	"sms/service/src/config"
	"sms/service/src/dao"
	"sms/service/src/dao/model"

	"github.com/huichen/wukong/engine"
	"github.com/huichen/wukong/types"
)

var (
	searcher = engine.Engine{}
)

func InitEngine() {
	searcher.Init(types.EngineInitOptions{
		SegmenterDictionaries: config.App.BasePath + "dictionary.txt",
	})
	ls := dao.AllOpenBlogs()
	for _, item := range ls {
		//utils.Log.Info("InitEngine:", item.Id, "|", item.Content)
		searcher.IndexDocument(uint64(item.Id), types.DocumentIndexData{Content: item.Content}, false)
	}
	searcher.FlushIndex()
}

func AddIndex(item *model.BlogCtx) {
	searcher.IndexDocument(uint64(item.Id), types.DocumentIndexData{Content: item.Content}, false)
	searcher.FlushIndex()
}

func Search(keywords string) ([]uint64, error) {
	res := []uint64{}
	resp := searcher.Search(types.SearchRequest{Text: keywords})
	//utils.Log.Info("Search:", keywords)
	//utils.Log.Info("Result:", resp)
	if resp.Timeout {
		return nil, errors.New("query timeout!")
	} else {
		if resp.NumDocs == 0 {
			return nil, errors.New("query no result!")
		} else {
			for _, data := range resp.Docs {
				res = append(res, data.DocId)
			}
			return res, nil
		}
	}
}

func CloseEngine() {
	searcher.Close()
}
