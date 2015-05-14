$( function () {

    var tools = $( '.tools' );
    var is_show = 0; // �����ж��Ƿ���ʾ tools
    var type;
    $( '#all' ).click( function () {
        type = $( this ).is( ':checked' );

        $( 'tbody > tr' ).each( function () {
            is_check ( $( this ) , tools , type );
        });

    });

    // �ļ��¼�
    $( 'tbody > tr' ).click( function () {
        is_check ( $( this ) , tools );
    });

    // ˫���ļ����¼�
    $( 'tr.folder' ).dblclick( function () {
        alert( 'double click folder' );
    });

    // ɾ���ļ��¼�
    $( '#delete' ).click( function () {
        $( '#confirm' ).modal({
            relatedTarget: this,
            onConfirm: function ( options ) {
                alert('confirm');
            }
        });
    });

    // �½��ļ����¼�
    $( '#add' ).click( function () {
        $( '#prompt' ).modal({
            relatedTarget: this,
            onConfirm: function( e ) {
                alert('confirm');
                console.log('���������Ϊ��'+ e.data);
            }
        });
    });

    // ����Ƿ�ѡ��
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