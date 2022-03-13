package actress

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/goccy/go-json"
	"golang.org/x/net/html/charset"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type DMMActress struct {
	Name         string   `json:"name",default:""`
	Kana         string   `json:"kana",default:""`
	Photo        string   `json:"photo"`
	Height       int      `json:"height"`
	Bust         int      `json:"bust"`
	Waist        int      `json:"waist"`
	Hips         int      `json:"hips"`
	Cup          string   `json:"cup"`
	Horoscope    string   `json:"horoscope"`
	Blood        string   `json:"blood"`
	BornCity     string   `json:"born_city"`
	BirthYear    int      `json:"birth_year"`
	BirthMonth   int      `json:"birth_month"`
	BirthDay     int      `json:"birth_day"`
	FullBirthday string   `json:"full_birthday"`
	Interests    []string `json:"interests"`
}

func Parse(reader io.Reader) (*DMMActress, error) {
	// https://github.com/djimenez/iconv-go 轉換出來的內容可能有問題，改用 charset 試試
	utfBody, e := charset.NewReader(reader, "text/html")
	if e != nil {
		return nil, e
	}
	if utfBody == nil {
		return nil, errors.New("reader cannot be initialized")
	}
	docs, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		return nil, err
		fmt.Printf(">>> Error: %s\n", err.Error())
	}

	// 初始化輸出物件
	actress := &DMMActress{}

	// 姓名
	nameField, e := docs.Find("td.t1 > h1").Html()
	if e != nil {
		return nil, e
	}
	// 為了避免正規式難懂，把全形括號換成底線
	replacer := strings.NewReplacer(
		"（", "_",
		"）", "_")
	nameField = replacer.Replace(nameField)
	patternName, _ := regexp.Compile("(?P<Name>[^（]*)_(?P<Kana>[^_]*)_")
	if patternName.MatchString(nameField) {
		match := patternName.FindStringSubmatch(nameField)
		for i, names := range patternName.SubexpNames() {
			switch names {
			case "Name":
				actress.Name = match[i]
			case "Kana":
				actress.Kana = match[i]
			}
		}
	}
	// 大頭照
	actress.Photo, _ = docs.Find("tr.area-av30.top > td:nth-child(1) > img").Attr("src")

	// 其他個人資料
	patternHoroscope, _ := regexp.Compile(`星座`)
	patternBlood, _ := regexp.Compile(`血液型`)
	patternBorn, _ := regexp.Compile(`出身地`)
	patternInterests, _ := regexp.Compile(`趣味・特技`)
	patternBody, _ := regexp.Compile(`サイズ`)
	patternDay, _ := regexp.Compile(`生年月日`)
	pattern3size, _ := regexp.Compile(`(T(?P<Height>[0-9]+)cm\s)?(B(?P<Bust>[0-9]+)cm)?(\((?P<Cup>[A-Z]{1,})カップ\))?(\sW(?P<Waist>[0-9]+)cm)?(\sH(?P<Hips>[0-9]+)cm)?`)
	patternBirthday, _ := regexp.Compile(`((?P<Year>[0-9]{4})年)?((?P<Month>[0-9]{1,})月)?((?P<Day>[0-9]{1,})日)?`)

	docs.Find("tr.area-av30.top > td:nth-child(2) > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		header, _ := s.Find("td:nth-child(1)").Html()
		value, _ := s.Find("td:nth-child(2)").Html()
		//fmt.Println(header)

		if patternHoroscope.MatchString(header) {
			if value != "----" {
				actress.Horoscope = value
			}
		} else if patternBlood.MatchString(header) {
			if value != "----" {
				actress.Blood = value
			}
		} else if patternBorn.MatchString(header) {
			if value != "----" {
				actress.BornCity = value
			}
		} else if patternInterests.MatchString(header) {
			actress.Interests = make([]string, 0)
			interests := strings.Split(value, "、")

			if len(interests) > 0 {
				for _, data := range interests {
					if data != "----" {
						actress.Interests = append(actress.Interests, data)
					}
				}
			}
		} else if patternBody.MatchString(header) {
			match3size := pattern3size.FindStringSubmatch(value)
			for i, names := range pattern3size.SubexpNames() {
				switch names {
				case "Height":
					actress.Height, _ = strconv.Atoi(match3size[i])
				case "Bust":
					actress.Bust, e = strconv.Atoi(match3size[i])
				case "Cup":
					actress.Cup = match3size[i]
				case "Waist":
					actress.Waist, _ = strconv.Atoi(match3size[i])
				case "Hips":
					actress.Hips, _ = strconv.Atoi(match3size[i])
				}
			}
		} else if patternDay.MatchString(header) {
			match_bday := patternBirthday.FindStringSubmatch(value)
			for i, names := range patternBirthday.SubexpNames() {
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
