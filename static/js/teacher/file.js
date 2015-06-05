$( function () {

    var tools = $( '.tools' );
    var is_show = 0; // 用于判断是否显示 tools
    var type;
    $( '#all' ).click( function () {
        type = $( this ).is( ':checked' );

        $( 'tbody > tr' ).each( function () {
            is_check ( $( this ) , tools , type );
        });

    });

    // 文件事件
    $( 'tbody > tr' ).click( function () {
        is_check ( $( this ) , tools );
    });

    // 双击文件夹事件
    $( 'tr.folder' ).dblclick( function () {
        alert( 'double click folder' );
    });

    // 删除文件事件
    $( '#delete' ).click( function () {
        $( '#confirm' ).modal({
            relatedTarget: this,
            onConfirm: function ( options ) {
                alert('confirm');
            }
        });
    });

    // 新建文件夹事件
    $( '#add' ).click( function () {
        $( '#prompt' ).modal({
            relatedTarget: this,
            onConfirm: function( e ) {
                alert('confirm');
                console.log('输入的内容为：'+ e.data);
            }
        });
    });

    // 检测是否被选中
    function is_check ( obj , tools , type ) {

        var oCb = obj.children( '.tb-cx' ).children( '.ipt-cx' );
        if( type || !obj.hasClass( 'active' ) ) {
            oCb.attr( 'checked' , true );
            obj.addClass( 'active' );
            is_show++;
            tools.addClass( 'am-show' );
            tools.removeClass( 'am-hide' );
        }
        else {
            oCb.removeAttr( 'checked' );
            obj.removeClass( 'active' );
            is_show--;
            if ( is_show <= 0 ) {
                tools.removeClass( 'am-show' );
                tools.addClass( 'am-hide' );
            }
        }

    }
});