package EDU

import (
	"errors"
	"fmt"
	"learn/models"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func getUserScheduleHref(data string) (href string, err error) {
	if len(data) <= 0 {
		return "", errors.New("getUserScheduleHref: data is null")
	}
	re, _ := regexp.Compile(`培养计划</a></li><li><a href="(.*)" target='zhuti' onclick=".*">学生选课情况查询</a></li>`)
	maps := re.FindStringSubmatch(data)
	if len(maps) > 0 && len(maps[0]) > 0 {
		return maps[1], nil
	}
	return "", errors.New("getUserScheduleHref: not find string")
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

func getUserScheduleXQHref(data string) (href string, err error) {
	if len(data) <= 0 {
		return "", errors.New("getUserScheduleXQHref: data is null")
	}
	re, _ := regexp.Compile(`<form name="Form1" method="post" action="(.*)" id="Form1">`)
	maps := re.FindStringSubmatch(data)
	if len(maps) > 0 && len(maps[0]) > 0 {
		return maps[1], nil
	}
	return "", errors.New("getUserScheduleXQHref: not find string")
}

func GetStudentSchedule(user *models.Student, data []byte, cookies []*http.Cookie) (schedules map[*models.Term][]*models.TeacherCourse, err error) {
	if data == nil {
		return nil, models.ErrorInfo("GetStudentSchedule", errors.New("data in nil"))
	}
	userschedulehref, _ := getUserScheduleHref(string(data))
	userschedule_data, _, err := client.SendRequest("GET", Host_url+"/"+userschedulehref, "", cookies, header)
	if err != nil {
		return nil, models.ErrorInfo("GetStudentSchedule", err)
	}

	//	获取当前学年度的学期
	xq, _ := getUserXQ(string(userschedule_data))
	//	根据用户的入学年份和当前学年度的学期 列出学期列表
	xqs, err := getXQS(xq, user.EnterSchoolYear)
	if err != nil {
		return nil, models.ErrorInfo("GetStudentSchedule", err)
	}
	//	查询历年课表
	schedules = make(map[*models.Term][]*models.TeacherCourse, len(xqs))
	userschedule_xq_href, _ := getUserScheduleXQHref(string(userschedule_data))
	vierstate, eventvalidation, _ := getLoginKey(string(userschedule_data))
	for _, v := range xqs {
		//	判断学期是否存在 如果不存在则添加
		if len(v["学年度"]) > 0 {
			if years := strings.Split(v["学年度"], "-"); len(years) == 2 {
				startyear, _ := strconv.Atoi(years[0])
				endyear, _ := strconv.Atoi(years[1])
				if len(v["学期"]) > 0 {
					num, _ := strconv.Atoi(v["学期"])
					if exist := models.TermExist(num, startyear, endyear); !exist {
						if err := models.AddTerm(&models.Term{Number: num, StartYear: startyear, EndYear: endyear}); err == nil {
							EduLogCreate(&EduLog{Student: user, Content: `{"学年度": "` + v["学年度"] + `", "学期": "` + v["学期"] + `"}`, Type: "term", Result: 1})
						} else {
							EduLogCreate(&EduLog{Student: user, Content: `{"学年度": "` + v["学年度"] + `", "学期": "` + v["学期"] + `"}`, Type: "term", Result: 0})
						}
					}
					//	查询学期
					term, err := models.SearchTerm(num, startyear, endyear)
					if err != nil {
						models.Info("EduGetStudentSchedule", err)
						continue
					}
					//	获取login key
					userschedule_temp_data, _, err := client.SendRequest("POST", Host_url+"/"+userschedule_xq_href, url.Values{
						"__EVENTTARGET":     {"ddlXQ"},
						"__EVENTARGUMENT":   {""},
						"__LASTFOCUS":       {""},
						"__VIEWSTATE":       {vierstate},
						"__EVENTVALIDATION": {eventvalidation},
						"ddlXN":             {fmt.Sprintf("%d-%d", term.StartYear, term.EndYear)},
						"ddlXQ":             {fmt.Sprintf("%d", term.Number)},
					}.Encode(), cookies, header)
					if err != nil {
						EduLogCreate(&EduLog{Student: user, Content: `获取` + v["学年度"] + ` ` + v["学期"] + `课表失败`, Type: "t_course", Result: 0})
						continue
					}
					if !strings.Contains(string(userschedule_temp_data), `<input type="hidden" name="__EVENTTARGET" id="__EVENTTARGET"`) {
						EduLogCreate(&EduLog{Student: user, Content: `获取` + v["学年度"] + ` ` + v["学期"] + `课表失败`, Type: "t_course", Result: 0})
						continue
					}
					usercourse, _, _ := getUserCourse(string(userschedule_temp_data), cookies)
					if usercourse == nil {
						EduLogCreate(&EduLog{Student: user, Content: `获取` + v["学年度"] + ` ` + v["学期"] + `课表失败`, Type: "t_course", Result: 0})
						continue
					}
					// schedules[term] = make([]*models.TeacherCourse, len(usercourse))
					for _, v := range usercourse {
						if len(v) > 0 {
							if exist := models.CourseExist(v["课程名称"]); !exist {
								if exist := models.DepartmentExist(v["教师所在系"]); !exist {
									if _, err := models.AddDepartment(&models.Department{Name: v["教师所在系"]}); err == nil {
										EduLogCreate(&EduLog{Student: user, Content: v["教师所在系"], Type: "department", Result: 1})
									} else {
										EduLogCreate(&EduLog{Student: user, Content: v["教师所在系"], Type: "department", Result: 0})
									}
								}
								if department, err := models.GetDepartmentByName(v["教师所在系"]); err == nil {
									if exist := models.TeacherExistByDep(v["教师姓名"], department.Id); !exist {
										if err := models.AddTeacher(&models.Teacher{Name: v["教师姓名"], Department: department, Password: "111111", Headimgurl: "/static/img/avatar.jpeg"}); err == nil {
											EduLogCreate(&EduLog{Student: user, Content: v["教师姓名"], Type: "teacher", Result: 1})
										} else {
											EduLogCreate(&EduLog{Student: user, Content: v["教师姓名"], Type: "teacher", Result: 0})
										}
									}
									f, _ := strconv.ParseFloat(v["学分"], 64)
									if err := models.AddCourse(&models.Course{Id: v["课程代码"], Name: v["课程名称"], Department: department, Type: v["课程性质"], Credit: f}); err == nil {
										EduLogCreate(&EduLog{Student: user, Content: v["课程名称"], Type: "course", Result: 1})
									} else {
										EduLogCreate(&EduLog{Student: user, Content: v["课程名称"], Type: "course", Result: 0})
									}
								}
							}
							if department, err := models.GetDepartmentByName(v["教师所在系"]); err == nil {
								if teacher, err := models.SearchTeacher(v["教师姓名"], department.Id); err == nil {
									//	search course
									if course, err := models.SearchCourse(v["课程名称"]); err == nil {
										if courseTime := getTime(v["上课时间"]); courseTime != nil {
											for _, v1 := range courseTime {
												end_week, _ := strconv.Atoi(v1["end_week"])
												start_week, _ := strconv.Atoi(v1["start_week"])
												spacing, _ := strconv.Atoi(v1["spacing"])
												time := `{"week_day":"` + v1["week_day"] + `", "week_time":"` + v1["time"] + `"}`
												//	add m2m between teacher and course
												if exist := models.TeacherCourseExist(course.Id, time, term.Id, teacher.Id); !exist {
													var teacher_course = &models.TeacherCourse{Course: course,
														Term:      term,
														Teacher:   teacher,
														Place:     v["上课地点"],
														StartWeek: start_week,
														EndWeek:   end_week,
														Time:      time,
														Remark:    v1["remark"],
														Spacing:   spacing,
													}
													if err := models.AddTeacherCourse(teacher_course); err == nil {
														EduLogCreate(&EduLog{Student: user, Content: `教师 ` + teacher_course.Teacher.Name + ` 教授的课程：` + teacher_course.Course.Name + `, 上课时间：` + v["上课时间"] + `, 上课地点：` + v["上课地点"], Type: "t_course", Result: 1})
													} else {
														models.Info("EduGetStudentSchedule", err)
														EduLogCreate(&EduLog{Student: user, Content: `教师 ` + teacher_course.Teacher.Name + ` 教授的课程：` + teacher_course.Course.Name + `, 上课时间：` + v["上课时间"] + `, 上课地点：` + v["上课地点"], Type: "t_course", Result: 0})
													}
												}
												//	search teacher course
												if t_course, err := models.SearchTeacherCourse(course.Id, time, term.Id, teacher.Id); err == nil {
													// add m2m between teacher_course and class
													if exist := models.ExistTeacherCourseAndClass(t_course.Id, user.Class.Id); !exist {
														if err := models.AddM2MBetweenTeacherCourseAndClass(t_course.Id, user.Class.Id); err == nil {
															EduLogCreate(&EduLog{Student: user, Content: `教师 ` + t_course.Teacher.Name + ` 教授的课程：` + t_course.Course.Name + `,班级：` + user.Class.Name, Type: "t_course_class", Result: 1})
														} else {
															EduLogCreate(&EduLog{Student: user, Content: `教师 ` + t_course.Teacher.Name + ` 教授的课程：` + t_course.Course.Name + `,班级：` + user.Class.Name, Type: "t_course_class", Result: 0})
														}
													}
													//	add m2m between teacher_course and student
													if exist := models.StudentCourseExist(user.Id, t_course.Id); !exist {
														//	add student course
														if err := models.AddStudentCourse(&models.StudentCourse{Student: user, TeacherCourse: t_course}); err == nil {
															EduLogCreate(&EduLog{Student: user, Content: `教师 ` + t_course.Teacher.Name + ` 教授的课程：` + t_course.Course.Name + `, 上课时间：` + v["上课时间"] + `, 上课地点：` + v["上课地点"], Type: "s_course", Result: 1})
														} else {
															EduLogCreate(&EduLog{Student: user, Content: `教师 ` + t_course.Teacher.Name + ` 教授的课程：` + t_course.Course.Name + `, 上课时间：` + v["上课时间"] + `, 上课地点：` + v["上课地点"], Type: "s_course", Result: 0})
														}
													}
													schedules[term] = append(schedules[term], t_course)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return schedules, nil
}

func getXQS(xq map[string]string, enterSchoolYear int) (xqs []map[string]string, err error) {
	if xq == nil || enterSchoolYear <= 0 {
		return nil, errors.New("getXQS: data is null")
	}
	//	get the max xnd
	arr := strings.Split(xq["学年度"], "-")
	var max = arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}
	max_xnd, _ := strconv.Atoi(max)
	for s := enterSchoolYear; s < max_xnd; s++ {
		if s == max_xnd-1 && xq["学期"] == "1" {
			xqs = append(xqs, map[string]string{"学年度": fmt.Sprintf("%d", s) + "-" + fmt.Sprintf("%d", s+1), "学期": "1"})
		} else {
			xqs = append(xqs, map[string]string{"学年度": fmt.Sprintf("%d", s) + "-" + fmt.Sprintf("%d", s+1), "学期": "1"})
			xqs = append(xqs, map[string]string{"学年度": fmt.Sprintf("%d", s) + "-" + fmt.Sprintf("%d", s+1), "学期": "2"})
		}
	}
	return xqs, nil
}

func getUserCourse(data string, cookies []*http.Cookie) (info []map[string]string, xq map[string]string, err error) {
	if len(data) <= 0 {
		return nil, nil, errors.New("getUserCourse: data is null")
	}
	xq = make(map[string]string, 2)
	//	取出学年度数据
	re, _ := regexp.Compile(`<option selected="selected" value="([\d-]*)">[\d-]*</option>`)
	maps := re.FindStringSubmatch(data)
	if len(maps) <= 0 {
		return nil, nil, errors.New("getUserCourse: school year is null")
	}
	xq["学年度"] = maps[1]
	//	取出学期数据
	re, _ = regexp.Compile(`<option selected="selected" value="([\d])">[\d]</option>`)
	maps = re.FindStringSubmatch(data)
	if len(maps) < 0 {
		return nil, nil, errors.New("getUserCourse: term is null")
	}
	xq["学期"] = maps[1]
	//	取出table内容
	re, _ = regexp.Compile(`<table class="datelist " cellspacing="0" cellpadding="3" border="0" id="DBGrid" style="width:100%;border-collapse:collapse;">([\s\S]*)</table>[\s\S]*</span>`)
	maps = re.FindStringSubmatch(data)
	if len(maps) <= 0 {
		return nil, nil, errors.New("getUserScore: table is null")
	}
	var key []string
	//	按 </tr> 截取
	arr := strings.Split(maps[1], "</tr>")
	for k, v := range arr {
		//	取出td中的内容
		re, _ = regexp.Compile(`<td>([\s\S]*)</td><td>([\s\S]*)</td><td>([\s\S]*)</td><td>([\s\S]*)</td><td>([\s\S]*)</td><td>([\s\S]*)</td><td>([\s\S]*)</td><td( title="([\s\S]*)"){0,1}>([\s\S]*)</td><td>([\s\S]*)</td><td>([\s\S]*)</td><td>([\s\S]*)</td><td>([\s\S]*)</td><td>([\s\S]*)</td><td>([\s\S]*)</td><td>([\s\S]*)</td>`)
		maps1 := re.FindStringSubmatch(v)
		if len(maps1) > 0 {
			if k == 0 {
				key = maps1[1:]
			} else {
				var s = make(map[string]string, 11)
				for k1, v1 := range maps1 {
					if k1 == 0 {
						continue
					}
					data = strings.Replace(data, "&nbsp;", "-", -1)
					if key[k1-1] == "上课时间" {
						if string(v1[len(v1)-1]) != "}" {
							s[key[k1-1]] = maps1[k1-1]
						} else {
							s[key[k1-1]] = v1
						}
					} else {
						s[key[k1-1]] = v1
					}
				}
				info = append(info, s)
			}
		}
	}
	for k, v := range info {
		re1, _ := regexp.Compile(`<a href='#' onclick="window.open\('kcxx\.aspx\?kcdm=(.*)','kcxx','toolbar=0,location=0,directories=0,status=0,menubar=0,scrollbars=1,resizable=0,width=490,height=500,left=200,top=50'\)">(.*)</a>`)
		maps2 := re1.FindStringSubmatch(v["课程名称"])
		if len(maps2) <= 0 {
			continue
		}
		info[k]["课程代码"] = maps2[1]
		info[k]["课程名称"] = maps2[2]
		//	获取教师姓名与教师所在学院
		re, _ = regexp.Compile(`<a href='#' onclick="window.open\('(.*)','jsxx','toolbar=0,location=0,directories=0,status=0,menubar=0,scrollbars=1,resizable=0,width=800,height=600,left=120,top=60'\)">(.*)</a>`)
		maps3 := re.FindStringSubmatch(v["教师姓名"])
		if len(maps3) <= 0 {
			continue
		}
		info[k]["教师姓名"] = maps3[2]
		//	发送连接获取教师所在的系
		data, _, err := client.SendRequest("GET", Host_url+"/"+maps3[1], "", cookies, header)
		if err != nil {
			return nil, nil, models.ErrorInfo("getUserCourse", err)
		}
		re3, _ := regexp.Compile(`<TD><span id="bm">(.*)</span></TD>`)
		maps4 := re3.FindStringSubmatch(string(data))
		if len(maps4) <= 0 {
			continue
		}
		info[k]["教师所在系"] = maps4[1]
	}
	return info, xq, nil
}

func getTime(date string) (result []map[string]string) {
	if len(date) <= 0 {
		return nil
	}
	arr := strings.Split(date, ";")
	for _, v := range arr {
		//	serach start_week and end_week
		re, _ := regexp.Compile("(\\d){1,2}")
		maps_week := re.FindAllString(v, 10)
		start_week := strings.Join(maps_week[len(maps_week)-1:], "")
		end_week := strings.Join(maps_week[len(maps_week)-2:len(maps_week)-1], "")
		remark := ""
		//	is spacing ?
		if string(v[len(v)-8]) == "|" {
			remark = v[len(v)-7 : len(v)-1]
		}
		//	serach week_day
		re, _ = regexp.Compile("周[一|二|三|四|五|六|日]第[\\d,]*节")
		maps := re.FindAllString(v, 10)
		for _, v := range maps {
			var result_tmp = make(map[string]string)
			time := v[:6]
			result_tmp["week_day"] = fmt.Sprintf("%d", models.SwitchDate(time))
			re, _ := regexp.Compile("(\\d){1,2}")
			maps := re.FindAllString(v, 10)
			if len(maps) > 0 {
				result_tmp["time"] = strings.Join(maps[:len(maps)], ",")
				result_tmp["end_week"] = start_week
				result_tmp["start_week"] = end_week
				result_tmp["remark"] = remark
				switch remark {
				case "":
					result_tmp["spacing"] = fmt.Sprintf("%d", 0)
				case "单周":
					result_tmp["spacing"] = fmt.Sprintf("%d", 1)
				case "双周":
					result_tmp["spacing"] = fmt.Sprintf("%d", 2)
				}

			}
			result = append(result, result_tmp)
		}
	}
	return result
}
