$(function() {
    var modal = $('#modal');
    $('.tb-class').each(function() {
        $(this).click(function() {
            var name = $(this).children('span').text(),
            s_course_id = $(this).attr("data-student-course-id"),
            now_week = getUrlParam("week");
            $('#course-name').text(name);
            //	add eventlistener for a
            modal.find("#course-info").click(function() {
            	window.location.href = "/view/student/courseInfo" + "?course=" + s_course_id;
            });
            modal.find("#student-course-homework").click(function() {
            	window.location.href = "/view/student/courseHomework" + "?course=" + s_course_id + (now_week ? "&week="+now_week : "");
            });
            modal.find("#student-course-history").click(function() {
            	window.location.href = "/view/student/courseHistory" + "?course=" + s_course_id + (now_week ? "&week="+now_week : "");
            });
            modal.modal();
        });
    });
});