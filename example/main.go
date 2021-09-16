package main

import (
	"encoding/json"
	"github.com/yeyudekuangxiang/opensearch/api"
	"github.com/yeyudekuangxiang/opensearch/sdk"
	"log"
)

func main() {
	client := sdk.Client{
		AccessKeyId:     "AccessKeyId",
		AccessKeySecret: "AccessKeySecret",
		EndPoint:        "EndPoint",
	}

	resp := api.NewSearchResponse()
	rest := api.NewSearchRequest()

	rest.SetAppName("app_name")
	rest.SetFetchFields("id;title")
	rest.SetQueryStr("search:'吃饭'")
	rest.SetStart(0)
	rest.SetHits(10)
	rest.SetSort(api.SortStruct{
		Key:   "id",
		Order: api.SortDecrease,
	})
	rest.SetFilter("id>0 AND id<1000000")

	err := client.DoAction(rest.Build(), resp)
	if err != nil {
		log.Fatal("系统错误", err)
	}

	if !resp.IsSuccess() {
		log.Fatal("请求异常")
	}

	result := &api.SearchResult{}
	err = json.Unmarshal(resp.GetContentBytes(), result)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", result)
}
