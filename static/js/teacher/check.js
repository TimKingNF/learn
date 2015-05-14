$( function () {

	// 点名按钮事件
	$( '#check' ).click( function () {

		var num = $( 'input[ name=num ]' );
		var maxnum = $( '.maxnum' ).text();
		var reg = /^[1-9]+[0-9]*]*$/;

		if( reg.test( num.val() ) && Number( num.val() ) <= Number( maxnum ) ) {

			// 点名列表生成
			$( '.checkarea' ).remove();
			var obj = '<table class="am-table am-table-bordered am-table-radius am-table-striped info">'+
				'<thead><tr><th>姓名</th><th>学号</th><th>班级</th><th></th></tr></thead>'+
				'<tbody>';
			var res = createRandomNumber( maxnum , num.val() );

			for ( var i = 0; i < data.length; i++ ) {

				// 判断是否在生成的名单里
				if ( in_array( res , i ) ) {
					obj += '<tr>'+
					'<td>' + data[i].name + '</td>' +
					'<td>' + data[i].id + '<input type="hidden" name="id[' + i + ']" value="' + data[i].id + '" /></td>' +
					'<td>' + data[i].class + '</td>' +
					'<td style="padding:5px;" width="30%">' +
					'<div class="am-btn-group am-btn-group-justify" data-am-button>' +
					'<label class="am-btn am-btn-success am-btn-sm am-active"><input type="radio" name="type[' + i + ']" value="0" checked>已到</label>' +
					'<label class="am-btn am-btn-warning am-btn-sm"><input type="radio" name="type[' + i + ']" value="1">迟到</label>' +
					'<label class="am-btn am-btn-danger am-btn-sm"><input type="radio" name="type[' + i + ']" value="2">未到</label>' +
					'<label class="am-btn am-btn-secondary am-btn-sm"><input type="radio" name="type[' + i + ']" value="3">请假</label>' +
					'</div>' +
					'</td>' +
					'</tr>';
				}
				else {
					obj += '<tr style="display: none;">' +
					'<td></td>' +
					'<td><input type="hidden" name="id[' + i + ']" value="' + data[i].id + '" /></td></td>' +
					'<td></td>' +
					'<td>' +
					'<input type="radio" style="display:none;" name="type[' + i + ']" checked value="0">' +
					'</td>' +
					'</tr>';
				}
			}

			obj += '</tbody></table><button type="submit" class="am-btn am-btn-primary am-btn-lg am-radius">确定提交</button>';
			$('.content').append(obj);


		}
		else {
			num.val('');
		}
	});

	// 生成不重复随机数
	// @param 取值范围 int num
	// @param 取数长度 int maxNum
	// return array
	function createRandomNumber ( num , maxNum ) {

		var flag = 0,
			i=0,
			arr=[];

		if ( maxNum - num < 0 ) {
			flag = maxNum;
			maxNum = num;
			num = flag;
		}

		for ( ; i < maxNum; i++ ) {
			arr[i] = i;
		}

		arr.sort( function ( a , b ) {
			return 0.5 - Math.random();
		});

		arr.length = num;

		// 从小到大排序
		arr.sort( function ( a , b ) {
			return a > b ? true : false;
		});

		return arr;
	}

	// 判断该值是否在数组内
	// @param 被遍历数组 array arr
	// @param 被查找值 object(除数组,json) find
	// return boolean
	function in_array ( arr , find ) {

		for ( var i = 0; i < arr.length; i++ ) {

			if ( arr[i] == find ) {
				return true;
			}

		}

		return false;
	}

});