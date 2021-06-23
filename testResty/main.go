package main

import (
	"github.com/go-resty/resty/v2"
	"log"
)

func main() {
	client := resty.New()
	type Body struct {
		Head struct {
			AppId   string `json:"app_id,omitempty"`
			DepotId int    `json:"depot_id,omitempty"`
		} `json:"head"`
		Data struct {
			TransforId     string `json:"transfor_id,omitempty"`
			TransforStatus string `json:"transfor_status,omitempty"`
			Owner          string `json:"owner,omitempty"`
			TesterInfo     struct {
				Tester string `json:"tester,omitempty"`
			} `json:"tester_info,omitempty"`
			CreateTime string `json:"create_time,omitempty"`
			EndTime    string `json:"end_time,omitempty"`
			PageInfo   struct {
				Start int `json:"start,omitempty"`
				Size  int `json:"size,omitempty"`
			} `json:"page_info"`
		} `json:"data"`
	}
	body := &Body{
		Head: struct {
			AppId   string `json:"app_id,omitempty"`
			DepotId int    `json:"depot_id,omitempty"`
		}{string(1234), 4567},
		Data: struct {
			TransforId     string `json:"transfor_id,omitempty"`
			TransforStatus string `json:"transfor_status,omitempty"`
			Owner          string `json:"owner,omitempty"`
			TesterInfo     struct {
				Tester string `json:"tester,omitempty"`
			} `json:"tester_info,omitempty"`
			CreateTime string `json:"create_time,omitempty"`
			EndTime    string `json:"end_time,omitempty"`
			PageInfo   struct {
				Start int `json:"start,omitempty"`
				Size  int `json:"size,omitempty"`
			} `json:"page_info"`
		}{
			Owner: "owner",
			TesterInfo: struct {
				Tester string `json:"tester,omitempty"`
			}{"handler"},
			TransforStatus: "statusName",
			PageInfo: struct {
				Start int `json:"start,omitempty"`
				Size  int `json:"size,omitempty"`
			}{Start: 1, Size: 1024},
		},
	}
	response, err := client.R().SetHeader("token", "81acef54586ee53ed1710b395dcaf1d0").SetBody(body).Post("http://test.zhiyan.oa.com/test_api/usecase_mgr/transfor_order/list")
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v\n", response)
}
