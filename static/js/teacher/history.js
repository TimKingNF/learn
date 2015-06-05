$( function () {

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
});