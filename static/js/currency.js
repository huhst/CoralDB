// 显示时间到前端
var t = null;
t = setTimeout(time, 10);//開始运行
function time() {
    clearTimeout(t);//清除定时器
    dt = new Date();
    var y = dt.getFullYear();
    var mt = dt.getMonth() + 1;
    var day = dt.getDate();
    var h = dt.getHours();//获取时
    var m = dt.getMinutes();//获取分
    var s = dt.getSeconds();//获取秒
    var week = dt.getDate();// 获取星期
    var str = "";
    if (week === 0) {
        str = "星期日";
    } else if (week === 1) {
        str = "星期一";
    } else if (week === 2) {
        str = "星期二";
    } else if (week === 3) {
        str = "星期三";
    } else if (week === 4) {
        str = "星期四";
    } else if (week === 5) {
        str = "星期五";
    } else if (week === 6) {
        str = "星期六";
    }
    var table = document.getElementById('Table');
    var num = table.rows.length -1;
    var rowNum = document.getElementById('rowsNum');
    if (rowNum !== null){
        rowNum.innerHTML = num;
    }
    // 为类名为`.showTime`的元素添加时间
    document.querySelector(".time").innerHTML = + y + "/" + mt + "/" + day + "  " + h + ":" + m + ":" + s + "  "+str;
    t = setTimeout(time, 1000); //设定定时器，循环运行
}

/*点击关闭按钮*/
function closeBox() {
    var popBox = document.getElementById("popBox");
    var popLayer = document.getElementById("popLayer");
    popBox.style.display = "none";
    popLayer.style.display = "none";
}