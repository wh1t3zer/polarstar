package news

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"polarstar/util"
)

var (
	logger = util.NewCustomLogger()
)

type JinSeLive struct {
	Id        int    `json:"id"`
	Content   string `json:"content"`
	Link      string `json:"link"`
	Grade     int    `json:"grade"`
	CreatedAt int64  `json:"created_at"`
}

type JinSeDate struct {
	Lives []JinSeLive `json:"lives"`
}

type JinSeResponse struct {
	List []JinSeDate `json:"list"`
}

var lastFetchTime int64 // 记录上一次爬虫的时间戳

// SpiderJinSe 金色财经爬虫
func SpiderJinSe() ([]JinSeLive, error) {
	client := &http.Client{}
	reqURL := "https://api.jinse.cn/noah/v2/lives?reading=false&sort=&flag=up&id=0&limit=20&_source=m"
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		fmt.Println("创建请求失败: ", err)
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Mobile Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取消息体失败: ", err)
		return nil, err
	}

	var result JinSeResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("JSON解码失败: ", err)
		return nil, err
	}
	if len(result.List) == 0 || len(result.List[0].Lives) == 0 {
		logger.ERROR("未获取到有效的新闻数据")
		return nil, nil
	}

	var newLives []JinSeLive
	for _, date := range result.List {
		for _, live := range date.Lives {
			if live.CreatedAt > lastFetchTime {
				newLives = append(newLives, live)
			}
		}
	}

	// 更新 lastFetchTime 为最新的新闻时间戳
	if len(newLives) > 0 {
		lastFetchTime = newLives[0].CreatedAt
	}

	return newLives, nil
}

func HandlerSpider() {
	logger.Info("获取金色财经信息模块初始化")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	lastFetchTime = time.Now().Unix()
	for range ticker.C {
		newLives, err := SpiderJinSe()
		if err != nil {
			fmt.Println("获取新闻数据失败: ", err)
			continue
		}
		if len(newLives) > 0 {
			pushInfo := make([]string, 0, 7)
			for _, live := range newLives {
				timeStr := time.Unix(live.CreatedAt, 0).Format("2006-01-02 15:04:05")
				info := fmt.Sprintf("ID: %d\n内容: %s\n链接: %s\n评分: %d\n文章时间: %s", live.Id, live.Content, live.Link, live.Grade, timeStr)
				pushInfo = append(pushInfo, info)
				if len(pushInfo) > 7 {
					pushInfo = pushInfo[len(pushInfo)-7:]
				}
			}
			resultStr := strings.Join(pushInfo, "\n\n")
			util.PushWeChatBusiness(resultStr)
		}
	}
}
