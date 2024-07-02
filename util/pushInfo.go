package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Notification struct {
	Body      string `json:"body"`
	Title     string `json:"title"`
	Badge     int    `json:"badge"`
	Level     string `json:"level"`
	Sound     string `json:"sound"`
	Icon      string `json:"icon"`
	Group     string `json:"group"`
	IsArchive string `json:"isArchive"`
	URL       string `json:"url"`
}
type BusinessData struct {
	MessageType string  `json:"msgtype"`
	Text        BusText `json:"text"`
}
type BusText struct {
	Content string `json:"content"`
}

var (
	logger = NewCustomLogger()
)

/*
推送bark,适用于ios
params: body:推送消息体 title:推送标题
*/
func PushMessageBark(body string, title string) {
	var err error
	notification := Notification{
		Body:      body,
		Title:     title,
		IsArchive: "1",
	}
	// 将结构体转换为 JSON 字节
	requestBody, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	reqURL := BarkURL + BarkKey
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		logger.Info("Bark信息推送成功")
		return
	} else {
		logger.Error("推送失败，请检查配置是否正确")
		return
	}
}

/*
企业微信推送
params : content:推送消息内容
*/
func PushWeChatBusiness(content string, args ...interface{}) {
	busdata := BusinessData{
		MessageType: "text",
		Text: BusText{
			Content: fmt.Sprintf(content, args...),
		},
	}
	client := &http.Client{}
	params := url.Values{}
	params.Add("key", BC.Push.BusWechat.BusinessKey)
	ReqURL := fmt.Sprintf("%s?%s", BC.Push.BusWechat.BusinessURL, params.Encode())
	busBody, err := json.Marshal(busdata)
	if err != nil {
		fmt.Println("json failed: ", err)
		return
	}
	req, err := http.NewRequest("POST", ReqURL, bytes.NewBuffer(busBody))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		fmt.Println("创建请求对象失败: ", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发起请求失败: ", err)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取消息体失败: ", err)
		return
	}
	defer resp.Body.Close()
	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return
	}
	if jsonResponse["errcode"].(float64) == 0 && jsonResponse["errmsg"].(string) == "ok" {
		logger.Info("企业微信信息推送成功")
		return
	} else {
		logger.Error("推送失败，请检查你的配置是否正确")
		return
	}
}
