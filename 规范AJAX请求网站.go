package dev

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"megaCrawler/crawlers"
	"megaCrawler/extractors"
	"net/http"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Response struct {
	Results []ResultItem `json:"results"`
}

type ResultItem struct {
	Alias string `json:"alias"`
}

func FetchAndVisitArticles(engine *crawlers.WebsiteEngine, page int) {
	url := fmt.Sprintf("https://theedgemalaysia.com/api/loadMoreCategories?offset=%d&categories=news", page)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应失败: %v\n", err)
		return
	}

	var jsonResp Response
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		log.Printf("解析 JSON 失败: %v\n", err)
		return
	}

	for _, item := range jsonResp.Results {
		fullURL := "https://theedgemalaysia.com/" + item.Alias
		engine.Visit(fullURL, crawlers.News)
	}
}

// AJAX请求采集
func init() {
	engine := crawlers.Register("N-0017", "The Edge Malaysia", "https://theedgemalaysia.com/")

	engine.SetStartingURLs([]string{"https://theedgemalaysia.com/categories/news"})

	extractorConfig := extractors.Config{
		Author:       true,
		Image:        false,
		Language:     true,
		PublishDate:  true,
		Tags:         true,
		Text:         false,
		Title:        true,
		TextLanguage: "",
	}

	extractorConfig.Apply(engine)

	engine.OnLaunch(func() {
		for page := 14; page <= 40000; page += 10 {
			FetchAndVisitArticles(engine, page)
		}
	})

	engine.OnHTML("div.newsTextDataWrapInner > p", func(element *colly.HTMLElement, ctx *crawlers.Context) {
		if ctx.PageType == crawlers.News {
			if !strings.Contains(element.Text, "English version") {
				ctx.Content += element.Text
			}
		}
	})
}
