$(function() {
    var modal = $('#modal');
    $('.tb-class').each(function() {
        $(this).click(function() {
            var name = $(this).children('span').text();
            $('#course-name').text(name);
            modal.modal();
        });
    });
});