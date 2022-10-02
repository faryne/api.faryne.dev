package hanime_api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"net/http"
)

func NewUpload(ctx *fiber.Ctx) error {
	var uri = "https://hanime1.me/search?query=&genre=H%E5%8B%95%E6%BC%AB&sort=%E6%9C%80%E6%96%B0%E5%85%A7%E5%AE%B9"
	client := http.Client{}
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.83 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.AddCookie(&http.Cookie{Name: "__cf_bm", Value: "WRMjDMb3DRXWsg8cgfpeXt_N4VpYMvCxk5JIQV7M0VQ-1648397193-0-AQ7YXt3CnjhSTgWETGefUMjyA4KrTiW/mk3OEetldjyrHMbhQ+QqW5AZVicGxiUtf7UPXQyb55YKc1akUlLaMQDMG31oaTlJEyqeIM9SwnJ4UbMHkRLdT1EkBcx0UF05Lw=="})
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(content))
	//reader, _ := goquery.NewDocumentFromReader(resp.Body)
	//reader.Find(`a`).Each(func(i int, s *goquery.Selection) {
	//	fmt.Println(s.Html())
	//})
	return nil
}
