package dmm

import (
	"fmt"
	"github.com/faryne/api-server/models/apireq"
	"io"
	"io/ioutil"
	"net/http"
)

func GetData(r *apireq.DMMCrawlerRequest) (interface{}, error) {
	c := http.Client{}
	var resp io.Reader
	req, _ := http.NewRequest("GET", r.Url, resp)
	req.AddCookie(&http.Cookie{
		Name:  "age_check_done",
		Value: "1",
	})
	out, err2 := c.Do(req)
	// 如果發生錯誤時
	if err2 != nil {
		fmt.Printf(">>> Error: %s\n", err2)
		return nil, err2
	}
	//
	// 不是 200 時
	if out.StatusCode != 200 {
		fmt.Printf(">>> Error (HTTP): %d\n", out.StatusCode)
		responseContent, _ := ioutil.ReadAll(out.Body)
		fmt.Printf(string(responseContent))
		return nil, nil
	}
	return nil, nil
}
