package EDU

import (
	"errors"
	"learn/models"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	Val *http.Client
}

type WsData struct {
	Done bool
	Data string
}

const (
	Host     string = "jwxt.nfsysu.cn"
	Host_url string = "http://" + Host
)

var (
	client *Client           = &Client{Val: new(http.Client)}
	header map[string]string = map[string]string{
		"Host":            Host,
		"User-Agent":      "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:37.0) Gecko/20100101 Firefox/37.0",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Charset":  "GBK,utf-8;q=0.7,*;q=0.3",
		"Accept-Language": "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
		"Accept-Encoding": "gzip, deflate",
		"Connection":      "keep-alive",
		"Content-Type":    "application/x-www-form-urlencoded; param=value",
		"Referer":         Host_url,
	}
)

func Sign_in(id, pwd, user_type string) ([]byte, bool, []*http.Cookie, error) {
	if user_type != "学生" && user_type != "教师" {
		return nil, false, nil, models.ErrorInfo("Sign_in", errors.New("user type is error："+user_type))
	}
	if len(id) <= 0 || len(pwd) <= 0 {
		return nil, false, nil, models.ErrorInfo("Sign_in", errors.New("id or pwd is error：["+id+"_"+pwd+"]"))
	}
	data, cookies, err := client.SendRequest("GET", Host_url, "", nil, nil)
	if err != nil {
		return nil, false, nil, models.ErrorInfo("Sign_in", err)
	}
	vierstate, eventvalidation, err := getLoginKey(string(data))
	if err != nil {
		return nil, false, nil, models.ErrorInfo("Sign_in", err)
	}
	//	sigin...
	post_arg := url.Values{
		"__VIEWSTATE":       {vierstate},
		"__EVENTVALIDATION": {eventvalidation},
		"TextBox1":          {id},
		"TextBox2":          {pwd},
		"RadioButtonList1":  {models.Utf82gbk(user_type)},
		"Button1":           {""},
	}
	//	登陆
	main_data, _, err := client.SendRequest("POST", Host_url+"/default2.aspx", post_arg.Encode(), cookies, header)
	if err != nil {
		return nil, false, nil, models.ErrorInfo("Sign_in", err)
	}
	//	判断是否登录失败
	if strings.Contains(string(main_data), "<title>登录</title>") {
		return nil, false, nil, models.ErrorInfo("Sign_in", errors.New("id or pwd not match"))
	}
	return main_data, true, cookies, nil
}
