package EDU

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func (client *Client) SendRequest(method, send_url, data string, cookie []*http.Cookie, header map[string]string) (bytes []byte, cookies []*http.Cookie, err error) {
	// 调用string包读取发送消息
	arg := strings.NewReader(data)

	//创建http请求，写入方法，地址和发送数据
	req, err := http.NewRequest(method, send_url, arg)
	if err != nil {
		return nil, nil, err
	}

	//	设置回复头
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	//	设置cookie
	if len(cookie) > 0 {
		for _, v := range cookie {
			req.AddCookie(v)
		}
	}

	client.Val.CheckRedirect = func(rq *http.Request, via []*http.Request) error {
		if len(via) == 1 {
			return errors.New("SendRequest: 302")
		}
		if len(via) >= 10 {
			return errors.New("SendRequest: stopped after 10 redirects")
		}
		return nil
	}

	//	发送http请求
	res, err := client.Val.Do(req)
	if err != nil {
		//	是否重定向
		if _, path, ok := judgeIsRedirect(err); ok {
			//	判断host 是否相同
			path1, err := url.Parse(path)
			if err != nil {
				return nil, nil, errors.New("SendRequest: redirect path is error")
			}
			if len(path1.Host) <= 0 {
				path2, err := url.Parse(send_url)
				if err != nil {
					return nil, nil, errors.New("SendRequest: send_url is error")
				}
				path = "http://" + path2.Host + path
			}
			return client.SendRequest("GET", path, data, cookie, header)
		}
		return nil, nil, err
	}
	defer res.Body.Close()
	//	读取回送消息体内容
	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	return resData, res.Cookies(), nil
}

func judgeIsRedirect(in error) (method, path string, ok bool) {
	if in == nil {
		return "", "", false
	}
	re, _ := regexp.Compile(`(Post|Get) (.*): S`)
	maps := re.FindStringSubmatch(in.Error())
	if len(maps) > 0 && len(maps[0]) > 0 {
		return maps[1], maps[2], true
	}
	return "", "", false
}

func getLoginKey(data string) (vierstate, eventvalidation string, err error) {
	if len(data) <= 0 {
		return "", "", errors.New("getLoginKey: data is null")
	}
	re, _ := regexp.Compile(`<input type="hidden" name="__VIEWSTATE" id="__VIEWSTATE" value="(.*)" />`)
	re1, _ := regexp.Compile(`<input type="hidden" name="__EVENTVALIDATION" id="__EVENTVALIDATION" value="(.*)" />`)
	maps := re.FindStringSubmatch(data)
	maps1 := re1.FindStringSubmatch(data)
	if len(maps[0]) > 0 && len(maps1[0]) > 0 {
		return maps[1], maps1[1], nil
	}
	return "", "", errors.New("getLoginKey: not find string")
}
