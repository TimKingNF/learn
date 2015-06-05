package models

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
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
	Edu_Host     string = "jwxt.nfsysu.cn"
	Edu_Host_url string = "http://" + Edu_Host
)

var (
	client *Client           = &Client{Val: new(http.Client)}
	header map[string]string = map[string]string{
		"Host":            Edu_Host,
		"User-Agent":      "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:37.0) Gecko/20100101 Firefox/37.0",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Charset":  "GBK,utf-8;q=0.7,*;q=0.3",
		"Accept-Language": "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
		"Accept-Encoding": "gzip, deflate",
		"Connection":      "keep-alive",
		"Content-Type":    "application/x-www-form-urlencoded; param=value",
		"Referer":         Edu_Host_url,
	}
)

func EduLogin(id, pwd, user_type string) ([]byte, bool, []*http.Cookie, error) {
	if user_type != "学生" && user_type != "教师" {
		return nil, false, nil, ErrorInfo("EduLogin", errors.New("user type is error："+user_type))
	}
	if len(id) <= 0 || len(pwd) <= 0 {
		return nil, false, nil, ErrorInfo("EduLogin", errors.New("id or pwd is error：["+id+"_"+pwd+"]"))
	}
	data, cookies, err := client.SendRequest("GET", Edu_Host_url, "", nil, nil)
	if err != nil {
		return nil, false, nil, ErrorInfo("EduLogin", err)
	}
	vierstate, eventvalidation, err := getLoginKey(string(data))
	if err != nil {
		return nil, false, nil, ErrorInfo("EduLogin", err)
	}
	//	sigin...
	post_arg := url.Values{
		"__VIEWSTATE":       {vierstate},
		"__EVENTVALIDATION": {eventvalidation},
		"TextBox1":          {id},
		"TextBox2":          {pwd},
		"RadioButtonList1":  {utf82gbk(user_type)},
		"Button1":           {""},
	}
	//	登陆
	main_data, _, err := client.SendRequest("POST", Edu_Host_url+"/default2.aspx", post_arg.Encode(), cookies, header)
	if err != nil {
		return nil, false, nil, ErrorInfo("EduLogin", err)
	}
	//	判断是否登录失败
	if strings.Contains(string(main_data), "<title>登录</title>") {
		return nil, false, nil, ErrorInfo("EduLogin", errors.New("id or pwd not match"))
	}
	return main_data, true, cookies, nil
}

func EduGetStudentProfile(data []byte, cookies []*http.Cookie) (map[string]string, error) {
	if data == nil {
		return nil, ErrorInfo("EduGetStudentProfile", errors.New("data in nil"))
	}
	userinfohref, _ := getUserInfoHref(string(data))
	userinfo_data, _, err := client.SendRequest("GET", Edu_Host_url+"/"+userinfohref, "", cookies, header)
	if err != nil {
		return nil, ErrorInfo("EduGetStudentProfile", err)
	}
	userinfo, _ := getUserInfo(string(userinfo_data))
	userimg, _, err := client.SendRequest("GET", Edu_Host_url+"/readimagexs.aspx?xh="+userinfo["学号"], "", cookies, nil)
	if err != nil {
		return nil, ErrorInfo("EduGetStudentProfile", errors.New("readimg is failed"))
	}
	if path, err := SaveFile("EDU/headimg/", userinfo["学号"], "jpg", userimg); err == nil {
		userinfo["头像"] = path
	}
	return userinfo, nil
}

func EduGetStudentSchedule(user *Student, data []byte, cookies []*http.Cookie) (schedules map[*Term][]*TeacherCourse, err error) {
	if data == nil {
		return nil, ErrorInfo("EduGetStudentSchedule", errors.New("data in nil"))
	}
	userschedulehref, _ := getUserScheduleHref(string(data))
	userschedule_data, _, err := client.SendRequest("GET", Edu_Host_url+"/"+userschedulehref, "", cookies, header)
	if err != nil {
		return nil, ErrorInfo("EduGetStudentSchedule", err)
	}
	//	获取当前学年度的学期
	xq, _ := getUserXQ(string(userschedule_data))
	//	获取历年的课程表
	userschedule_xq_href, _ := getUserScheduleXQHref(string(userschedule_data))
	if xqs, err := EduGetXQ(string(userschedule_data)); err == nil {
		vierstate, eventvalidation, _ := getLoginKey(string(userschedule_data))
		//	查询历年课表
		schedules = make(map[*Term][]*TeacherCourse, len(xqs))
		for _, v := range xqs {
			numtemp, _ := strconv.Atoi(xq["学期"])
			numtemp1, _ := strconv.Atoi(v["学期"])
			if v["学年度"] == xq["学年度"] && numtemp < numtemp1 {
				continue
			}
			userschedule_temp_data, _, err := client.SendRequest("POST", Edu_Host_url+"/"+userschedule_xq_href, url.Values{
				"__EVENTTARGET":     {"xnd"},
				"__EVENTARGUMENT":   {""},
				"__LASTFOCUS":       {""},
				"__VIEWSTATE":       {vierstate},
				"__EVENTVALIDATION": {eventvalidation},
				"xnd":               {v["学年度"]},
				"xqd":               {v["学期"]},
			}.Encode(), cookies, header)
			if err != nil {
				fmt.Println(err)
				return nil, ErrorInfo("EduGetStudentSchedule", err)
			}
			if !strings.Contains(string(userschedule_temp_data), `<input type="hidden" name="__VIEWSTATE" id="__VIEWSTATE"`) {
				continue
			}

			vierstate, eventvalidation, _ = getLoginKey(string(userschedule_temp_data))
			userschedule, xq, _ := getUserSchedule(string(userschedule_temp_data))
			if len(xq["学年度"]) > 0 {
				if years := strings.Split(xq["学年度"], "-"); len(years) == 2 {
					startyear, _ := strconv.Atoi(years[0])
					endyear, _ := strconv.Atoi(years[1])
					if len(xq["学期"]) > 0 {
						num, _ := strconv.Atoi(xq["学期"])
						if exist := TermExist(num, startyear, endyear); !exist {
							if err := AddTerm(&Term{Number: num, StartYear: startyear, EndYear: endyear}); err != nil {
								return nil, ErrorInfo("EduGetStudentSchedule", err)
							}
						}
						term, err := SearchTerm(num, startyear, endyear)
						if err != nil {
							return nil, ErrorInfo("EduGetStudentSchedule", err)
						}
						var result []*TeacherCourse
						for _, v := range userschedule {
							if tercher_course, err := CourseTemp2TeacherCourse(term, v); err == nil {
								//	add tercher_course
								if exist := TeacherCourseExist(tercher_course.Course.Id, tercher_course.Term.Id, tercher_course.Teacher.Id); !exist {
									if err := AddTeacherCourse(tercher_course); err != nil {
										Info("EduGetStudentSchedule", err)
									}
								}
								//	search
								if t_course, err := SearchTeacherCourse(tercher_course.Course.Id, tercher_course.Term.Id, tercher_course.Teacher.Id); err == nil {
									// add m2m between teacher_course and class
									if exist := ExistTeacherCourseAndClass(t_course.Id, user.Class.Id); !exist {
										AddM2MBetweenTeacherCourseAndClass(t_course.Id, user.Class.Id)
									}
									result = append(result, t_course)
								}
							}
						}
						schedules[term] = result
					}
				}
			}
		}
	}
	return schedules, nil
}

func EduGetStudentScore(user *Student, data []byte, cookies []*http.Cookie) (map[string]map[string]string, error) {
	userscorehref, _ := getUserScoreHref(string(data))
	userscore_param, _, err := client.SendRequest("GET", Edu_Host_url+"/"+userscorehref, "", cookies, header)
	if err != nil {
		return nil, ErrorInfo("EduGetStudentScore", errors.New("连接失败,错误原因："+err.Error()))
	}
	vierstate, eventvalidation, _ := getLoginKey(string(userscore_param))
	userscore_info, _, err := client.SendRequest("POST", Edu_Host_url+"/"+userscorehref, url.Values{
		"__VIEWSTATE":       {vierstate},
		"__EVENTVALIDATION": {eventvalidation},
		"ddIXN":             {""},
		"ddIXQ":             {""},
		"Button1":           {utf82gbk("按学期查询")},
	}.Encode(), cookies, header)
	if err != nil {
		return nil, ErrorInfo("EduGetStudentScore", errors.New("连接失败,错误原因："+err.Error()))
	}
	userscore, _ := getUserScore(string(userscore_info))
	for _, v := range userscore {
		if len(v) > 0 {
			if exist := CourseExist(v["课程名称"]); !exist {
				if exist := DepartmentExist(v["学院名称"]); !exist {
					AddDepartment(&Department{Name: v["学院名称"]})
				}
				if department, err := GetDepartmentByName(v["学院名称"]); err == nil {
					f, _ := strconv.ParseFloat(v["学分"], 64)
					AddCourse(&Course{Id: v["课程代码"], Name: v["课程名称"], Department: department, Remark: v["课程归属"], Type: v["课程性质"], Credit: f})
				}
			}
		}
	}
	return userscore, nil
}

func EduUpdateStudentProfile(user *Student, userinfo map[string]string) error {
	if userinfo == nil {
		return ErrorInfo("EduUpdateStudentProfile", errors.New("userinfo is nil"))
	}
	if user == nil || user.Id != userinfo["学号"] {
		return ErrorInfo("EduUpdateStudentProfile", errors.New("用户与学号不匹配"))
	}
	user.Name = userinfo["姓名"]
	user.IsEdu = true
	if userinfo["性别"] == "男" {
		user.Sex = 1
	} else if userinfo["性别"] == "女" {
		user.Sex = 2
	}
	user.Headimgurl = userinfo["头像"]
	enterSchoolYear, _ := strconv.Atoi(userinfo["当前所在级"])
	user.EnterSchoolYear = enterSchoolYear
	enterSchoolDate, _ := strconv.Atoi(userinfo["入学日期"])
	user.EnterSchoolDate = enterSchoolDate
	idcard, _ := strconv.ParseInt(userinfo["身份证号"], 10, 64)
	user.IdCard = idcard

	if len(userinfo["系"]) > 0 {
		if exist := DepartmentExist(userinfo["系"]); !exist {
			_, err := AddDepartment(&Department{Name: userinfo["系"]})
			if err != nil {
				return ErrorInfo("EduUpdateStudentProfile", err)
			}
		}
	}
	if len(userinfo["专业名称"]) > 0 {
		if exist := MajorExist(userinfo["专业名称"]); !exist {
			department, err := GetDepartmentByName(userinfo["系"])
			if err != nil {
				return ErrorInfo("EduUpdateStudentProfile", err)
			}
			_, err = AddMajor(&Major{Name: userinfo["专业名称"], Department: department})
			if err != nil {
				return ErrorInfo("EduUpdateStudentProfile", err)
			}
		}
	}
	if len(userinfo["行政班"]) > 0 {
		if exist := ClassExist(userinfo["行政班"]); !exist {
			department, err := GetDepartmentByName(userinfo["系"])
			if err != nil {
				return ErrorInfo("EduUpdateStudentProfile", err)
			}
			major, err := GetMajorByName(userinfo["专业名称"])
			if err != nil {
				return ErrorInfo("EduUpdateStudentProfile", err)
			}
			_, err = AddClass(&Class{Name: userinfo["行政班"], Department: department, Major: major})
			if err != nil {
				return ErrorInfo("EduUpdateStudentProfile", err)
			}
		}
	}
	department, _ := GetDepartmentByName(userinfo["系"])
	major, _ := GetMajorByName(userinfo["专业名称"])
	class, _ := GetClassByName(userinfo["行政班"])
	user.Department = department
	user.Major = major
	user.Class = class
	if err := UpdateStudent(user); err != nil {
		return ErrorInfo("EduUpdateStudentProfile", err)
	}
	return nil
}

// func EduAddStudentProfile(userinfo map[string]string) error {
// 	if userinfo == nil {
// 		return ErrorInfo("EduAddStudentProfile", errors.New("userinfo is nil"))
// 	}
// 	var user = &Student{Id: userinfo["学号"], Name: userinfo["姓名"], IsEdu: true}
// 	if userinfo["性别"] == "男" {
// 		user.Sex = 1
// 	} else if userinfo["性别"] == "女" {
// 		user.Sex = 2
// 	}
// 	user.Headimgurl = userinfo["头像"]
// 	enterSchoolYear, _ := strconv.Atoi(userinfo["当前所在级"])
// 	user.EnterSchoolYear = enterSchoolYear
// 	enterSchoolDate, _ := strconv.Atoi(userinfo["入学日期"])
// 	user.EnterSchoolDate = enterSchoolDate
// 	idcard, _ := strconv.ParseInt(userinfo["身份证号"], 10, 64)
// 	user.IdCard = idcard

// 	if len(userinfo["系"]) > 0 {
// 		if exist := DepartmentExist(userinfo["系"]); !exist {
// 			_, err := AddDepartment(&Department{Name: userinfo["系"]})
// 			if err != nil {
// 				return ErrorInfo("EduAddStudentProfile", err)
// 			}
// 		}
// 	}
// 	if len(userinfo["专业名称"]) > 0 {
// 		if exist := MajorExist(userinfo["专业名称"]); !exist {
// 			department, err := GetDepartmentByName(userinfo["系"])
// 			if err != nil {
// 				return ErrorInfo("EduAddStudentProfile", err)
// 			}
// 			_, err = AddMajor(&Major{Name: userinfo["专业名称"], Department: department})
// 			if err != nil {
// 				return ErrorInfo("EduAddStudentProfile", err)
// 			}
// 		}
// 	}
// 	if len(userinfo["行政班"]) > 0 {
// 		if exist := ClassExist(userinfo["行政班"]); !exist {
// 			department, err := GetDepartmentByName(userinfo["系"])
// 			if err != nil {
// 				return ErrorInfo("EduAddStudentProfile", err)
// 			}
// 			major, err := GetMajorByName(userinfo["专业名称"])
// 			if err != nil {
// 				return ErrorInfo("EduAddStudentProfile", err)
// 			}
// 			_, err = AddClass(&Class{Name: userinfo["行政班"], Department: department, Major: major})
// 			if err != nil {
// 				return ErrorInfo("EduAddStudentProfile", err)
// 			}
// 		}
// 	}
// 	department, _ := GetDepartmentByName(userinfo["系"])
// 	major, _ := GetMajorByName(userinfo["专业名称"])
// 	class, _ := GetClassByName(userinfo["行政班"])
// 	user.Department = department
// 	user.Major = major
// 	user.Class = class
// 	if err := AddStudent(user); err != nil {
// 		return ErrorInfo("EduAddStudentProfile", err)
// 	}
// 	return nil
// }

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

func getUserInfoHref(data string) (href string, err error) {
	if len(data) <= 0 {
		return "", errors.New("getUserInfoHref: data is null")
	}
	re, _ := regexp.Compile(`信息维护</span><!--\[if gte IE 7\]><!--></a><!--<!\[endif\]--><!--\[if lte IE 6\]><table><tr><td><!\[endif\]--><ul class='sub'><!--\[if lte IE 6\]><iframe class='navbug'></iframe><!\[endif\]--><li><a href="(.*)" target='zhuti' onclick=".*">个人信息</a></li>`)
	maps := re.FindStringSubmatch(data)
	if len(maps) > 0 && len(maps[0]) > 0 {
		return maps[1], nil
	}
	return "", errors.New("getUserInfoHref: not find string")
}

func getUserInfo(data string) (info map[string]string, err error) {
	if len(data) <= 0 {
		return nil, errors.New("getUserInfo: data is null")
	}
	//	取出table 中所有内容
	re, _ := regexp.Compile(`<table class="formlist" width="100%" align="center">([\s\S]*)</table>`)
	maps := re.FindStringSubmatch(data)
	if len(maps) <= 0 {
		return nil, errors.New("getUserInfo: table is null")
	}
	//	去除所有html标签
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	data = re.ReplaceAllString(maps[1], "\n")
	//	去除多余换行符 和空格
	re, _ = regexp.Compile("\\s{2,}")
	data = re.ReplaceAllString(data, "\n")
	data = strings.Replace(data, "&nbsp;", "", -1)
	//	根据换行符分割data
	arr := strings.Split(data, "\n")
	if len(arr) < 0 {
		return nil, errors.New("getUserInfo: userinfo is null")
	}
	var length = math.Ceil(float64(len(arr) / 2))
	info = make(map[string]string, int(length))
	for k := 0; k < len(arr); k++ {
		if strings.LastIndex(arr[k], "：") == len(arr[k])-3 && strings.LastIndex(arr[k+1], "：") < 0 {
			info[strings.Replace(arr[k], "：", "", -1)] = arr[k+1]
		}
	}
	return info, nil
}

func getUserScheduleHref(data string) (href string, err error) {
	if len(data) <= 0 {
		return "", errors.New("getUserScheduleHref: data is null")
	}
	re, _ := regexp.Compile(`专业推荐课表查询</a></li><li><a href="(.*)" target='zhuti' onclick=".*">学生个人课表</a>`)
	maps := re.FindStringSubmatch(data)
	if len(maps) > 0 && len(maps[0]) > 0 {
		return maps[1], nil
	}
	return "", errors.New("getUserScheduleHref: not find string")
}

func getUserScheduleXQHref(data string) (href string, err error) {
	if len(data) <= 0 {
		return "", errors.New("getUserScheduleXQHref: data is null")
	}
	re, _ := regexp.Compile(`<form name="xskb_form" method="post" action="(.*)" id="xskb_form">`)
	maps := re.FindStringSubmatch(data)
	if len(maps) > 0 && len(maps[0]) > 0 {
		return maps[1], nil
	}
	return "", errors.New("getUserScheduleXQHref: not find string")
}

func getUserSchedule(data string) (info [][]string, xq map[string]string, err error) {
	if len(data) <= 0 {
		return nil, nil, errors.New("getUserSchedule: data is null")
	}
	xq = make(map[string]string, 2)
	//	取出学年度数据
	re, _ := regexp.Compile(`<option selected="selected" value="([\d-]*)">[\d-]*</option>`)
	maps := re.FindStringSubmatch(data)
	if len(maps) <= 0 {
		return nil, nil, errors.New("getUserSchedule: school year is null")
	}
	xq["学年度"] = maps[1]
	//	取出学期数据
	re, _ = regexp.Compile(`<option selected="selected" value="([\d])">[\d]</option>`)
	maps = re.FindStringSubmatch(data)
	if len(maps) < 0 {
		return nil, nil, errors.New("getUserSchedule: term is null")
	}
	xq["学期"] = maps[1]
	//	取出table内容
	re, _ = regexp.Compile(`<td colspan="2" rowspan="1" style="width:2%;">时间</td><td align="center" style="width:14%;">星期一</td><td align="center" style="width:14%;">星期二</td><td align="center" style="width:14%;">星期三</td><td align="center" style="width:14%;">星期四</td><td align="center" style="width:14%;">星期五</td><td align="center" style="width:14%;">星期六</td><td align="center" style="width:14%;">星期日</td>([\s\S]*)</table>[\s\S]*<br>`)
	maps = re.FindStringSubmatch(data)
	if len(maps) <= 0 {
		return nil, xq, errors.New("getUserSchedule: table is null")
	}
	//	去除所有html标签
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	data = re.ReplaceAllString(maps[1], "\n")
	data = strings.Replace(data, "&nbsp;", "", -1)
	//	根据换行符分割data
	arr := strings.Split(data, "\n")
	if len(arr) < 0 {
		return nil, xq, errors.New("getUserSchedule: data is null")
	}
	for k, _ := range arr {
		if k >= len(arr)-3 {
			break
		}
		if len(arr[k]) > 0 && len(arr[k+1]) > 0 && len(arr[k+2]) > 0 && len(arr[k+3]) > 0 {
			info = append(info, []string{arr[k], arr[k+1], arr[k+2], arr[k+3]})
		}
	}
	return info, xq, nil
}

func getUserScoreHref(data string) (href string, err error) {
	if len(data) <= 0 {
		return "", errors.New("getUserScoreHref: data is null")
	}
	re, _ := regexp.Compile(`学生考试查询</a></li><li><a href="(.*)" target='zhuti' onclick=".*">成绩查询</a></li>`)
	maps := re.FindStringSubmatch(data)
	if len(maps) > 0 && len(maps[0]) > 0 {
		return maps[1], nil
	}
	return "", errors.New("getUserScoreHref: not find string")
}

func getUserScore(data string) (info map[string]map[string]string, err error) {
	if len(data) <= 0 {
		return nil, errors.New("getUserScore: data is null")
	}
	//	取出table内容
	re, _ := regexp.Compile(`<table class="datelist" cellspacing="0" cellpadding="3" border="0" id="Datagrid1" style="width:100%;border-collapse:collapse;">([\s\S]*)</table>[\s\S]*<table width="100%">`)
	maps := re.FindStringSubmatch(data)
	if len(maps) <= 0 {
		return nil, errors.New("getUserScore: table is null")
	}
	//	去除所有html标签
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	data = re.ReplaceAllString(maps[1], "\n")
	data = strings.Replace(data, "&nbsp;", "-", -1)
	//	去除多余换行符 和空格
	re, _ = regexp.Compile("\\s{2,}")
	data = re.ReplaceAllString(data, "\n")
	//	根据换行符分割data
	arr := strings.Split(data, "\n")
	if len(arr) < 0 {
		return nil, errors.New("getUserScore: data is null")
	}
	arr = arr[1:]
	info = make(map[string]map[string]string, len(arr)-1)
	for k, _ := range arr {
		if k >= 15 && k%15 == 0 {
			k1 := k + 15
			if k1 >= len(arr)-15 {
				break
			}
			var s = make(map[string]string, k1-k)
			for i := 0; k < k1; k++ {
				s[arr[i]] = arr[k]
				i++
			}
			info[s["课程代码"]] = s
		}
	}
	return info, nil
}

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

func CourseTemp2TeacherCourse(xq *Term, in []string) (*TeacherCourse, error) {
	if len(in) != 4 || xq == nil {
		return nil, ErrorInfo("CourseTemp2TeacherCourse", errors.New("in or xq is nil"))
	}
	courseName := in[0]
	time := in[1]
	teacherName := in[2]
	if exist := CourseExist(courseName); !exist {
		if err := AddCourse(&Course{Name: courseName}); err != nil {
			return nil, ErrorInfo("CourseTemp2TeacherCourse", err)
		}
	}
	course, err := GetCourseByName(courseName)
	if err != nil {
		return nil, ErrorInfo("CourseTemp2TeacherCourse", err)
	}
	if exist := TeacherExistByDep(teacherName, course.Department.Id); !exist {
		if err := AddTeacher(&Teacher{Name: teacherName, Password: "111111", Department: course.Department, Headimgurl: "/static/img/avatar.jpeg"}); err != nil {
			return nil, ErrorInfo("CourseTemp2TeacherCourse", err)
		}
	}
	teacher, err := SearchTeacher(teacherName, course.Department.Id)
	if err != nil {
		return nil, ErrorInfo("CourseTemp2TeacherCourse", err)
	}
	if courseTime := EduGetTime(time); courseTime != nil {
		end_week, _ := strconv.Atoi(courseTime["end_week"])
		start_week, _ := strconv.Atoi(courseTime["start_week"])
		return &TeacherCourse{Course: course, Term: xq, Teacher: teacher, Place: in[3], StartWeek: start_week, EndWeek: end_week, Time: `{"week_day":"` + courseTime["week_day"] + `", "week_time":"` + courseTime["time"] + `"}`}, nil
	}
	return nil, ErrorInfo("CourseTemp2TeacherCourse", "操作失败")
}

func EduGetTime(date string) map[string]string {
	if len(date) <= 0 {
		return nil
	}
	var result = make(map[string]string)
	time := date[:6]
	result["week_day"] = fmt.Sprintf("%d", SwitchDate(time))
	re, _ := regexp.Compile("(\\d){1,2}")
	maps := re.FindAllString(date, 10)
	if len(maps) > 0 {
		result["end_week"] = strings.Join(maps[len(maps)-1:], "")
		result["start_week"] = strings.Join(maps[len(maps)-2:len(maps)-1], "")
		result["time"] = strings.Join(maps[:len(maps)-2], ",")
	}
	return result
}

func EduGetXQ(data string) (xqs []map[string]string, err error) {
	if len(data) <= 0 {
		return nil, errors.New("EduGetXQ: data is null")
	}
	//	取出学年度数据
	re, _ := regexp.Compile(`<option( selected="selected"){0,1} value="([\d-]*)">[\d-]*</option>`)
	maps := re.FindAllStringSubmatch(data, 4)
	for _, v := range maps {
		if len(v[2]) > 1 {
			xqs = append(xqs, map[string]string{"学年度": v[2], "学期": "1"})
			xqs = append(xqs, map[string]string{"学年度": v[2], "学期": "2"})
		}
	}
	return xqs, nil
}

func getUserXQ(data string) (xq map[string]string, err error) {
	if len(data) <= 0 {
		return nil, errors.New("getUserXQ: data is null")
	}
	xq = make(map[string]string, 2)
	//	取出学年度数据
	re, _ := regexp.Compile(`<option selected="selected" value="([\d-]*)">[\d-]*</option>`)
	maps := re.FindStringSubmatch(data)
	if len(maps) < 0 {
		return nil, errors.New("getUserXQ: school year is null")
	}
	xq["学年度"] = maps[1]
	//	取出学期数据
	re, _ = regexp.Compile(`<option selected="selected" value="([\d])">[\d]</option>`)
	maps = re.FindStringSubmatch(data)
	if len(maps) < 0 {
		return nil, errors.New("getUserXQ: term is null")
	}
	xq["学期"] = maps[1]
	return xq, nil
}

func EduControlSQLByScheduleAndScore(user *Student, schedules map[*Term][]*TeacherCourse, score map[string]map[string]string) error {
	if user == nil || len(schedules) <= 0 || len(score) <= 0 {
		return ErrorInfo("EduControlSQLByScheduleAndScore", errors.New("user is nil"))
	}
	for k, v := range schedules {
		//	取出学年度 学期
		if k == nil || v == nil {
			continue
		}
		if exist := TermExist(k.Number, k.StartYear, k.EndYear); !exist {
			AddTerm(k)
		}
		_, err := SearchTerm(k.Number, k.StartYear, k.EndYear)
		if err != nil {
			Info("EduControlSQLByScheduleAndScore", err)
			continue
		}
		for _, t_course := range v {
			if exist := TeacherCourseExist(t_course.Course.Id, t_course.Term.Id, t_course.Teacher.Id); !exist {
				//	add teacher course
				AddTeacherCourse(t_course)
			}
			if temp_t_course, err := SearchTeacherCourse(t_course.Course.Id, t_course.Term.Id, t_course.Teacher.Id); err == nil {
				if exist := StudentCourseExist(user.Id, temp_t_course.Id); !exist {
					//	add student course
					AddStudentCourse(&StudentCourse{Student: &Student{Id: user.Id}, TeacherCourse: temp_t_course})
				}
				if temp_s_course, err := SearchStudentCourse(user.Id, temp_t_course.Id); err == nil {
					if exist := StudentGradeExist(temp_s_course.Id); !exist {
						//	取出成绩
						var (
							score = score[t_course.Course.Id]
							grade = 0
						)
						if score["重修成绩"] != "-" {
							num, _ := strconv.Atoi(score["重修成绩"])
							grade = num
						} else if score["补考成绩"] != "-" {
							num, _ := strconv.Atoi(score["补考成绩"])
							grade = num
						} else if len(score["成绩"]) > 0 {
							num, _ := strconv.Atoi(score["成绩"])
							grade = num
						}
						gpa, _ := strconv.ParseFloat(score["绩点"], 64)
						AddStudentGrade(&StudentGrade{StudentCourse: temp_s_course, GradePointAverage: gpa, Grade: grade})
					}
				}
			}
		}
	}
	return nil
}
