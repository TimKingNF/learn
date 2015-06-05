package EDU

import (
	"errors"
	"learn/models"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

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

func GetStudentProfile(data []byte, cookies []*http.Cookie) (map[string]string, error) {
	if data == nil {
		return nil, models.ErrorInfo("GetStudentProfile", errors.New("data in nil"))
	}
	userinfohref, _ := getUserInfoHref(string(data))
	userinfo_data, _, err := client.SendRequest("GET", Host_url+"/"+userinfohref, "", cookies, header)
	if err != nil {
		return nil, models.ErrorInfo("GetStudentProfile", err)
	}
	userinfo, _ := getUserInfo(string(userinfo_data))
	userimg, _, err := client.SendRequest("GET", Host_url+"/readimagexs.aspx?xh="+userinfo["学号"], "", cookies, nil)
	if err != nil {
		return nil, models.ErrorInfo("GetStudentProfile", errors.New("readimg is failed"))
	}
	if path, err := models.SaveFile("EDU/headimg/", userinfo["学号"], "jpg", userimg); err == nil {
		userinfo["头像"] = path
	}
	return userinfo, nil
}

func UpdateStudentProfile(user *models.Student, userinfo map[string]string) error {
	if userinfo == nil {
		return models.ErrorInfo("UpdateStudentProfile", errors.New("userinfo is nil"))
	}
	if user == nil || user.Id != userinfo["学号"] {
		return models.ErrorInfo("UpdateStudentProfile", errors.New("用户与学号不匹配"))
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
		if exist := models.DepartmentExist(userinfo["系"]); !exist {
			_, err := models.AddDepartment(&models.Department{Name: userinfo["系"]})
			if err != nil {
				EduLogCreate(&EduLog{Student: user, Content: userinfo["系"], Type: "department", Result: 0})
				return models.ErrorInfo("UpdateStudentProfile", err)
			}
			EduLogCreate(&EduLog{Student: user, Content: userinfo["系"], Type: "department", Result: 1})
		}
	}
	if len(userinfo["专业名称"]) > 0 {
		if exist := models.MajorExist(userinfo["专业名称"]); !exist {
			department, err := models.GetDepartmentByName(userinfo["系"])
			if err != nil {
				return models.ErrorInfo("UpdateStudentProfile", err)
			}
			_, err = models.AddMajor(&models.Major{Name: userinfo["专业名称"], Department: department})
			if err != nil {
				EduLogCreate(&EduLog{Student: user, Content: userinfo["专业名称"], Type: "major", Result: 0})
				return models.ErrorInfo("UpdateStudentProfile", err)
			}
			EduLogCreate(&EduLog{Student: user, Content: userinfo["专业名称"], Type: "major", Result: 1})
		}
	}
	if len(userinfo["行政班"]) > 0 {
		if exist := models.ClassExist(userinfo["行政班"]); !exist {
			department, err := models.GetDepartmentByName(userinfo["系"])
			if err != nil {
				return models.ErrorInfo("UpdateStudentProfile", err)
			}
			major, err := models.GetMajorByName(userinfo["专业名称"])
			if err != nil {
				return models.ErrorInfo("UpdateStudentProfile", err)
			}
			_, err = models.AddClass(&models.Class{Name: userinfo["行政班"], Department: department, Major: major})
			if err != nil {
				EduLogCreate(&EduLog{Student: user, Content: userinfo["行政班"], Type: "class", Result: 0})
				return models.ErrorInfo("UpdateStudentProfile", err)
			}
			EduLogCreate(&EduLog{Student: user, Content: userinfo["行政班"], Type: "class", Result: 1})
		}
	}
	department, _ := models.GetDepartmentByName(userinfo["系"])
	major, _ := models.GetMajorByName(userinfo["专业名称"])
	class, _ := models.GetClassByName(userinfo["行政班"])
	user.Department = department
	user.Major = major
	user.Class = class
	if err := models.UpdateStudent(user); err != nil {
		return models.ErrorInfo("UpdateStudentProfile", err)
	}
	return nil
}
