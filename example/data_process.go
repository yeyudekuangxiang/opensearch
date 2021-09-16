package main

import (
	"fmt"
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
	request := api.NewDataProcessRequest()
	request.SetAppName("app")
	request.SetTable("table_name")
	request.AddAction(map[string]interface{}{
		"cmd": "delete",
		"fields": map[string]interface{}{
			"id": 123123123,
			"no": "no123",
		},
	})
	resp := api.NewDataProcessResponse()
	err := client.DoAction(request.Build(), resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
