$(function() {

    var contact = $('.contact');
    var content = $('.content');

    setHeight();
    $(window).resize(function() {
        setHeight();
    });
    function setHeight() {
        contact.height($(window).height() - content.offset().top);
        content.height($(window).height() - content.offset().top - 40);
    }
});