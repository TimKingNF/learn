$(function() {
    var modal = $('#modal');
    $('.tb-class').each(function() {
        $(this).click(function() {
            var name = $(this).children('span').text(), 
            t_course_id = $(this).attr("data-teacher-course-id"),
            now_week = getUrlParam("week");
            $('#course-name').text(name);
            //	add eventlistener for a
            modal.find("#course-info").click(function() {
            	window.location.href = "/view/teacher/courseInfo" + "?course=" + t_course_id;
            });
            modal.find("#teacher-course-homework").click(function() {
            	window.location.href = "/view/teacher/courseHomework" + "?course=" + t_course_id;
            });
            modal.find("#teacher-course-history").click(function() {
            	window.location.href = "/view/teacher/courseHistory" + "?course=" + t_course_id + (now_week ? "&week="+now_week : "");
            });
            modal.modal();
        });
    });
});