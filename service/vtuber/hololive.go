package vtuber

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"net/url"
)

func HololiveSchedule() {
	req, _ := http.Get("https://schedule.hololive.tv/")
	defer req.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(req.Body)

	videos := make([]*youtube.VideoListResponse, 0)
	doc.Find("a[href*=youtube]").Each(func(i int, obj *goquery.Selection) {
		href, _ := obj.Attr("href")
		parser, _ := url.ParseRequestURI(href)
		ytId := parser.Query().Get("v")
		videoDetail, _ := getVideoDetail(ytId)
		videos = append(videos, videoDetail)
	})

	jsonContent, _ := json.Marshal(videos)
	fmt.Fprint(w, string(jsonContent))
}
