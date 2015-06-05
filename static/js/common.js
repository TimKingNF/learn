function getUrlParam(name){
	var reg = new RegExp("(^|&)"+ name +"=([^&]*)(&|$)"); //构造一个含有目标参数的正则表达式对象
	var r = window.location.search.substr(1).match(reg);  //匹配目标参数
	if (r!=null) return unescape(r[2]); 
	return null; //返回参数值
} 

function switchWeekDay(name){
	switch(name) {
		case "1":
			return "一";
			break;
		case "2":
			return "二";
			break;
		case "3":
			return "三";
			break;
		case "4":
			return "四";
			break;
		case "5":
			return "五";
			break;
		case "6":
			return "六";
			break;
		case "7":
			return "日";
			break;
	}
	return ""
}