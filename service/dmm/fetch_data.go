package dmm

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/faryne/api-server/models/entity"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func GetData(uri string) (*http.Response, error) {
	c := http.Client{}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{
		Name:  "age_check_done",
		Value: "1",
	})
	return c.Do(req)
}

func GetVideos(response *http.Response) (entity.DmmVideosList, error) {
	defer response.Body.Close()
	docs, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return entity.DmmVideosList{}, err
	}
	imgs := docs.Find("p.tmb > a")

	// 輸出內容定在這裡
	var videos []entity.DmmVideo
	var ch = make(chan entity.DmmVideoBody)

	imgs.Each(func(idx int, s *goquery.Selection) {
		// 解出網址 / 縮圖網址 / 標題
		link, _ := s.Attr("href")
		thumb, _ := s.Find("span > img").Attr("src")
		title, _ := s.Find("span > img").Attr("alt")

		// 解出番號
		NoPattern, _ := regexp.Compile("cid=([^/]+)")
		NoString := NoPattern.FindStringSubmatch(link)

		VideoHeader := entity.DmmVideoHeader{
			No:    strings.ToUpper(NoString[1]),
			Url:   link,
			Title: title,
			Thumb: thumb,
		}

		go func() {
			r, errParseVideo := ParseVideoPage(link)
			if errParseVideo == nil {
				ch <- r
			}
		}()

		videos = append(videos, entity.DmmVideo{
			VideoHeader,
			<-ch,
		})
	})

	output := entity.DmmVideosList{
		Videos: videos,
	}
	return output, nil
}

func ParseVideoPage(PageUrl string) (entity.DmmVideoBody, error) {
	// 定義 pattern
	Patterns := map[string]*regexp.Regexp{}
	Patterns["vod_date"], _ = regexp.Compile("配信開始")
	Patterns["publish_date"], _ = regexp.Compile("商品発売")
	Patterns["duration"], _ = regexp.Compile("収録時間")
	Patterns["directors"], _ = regexp.Compile("監督")
	Patterns["series"], _ = regexp.Compile("シリーズ")
	Patterns["makers"], _ = regexp.Compile("メーカー")
	Patterns["labels"], _ = regexp.Compile("レーベル")
	Patterns["tags"], _ = regexp.Compile("ジャンル")

	var VideoBody entity.DmmVideoBody
	// 解出詳細資訊
	out, err := GetData(PageUrl)
	if err != nil {
		return entity.DmmVideoBody{}, err
	}
	defer out.Body.Close()

	v, _ := goquery.NewDocumentFromReader(out.Body)

	VideoBody.Actresses = ParseActresses(*v, PageUrl)
	VideoBody.Images = parseImages(*v)

	rows := v.Find("table.mg-b20 > tbody > tr")
	rows.Each(func(idx int, e *goquery.Selection) {
		rowTitle, _ := e.Find("td.nw").Html()
		rowValue := e.Find("td:last-child")

		for k, v := range Patterns {
			if v.MatchString(rowTitle) {
				rowContent, _ := rowValue.Html()
				switch k {
				case "vod_date":
					VideoBody.VodDate = strings.TrimSpace(rowContent)
					break
				case "publish_date":
					VideoBody.PublishDate = strings.TrimSpace(rowContent)
					break
				case "duration":
					DurationPattern, _ := regexp.Compile("([0-9]+)")
					duration := DurationPattern.FindStringSubmatch(rowContent)
					VideoBody.Duration, _ = strconv.Atoi(strings.TrimSpace(duration[1]))
					break
				case "directors":
					VideoBody.Directors = make([]string, 0)
					rowValue.Find("a").Each(func(i int, d *goquery.Selection) {
						h, _ := d.Html()
						VideoBody.Directors = append(VideoBody.Directors, strings.TrimSpace(h))
					})
					break
				case "series":
					VideoBody.Series = make([]string, 0)
					rowValue.Find("a").Each(func(i int, d *goquery.Selection) {
						h, _ := d.Html()
						VideoBody.Series = append(VideoBody.Series, strings.TrimSpace(h))
					})
					break
				case "makers":
					VideoBody.Makers = make([]string, 0)
					rowValue.Find("a").Each(func(i int, d *goquery.Selection) {
						h, _ := d.Html()
						VideoBody.Makers = append(VideoBody.Makers, strings.TrimSpace(h))
					})
				case "labels":
					VideoBody.Labels = make([]string, 0)
					rowValue.Find("a").Each(func(i int, d *goquery.Selection) {
						h, _ := d.Html()
						VideoBody.Labels = append(VideoBody.Labels, strings.TrimSpace(h))
					})
				case "tags":
					VideoBody.Tags = make([]string, 0)
					rowValue.Find("a").Each(func(i int, d *goquery.Selection) {
						h, _ := d.Html()
						VideoBody.Tags = append(VideoBody.Tags, strings.TrimSpace(h))
					})
				default:
					fmt.Printf("----\n")
				}
			}
		}

	})

	return VideoBody, nil
}

func parseImages(document goquery.Document) []entity.DmmVideoImage {
	images := document.Find("img.mg-b6")

	var OututImages = []entity.DmmVideoImage{}

	images.Each(func(idx int, s *goquery.Selection) {
		thumb, _ := s.Attr("src")

		pattern, _ := regexp.Compile("(\\-[0-9]+\\.jpg)$")
		preview := pattern.ReplaceAllString(thumb, `jp${1}`)

		imageBody := entity.DmmVideoImage{
			Preview: preview,
			Thumb:   thumb,
		}

		OututImages = append(OututImages, imageBody)
	})

	return OututImages
}

func ParseActresses(document goquery.Document, url string) []string {
	element := document.Find("a#a_performer")

	var actresses = make([]string, 0)
	if element.Length() == 0 {
		document.Find("span#performer > a").Each(func(i int, s *goquery.Selection) {
			h, _ := s.Html()
			actresses = append(actresses, strings.TrimSpace(h))
		})
	} else {
		// 在網頁中找看看有沒有這個網址 pattern
		reg, _ := regexp.Compile(`'(/digital/videoa/-/detail/ajax-performer[^']+)'`)
		html, _ := document.Html()
		matches := reg.FindStringSubmatch(html)

		// 取得標籤列表網頁
		TagUrl := "https://www.dmm.co.jp" + matches[1]

		Client := http.Client{}
		req, _ := http.NewRequest("GET", TagUrl, strings.NewReader(""))
		req.Header.Add("Referer", url)
		resp, _ := Client.Do(req)
		defer resp.Body.Close()

		parser, _ := goquery.NewDocumentFromReader(resp.Body)

		parser.Find("a").Each(func(idx int, s *goquery.Selection) {
			h, _ := s.Html()
			actresses = append(actresses, strings.TrimSpace(h))
		})
	}
	return actresses
}

// GetActresses
// @TODO qoquery selectors must be reviewd
// url must be like: https://actress.dmm.co.jp/-/detail/=/actress_id=1000341/
func GetActresses(response *http.Response) (*entity.DMMActress, error) {
	defer response.Body.Close()
	docs, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}
	actress := &entity.DMMActress{}

	// 姓名
	name_field := docs.Find("h1.c-tx-actressName").Text()
	// 為了避免正規式難懂，把全形括號換成底線
	replacer := strings.NewReplacer(
		"（", "_",
		"）", "_")
	name_field = replacer.Replace(name_field)
	pattern_name, _ := regexp.Compile("\\s+(?P<Name>[^(\\n|\\s)]*)\\n(?P<Kana>.+)")
	if pattern_name.MatchString(name_field) {
		match := pattern_name.FindStringSubmatch(name_field)
		for i, names := range pattern_name.SubexpNames() {
			switch names {
			case "Name":
				actress.Name = strings.TrimSpace(match[i])
			case "Kana":
				actress.Kana = strings.TrimSpace(match[i])
			}
		}
	}
	// 大頭照
	actress.Photo, _ = docs.Find("span.p-section-profile__image > img").Attr("src")

	// 其他個人資料
	pattern_horoscope, _ := regexp.Compile(`星座`)
	pattern_blood, _ := regexp.Compile(`血液型`)
	pattern_born, _ := regexp.Compile(`出身地`)
	pattern_interests, _ := regexp.Compile(`趣味・特技`)
	pattern_body, _ := regexp.Compile(`サイズ`)
	pattern_day, _ := regexp.Compile(`生年月日`)
	pattern_3size, _ := regexp.Compile(`(T(?P<Height>[0-9]+)cm\s)?(B(?P<Bust>[0-9]+)cm)?(\((?P<Cup>[A-Z]{1,})カップ\))?(\sW(?P<Waist>[0-9]+)cm)?(\sH(?P<Hips>[0-9]+)cm)?`)
	pattern_birthday, _ := regexp.Compile(`((?P<Year>[0-9]{4})年)?((?P<Month>[0-9]{1,})月)?((?P<Day>[0-9]{1,})日)?`)

	values := make(map[int]string, 0)
	docs.Find(".p-list-profile__description").Each(func(i int, s *goquery.Selection) {
		v, _ := s.Html()
		values[i] = strings.TrimSpace(v)
	})
	docs.Find("dt.p-list-profile__heading").Each(func(i int, s *goquery.Selection) {
		header, _ := s.Html()
		value := values[i]
		header = strings.TrimSpace(header)

		if pattern_horoscope.MatchString(strings.TrimSpace(header)) {
			if value != "----" {
				actress.Horoscope = value
			}
		} else if pattern_blood.MatchString(strings.TrimSpace(header)) {
			if value != "----" {
				actress.Blood = value
			}
		} else if pattern_born.MatchString(header) {
			if value != "----" {
				actress.BornCity = value
			}
		} else if pattern_interests.MatchString(header) {
			actress.Interests = make([]string, 0)
			interests := strings.Split(value, "、")

			if len(interests) > 0 {
				for _, data := range interests {
					if data != "----" {
						actress.Interests = append(actress.Interests, data)
					}
				}
			}
		} else if pattern_body.MatchString(header) {
			match_3size := pattern_3size.FindStringSubmatch(value)
			for i, names := range pattern_3size.SubexpNames() {
				switch names {
				case "Height":
					actress.Height, _ = strconv.Atoi(match_3size[i])
				case "Bust":
					actress.Bust, _ = strconv.Atoi(match_3size[i])
				case "Cup":
					actress.Cup = match_3size[i]
				case "Waist":
					actress.Waist, _ = strconv.Atoi(match_3size[i])
				case "Hips":
					actress.Hips, _ = strconv.Atoi(match_3size[i])
				}
			}
		} else if pattern_day.MatchString(header) {
			match_bday := pattern_birthday.FindStringSubmatch(value)
			for i, names := range pattern_birthday.SubexpNames() {
				switch names {
				case "Year":
					actress.BirthYear, _ = strconv.Atoi(match_bday[i])
				case "Month":
					actress.BirthMonth, _ = strconv.Atoi(match_bday[i])
				case "Day":
					actress.BirthDay, _ = strconv.Atoi(match_bday[i])
				}
			}
			if actress.BirthYear > 0 && actress.BirthMonth > 0 && actress.BirthDay > 0 {
				actress.FullBirthday = fmt.Sprintf("%04d/%02d/%02d", actress.BirthYear, actress.BirthMonth, actress.BirthDay)
			}
		}
	})

	return actress, nil
}
