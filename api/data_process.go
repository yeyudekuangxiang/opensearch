package api

import (
	"encoding/json"
	"fmt"
	"github.com/yeyudekuangxiang/opensearch/sdk"
)

type Cmd string

func NewDataProcessRequest() *DataProcessRequest {
	return &DataProcessRequest{
		params:      make([]map[string]interface{}, 0),
		OpenRequest: sdk.NewOpenRequest(),
	}
}

type DataProcessRequest struct {
	*sdk.OpenRequest
	app    string
	table  string
	params []map[string]interface{}
}

func (d *DataProcessRequest) SetAppName(app string) *DataProcessRequest {
	d.app = app
	return d
}
func (d *DataProcessRequest) SetTable(table string) *DataProcessRequest {
	d.table = table
	return d
}

//每次调用最多1000条 100条效果最好
//body大小不能超过2M
//map[string]interface{}{
//		"cmd":"delete",
//		"timestamp": 1401342874778,//标准版应用不支持 timestamp 参数。如果指定 timestamp 选项，推送会报4007错误码
//		"fields":map[string]interface{}{
//			"id":123123123,
//			"no":"ajin111",
//		},
//	}
func (d *DataProcessRequest) AddAction(action map[string]interface{}) {
	d.params = append(d.params, action)
}
func (d *DataProcessRequest) Delete(data map[string]interface{}) {
	d.AddAction(map[string]interface{}{
		"cmd":    "delete",
		"fields": data,
	})
}
func (d *DataProcessRequest) Update(data map[string]interface{}) {
	d.AddAction(map[string]interface{}{
		"cmd":    "update",
		"fields": data,
	})
}
func (d *DataProcessRequest) Add(data map[string]interface{}) {
	d.AddAction(map[string]interface{}{
		"cmd":    "add",
		"fields": data,
	})
}
func (d *DataProcessRequest) Build() *DataProcessRequest {
	data, _ := json.Marshal(d.params)
	d.Content = data
	d.Method = "POST"
	d.Path = fmt.Sprintf("/v3/openapi/apps/%s/%s/actions/bulk", d.app, d.table)
	return d
}

func NewDataProcessResponse() *DataResponse {
	return &DataResponse{}
}

type DataResponse struct {
	sdk.OpenResponse
}

//数据处理请求返回数据结构
// https://help.aliyun.com/document_detail/57154.html#h3-f0q-hqz-gh7
type ProcessResult struct {
	Errors []struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Params  interface{} `json:"params"`
	} `json:"errors"`
	Result    bool   `json:"result"`
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}
