package api

import (
	"fmt"
	"github.com/yeyudekuangxiang/opensearch/sdk"
	"strings"
)

type SortType int

const (
	SortDecrease SortType = iota //降序
	SortIncrease                 //升序

)

type SortStruct struct {
	Key   string
	Order SortType
}
type SearchRequest struct {
	*sdk.OpenRequest
	Start          int
	Hits           int
	Format         string
	AppName        string
	QueryStr       string
	FetchFields    string
	Filter         string
	Sort           []SortStruct
	FirstRankName  string
	SecondRankName string
	QP             string
}

func NewSearchRequest() *SearchRequest {
	return &SearchRequest{
		OpenRequest: sdk.NewOpenRequest(),
		Sort:        make([]SortStruct, 0),
		Start:       0,
		Hits:        10,
		Format:      "fulljson",
	}
}

//SetStart 设置返回结果的偏移量 范围[0,5000] 默认0
func (sq *SearchRequest) SetStart(start int) {
	sq.Start = start
}

//设置返回结果的条数。范围[0,500] 默认10
func (sq *SearchRequest) SetHits(hits int) {
	sq.Hits = hits
}

//指定要搜索的应用名称或ID
func (sq *SearchRequest) SetAppName(appName string) {
	sq.AppName = appName
}

//设置的搜索关键词，格式为：索引名:'关键词' [AND|OR ...]
func (sq *SearchRequest) SetQueryStr(queryStr string) {
	sq.QueryStr = queryStr
}

//指定的返回字段的列表，例如 id;title
func (sq *SearchRequest) SetFetchFields(fetchFields string) {
	sq.FetchFields = fetchFields
}

//过滤，例如id>1 AND id<100 OR id=1000
func (sq *SearchRequest) SetFilter(filter string) {
	sq.Filter = filter
}

//返回结果的格式，有json、fulljson和xml格式 默认fulljson
func (sq *SearchRequest) SetFormat(format string) {
	sq.Format = format
}

//排序策略
func (sq *SearchRequest) SetSort(sort ...SortStruct) {
	sq.Sort = append(sq.Sort, sort...)
}

//指定的粗排表达式名称
func (sq *SearchRequest) SetFirstRankName(firstRankName string) {
	sq.FirstRankName = firstRankName
}

//指定的精排表达式名称
func (sq *SearchRequest) SetSecondRankName(secondRankName string) {
	sq.SecondRankName = secondRankName
}
func (sq *SearchRequest) SetQP(qp string) {
	sq.QP = qp
}

//格式化搜索参数
func (sq *SearchRequest) Build() *SearchRequest {
	sq.Path = fmt.Sprintf("/v3/openapi/apps/%s/search", sq.AppName)
	sq.Query["query"] = ""

	if len(sq.FetchFields) > 0 {
		sq.Query["fetch_fields"] = sq.FetchFields
	}
	if len(sq.FirstRankName) > 0 {
		sq.Query["first_rank_name"] = sq.FirstRankName
	}
	if len(sq.SecondRankName) > 0 {
		sq.Query["second_rank_name"] = sq.SecondRankName
	}
	if len(sq.QP) > 0 {
		sq.Query["qp"] = sq.QP
	}

	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("query=%s&&config=start:%d,hit:%d,format:%s", sq.QueryStr, sq.Start, sq.Hits, sq.Format))
	if len(sq.Filter) > 0 {
		str.WriteString(fmt.Sprintf("&&filter=%s", sq.Filter))
	}
	sortStr := ""
	for _, sort := range sq.Sort {
		if sort.Order == SortIncrease {
			sortStr += "+" + sort.Key + ";"
		} else {
			sortStr += "-" + sort.Key + ";"
		}
	}
	if len(sortStr) > 0 {
		str.WriteString(fmt.Sprintf("&&sort=%s", sortStr))
	}
	sq.Query["query"] = str.String()
	sq.Headers["Content-MD5"] = ""
	sq.Headers["Content-Type"] = "application/json"
	return sq
}

type SearchResponse struct {
	*sdk.OpenResponse
}

func NewSearchResponse() *SearchResponse {
	return &SearchResponse{
		&sdk.OpenResponse{},
	}
}

//搜索处理返回数据结构
//https://help.aliyun.com/document_detail/57155.html#h2-u8FD4u56DEu7ED3u679C7
type SearchResult struct {
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
	Result    struct {
		Searchtime  float64 `json:"searchtime"`
		Total       int     `json:"total"`
		Num         int     `json:"num"`
		Viewtotal   int     `json:"viewtotal"`
		ComputeCost []struct {
			IndexName string  `json:"index_name"`
			Value     float64 `json:"value"`
		} `json:"compute_cost"`
		Items []struct {
			Fields   interface{} `json:"fields"`
			Property struct {
			} `json:"property"`
			Attribute struct {
			} `json:"attribute"`
			VariableValue struct {
			} `json:"variableValue"`
			SortExprValues []string `json:"sortExprValues"`
		} `json:"items"`
		Facet []interface{} `json:"facet"`
	} `json:"result"`
	Errors         []interface{} `json:"errors"`
	Tracer         string        `json:"tracer"`
	OpsRequestMisc string        `json:"ops_request_misc"`
}
