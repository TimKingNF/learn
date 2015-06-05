$( function () {
    upload($(".upload-btn"));
    updateUpload($(".update-upload-btn"));

    echarts.init( document.getElementById( 'line') ).setOption( optionLine );
    $( window).resize( function() {
        echarts.init( document.getElementById( 'line' ) ).setOption( optionLine );
    });

    $(".select-form").find("select[name='xnd']").change(function() {
        $.post("/api/teacher/getTermNumberByTeacher", {
            "key": key,
            "signature": GetSignature(getCookie("token"), appid, sessid),
            "xnd": $(this).children('option:selected').val(),
        }, function(data){
            if (data.Status != 0) {
                alert(data.Data);
                return
            } else {
                var objs = data.Data, html = '<option value="0">-未选择-</option>';
                if (objs) {
                    for (var i = 0; i < objs.length; i++) {
                        html += '<option value="'+objs[i].Number+'">'+objs[i].Number+'</option>';
                    }
                    $(".select-form").find("select[name='xqd']").html(html);
                    $(".select-form").find("select[name='course']").html('<option value="0">-未选择-</option>');
                } else {
                    alert("操作失败");
                    return
                }
            }
        }, "json");
    });

    $(".select-form").find("select[name='xqd']").change(function() {
        $.post("/api/teacher/getTeacherCourseByTerm", {
            "key": key,
            "signature": GetSignature(getCookie("token"), appid, sessid),
            "xnd": $(".select-form").find("select[name='xnd']").children('option:selected').val(),
            "xqd": $(this).children('option:selected').val(),
        }, function(data){
            if (data.Status != 0) {
                alert(data.Data);
                return
            } else {
                var objs = data.Data, html = "";
                if (objs) {
                    for (var i = 0; i < objs.length; i++) {
                        var time = JSON.parse(objs[i].Time);
                        html += '<option value="'+objs[i].Id+'">'+objs[i].Course.Name+' (周'+switchWeekDay(time.week_day)+' 第'+time.week_time+'节)</option>';
                    }
                    $(".select-form").find("select[name='course']").html(html);
                } else {
                    alert("操作失败");
                    return
                }
            }
        }, "json");
    });

    $('#addhomework').click(function() {
        $('#prompt').modal({
            relatedTarget: this,
            onConfirm: function(e) {
                var data = e.data;
                //  add homework
                $.post("/api/teacher/publisHomework", {
                    "key": key,
                    "signature": GetSignature(getCookie("token"), appid, sessid),
                    "title": data[0],
                    "remark": data[1],
                    "publish_week": data[2],
                    "as_of_time": data[3],
                    "attachment": data[4],
                    "t_course": data[5],
                }, function(data) {
                    if (data.Status != 0) {
                        alert(data.Data);
                    } else {
                        alert("发布成功");
                    }
                }, "json");
            }
        });
    });

    $('.update-btn').click(function() {
        $.post("/api/teacher/getTeacherCourseHomework", {
            "key": key,
            "signature": GetSignature(getCookie("token"), appid, sessid),
            "t_course_homework": $(this).attr("data-teacher-course-homework-id"),
        }, function(data) {
            if (data.Status != 0) {
                alert("获取作业失败");
                return
            } else {
                var html = "", obj = data.Data, btn = $('#update');
                if (obj) {
                    btn.find("#update-homework-tilte").val(obj.Title);
                    btn.find("#update-homework-remark").html(obj.Remark);
                    btn.find("#update-homework-publishweek").val(obj.PublishWeek);
                    btn.find("#update-homework-as_of_time").val(obj.AsOfTime);
                    btn.find(".update-homework-id").val(obj.Id);
                    if (!obj.Attachment)  btn.find(".update-attachment-file").css({"display": "none"});
                    else {
                        btn.find(".update-attachment").val(obj.Attachment);
                        btn.find(".update-attachment-file>a").html(obj.Attachment.substring(20, obj.Attachment.length));
                        btn.find(".update-attachment-file>a").attr("href", "/view/teacher/download" + "?filepath=" + obj.Attachment);
                        btn.find(".update-attachment-file").css({"display": "block"});
                    }
                }
                btn.modal({
                    relatedTarget: this,
                    onConfirm: function(e) {
                        var data = e.data;
                        //  update homework
                        $.post("/api/teacher/updateHomework", {
                            "key": key,
                            "signature": GetSignature(getCookie("token"), appid, sessid),
                            "teacher": teacher,
                            "title": data[0],
                            "remark": data[1],
                            "publish_week": data[2],
                            "as_of_time": data[3],
                            "attachment": data[4],
                            "t_course": data[5],
                            "t_course_homework": data[6],
                        }, function(data) {
                            if (data.Status != 0) {
                                alert(data.Data);
                            } else {
                                alert("更新成功");
                            }
                        }, "json");
                    }
                });
            }
        }, "json");
    });

    $(".delete-btn").click(function() {
        if (!window.confirm("确定要删除该作业么")) {
            return false;
        }
        var t_course_homework_id = $(this).attr("data-teacher-course-homework-id");
        $.post("/api/teacher/delTeacherCourseHomework", {
            "key": key,
            "signature": GetSignature(getCookie("token"), appid, sessid),
            "t_course_homework": t_course_homework_id,
        }, function(data) {
            if (data.Status != 0) {
                alert(data.Data);
            } else {
                $("tr[data-teacher-course-homework-id='"+t_course_homework_id+"']").remove();
            }
        }, "json");
    });

    $('.score').click(function() {
        $.post("/api/teacher/getStudentHomeworksByTeacherCourseHomework", {
            "key": key,
            "signature": GetSignature(getCookie("token"), appid, sessid),
            "t_course_homework": $(this).attr("data-teacher-course-homework-id"),
        }, function(data) {
            if (data.Status != 0) {
                alert("获取学生上交作业记录失败");
                return
            } else {
                var html = "", objs = data.Data;
                if (objs) {
                    for (var i = 0; i < objs.length; i++) {
                        html += "<tr><td>"+objs[i].Student.Id+"</td><td>"+objs[i].Student.Name+"</td><td>"+objs[i].Student.Class.Name+"</td>"
                        html += '<td><a href="/view/teacher/download?filepath='+objs[i].Attachment+'" class="am-btn am-btn-secondary am-btn-sm all-download-btn"><span class="am-icon-download"></span></a></td>'
                        html += "<td><select class='setGrade' data-student-homework-id='"+objs[i].Id+"'>"
                        html += "<option>-未选择-</option>"
                        html += "<option value='A' "+ (objs[i].Grade == "A" ? "selected" : "") +">A</option>"
                        html += "<option value='B' "+ (objs[i].Grade == "B" ? "selected" : "") +">B</option>"
                        html += "<option value='C' "+ (objs[i].Grade == "C" ? "selected" : "") +">C</option>"
                        html += "<option value='D' "+ (objs[i].Grade == "D" ? "selected" : "") +">D</option>"
                        html += "<option value='E' "+ (objs[i].Grade == "E" ? "selected" : "") +">E</option></select></td></tr>"
                    }
                }
                $("#popup tbody").html(html);
                $(".setGrade").change(function() {
                    $.post("/api/teacher/setGradeByStudentHomework", {
                        "key": key,
                        "signature": GetSignature(getCookie("token"), appid, sessid),
                        "grade" : $(this).children('option:selected').val(),
                        "s_homework": $(this).attr("data-student-homework-id"),
                    }, function(data) {
                        if (data.Status != 0) {
                            alert(data.Data);return
                        } 
                    }, "json");
                });
                $('#popup').modal();
            }
        }, "json")
    });

    var modal = $('#modal');
    $('.modal-btn').click(function() {
        var remark = $(this).attr("data-remark"), attachment = $(this).attr("data-attachment");
        modal.find(".am-padding-vertical-lg").html(remark || "暂无说明");
        if (!attachment) modal.find("a").css({"display":"none"});
        else {
            modal.find("a").attr("href", "/view/teacher/download" + "?filepath=" + attachment);
            modal.find("a").css({"display":"block"});
        }
        modal.modal();
    });
});