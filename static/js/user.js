// 用户退出
function exit(){
    var power = document.getElementById("power").innerHTML
    if (power === "0"){
        alert("服务器未连接，无法退出")
        return;
    }else if(power === "1"){
        alert("用户未登录，无法退出")
        return;
    }else if(power === "2") {
        var json = {}
        json.key = "exit"
        var a = JSON.stringify(json)
        $.ajax({
            url:"/user",
            type:"POST",
            data:a,
            success:function (data){
                if (data["info"]!=="退出成功"){
                    alert(data["info"])
                    return
                }else{
                    alert(data["info"])
                    document.getElementById("power").innerHTML = "1";
                    location.reload();
                }
            },
            fail:function (a){
                console.log(a)
            }
        })
    }
}

// 断开连接
function closeServer(){
    var power = document.getElementById("power").innerHTML
    if (power === "0"){
        alert("服务器已断开，无需再次断开")
        return;
    }else if(power === "1"){
        var json = {}
        json.key = "close"
        var a = JSON.stringify(json)
        $.ajax({
            url:"/user",
            type:"POST",
            data:a,
            success:function (data){
                if (data["info"]!=="断开连接成功"){
                    alert(data["info"])
                    return
                }else{
                    alert(data["info"])
                    document.getElementById("power").innerHTML = "0";
                    location.reload();
                }

            },
            fail:function (a){
                console.log(a)
            }
        })

    }else if(power === "2") {
        alert("用户已登录，无法断开连接")
        return;
    }
}

// 创建共享表
function creatTable(){
    var power = document.getElementById("power").innerHTML
    if (power === "0"){
        alert("未连接服务器，请先连接服务器")
        return;
    }else if(power === "1"){
        alert("用户未登录，请先登录")
        return;
    }else if(power === "2")  {
        document.getElementById('in1_for').style.display = "inline-table";
        document.getElementById('in2_for').style.display = "inline-table";
        document.getElementById('in3_for').style.display = "none";
        document.getElementById('in1').style.display = "inline-table";
        document.getElementById('in2').style.display = "inline-table";
        document.getElementById('in3').style.display = "none";

        document.getElementById('in1_for').innerText = "表     名:"
        document.getElementById('in2_for').innerText = "用户及权限:"

        var popBox = document.getElementById("popBox");
        var popLayer = document.getElementById("popLayer");
        popBox.style.display = "block";
        popLayer.style.display = "block";

        // 显示提交1按钮
        document.getElementById("up").style.display = "block";
        document.getElementById("up1").style.display = "none";
        document.getElementById("up2").style.display = "none";
        document.getElementById("up3").style.display = "none";
        document.getElementById("up4").style.display = "none";
        document.getElementById("up5").style.display = "none";
    }

}

// 添加数据
function addDates(){
    var power = document.getElementById("power").innerHTML
    if (power === "0"){
        alert("未连接服务器，请先连接服务器")
        return;
    }else if(power === "1"){
        alert("用户未登录，请先登录")
        return;
    }else if(power === "2") {
        document.getElementById('in1_for').style.display = "inline-table";
        document.getElementById('in2_for').style.display = "inline-table";
        document.getElementById('in3_for').style.display = "inline-table";
        document.getElementById('in1').style.display = "inline-table";
        document.getElementById('in2').style.display = "inline-table";
        document.getElementById('in3').style.display = "inline-table";

        document.getElementById('in1_for').innerText = "key:"
        document.getElementById('in2_for').innerText = "value:"
        document.getElementById('in3_for').innerText = "表名:"

        var popBox = document.getElementById("popBox");
        var popLayer = document.getElementById("popLayer");
        popBox.style.display = "block";
        popLayer.style.display = "block";

        // 显示提交1按钮
        document.getElementById("up").style.display = "none";
        document.getElementById("up1").style.display = "block";
        document.getElementById("up2").style.display = "none";
        document.getElementById("up3").style.display = "none";
        document.getElementById("up4").style.display = "none";
        document.getElementById("up5").style.display = "none";
    }

}

// 修改数据
function modifyDates(){
    var power = document.getElementById("power").innerHTML
    if (power === "0"){
        alert("未连接服务器，请先连接服务器")
        return;
    }else if(power === "1"){
        alert("用户未登录，请先登录")
        return;
    }else if(power === "2") {
        document.getElementById('in1_for').style.display = "inline-table";
        document.getElementById('in2_for').style.display = "inline-table";
        document.getElementById('in3_for').style.display = "inline-table";
        document.getElementById('in1').style.display = "inline-table";
        document.getElementById('in2').style.display = "inline-table";
        document.getElementById('in3').style.display = "inline-table";

        document.getElementById('in1_for').innerText = "key:"
        document.getElementById('in2_for').innerText = "value:"
        document.getElementById('in3_for').innerText = "表名:"

        var popBox = document.getElementById("popBox");
        var popLayer = document.getElementById("popLayer");
        popBox.style.display = "block";
        popLayer.style.display = "block";

        // 显示提交1按钮
        document.getElementById("up").style.display = "none";
        document.getElementById("up1").style.display = "none";
        document.getElementById("up2").style.display = "block";
        document.getElementById("up3").style.display = "none";
        document.getElementById("up4").style.display = "none";
        document.getElementById("up5").style.display = "none";
    }
}

// 登录
function login(){
    var power = document.getElementById("power").innerHTML
    if (power === "0"){
        alert("未连接服务器，请先连接服务器")
        return;
    }else if(power === "1"){
        document.getElementById('in1_for').style.display = "inline-table";
        document.getElementById('in2_for').style.display = "inline-table";
        document.getElementById('in3_for').style.display = "none";
        document.getElementById('in1').style.display = "inline-table";
        document.getElementById('in2').style.display = "inline-table";
        document.getElementById('in3').style.display = "none";

        document.getElementById('in1_for').innerText = "用户名:"
        document.getElementById('in2_for').innerText = "密  码:"

        var popBox = document.getElementById("popBox");
        var popLayer = document.getElementById("popLayer");
        popBox.style.display = "block";
        popLayer.style.display = "block";

        // 显示提交1按钮
        document.getElementById("up").style.display = "none";
        document.getElementById("up1").style.display = "none";
        document.getElementById("up2").style.display = "none";
        document.getElementById("up3").style.display = "block";
        document.getElementById("up4").style.display = "none";
        document.getElementById("up5").style.display = "none";
    }else if(power === "2") {
        alert("用户已登录，请退出后再试")
        return;
    }
}

// 连接服务器
function connectServer(){
    var power = document.getElementById("power").innerHTML
    if (power === "0"){
        document.getElementById('in1_for').style.display = "inline-table";
        document.getElementById('in2_for').style.display = "inline-table";
        document.getElementById('in3_for').style.display = "none";
        document.getElementById('in1').style.display = "inline-table";
        document.getElementById('in2').style.display = "inline-table";
        document.getElementById('in3').style.display = "none";

        document.getElementById('in1_for').innerText = "服务器ip:"
        document.getElementById('in2_for').innerText = "密   码:"

        var popBox = document.getElementById("popBox");
        var popLayer = document.getElementById("popLayer");
        popBox.style.display = "block";
        popLayer.style.display = "block";

        // 显示提交1按钮
        document.getElementById("up").style.display = "none";
        document.getElementById("up1").style.display = "none";
        document.getElementById("up2").style.display = "none";
        document.getElementById("up3").style.display = "none";
        document.getElementById("up4").style.display = "block";
        document.getElementById("up5").style.display = "none";
    }else if(power === "1"){
        alert("服务器已连接，无需再次连接")
        return;
    }else if(power === "2") {
        alert("服务器已连接，无需再次连接")
        return;
    }
}

// 修改共享表
function modifyTable(){
    var power = document.getElementById("power").innerHTML
    if (power === "0"){
        alert("未连接服务器，请先连接服务器")
        return;
    }else if(power === "1"){
        alert("用户未登录，请先登录")
        return;
    }else if(power === "2")  {
        document.getElementById('in1_for').style.display = "inline-table";
        document.getElementById('in2_for').style.display = "inline-table";
        document.getElementById('in3_for').style.display = "none";
        document.getElementById('in1').style.display = "inline-table";
        document.getElementById('in2').style.display = "inline-table";
        document.getElementById('in3').style.display = "none";

        document.getElementById('in1_for').innerText = "表     名:"
        document.getElementById('in2_for').innerText = "用户及权限:"

        var popBox = document.getElementById("popBox");
        var popLayer = document.getElementById("popLayer");
        popBox.style.display = "block";
        popLayer.style.display = "block";

        // 显示提交1按钮
        document.getElementById("up").style.display = "none";
        document.getElementById("up1").style.display = "none";
        document.getElementById("up2").style.display = "none";
        document.getElementById("up3").style.display = "none";
        document.getElementById("up4").style.display = "none";
        document.getElementById("up5").style.display = "block";
    }
}

// 创建共享表
function user_save(){
    var tablename = document.getElementById('in1').value
    var userpower = document.getElementById('in2').value
    if (tablename === ""){
        alert("表名不能为空！");
        return;
    }
    if (userpower === ""){
        alert("权限不能为空！");
        return;
    }
    var json = {}
    json.key = "creatTable"
    json.tablename = tablename
    json.userpower = userpower
    var a = JSON.stringify(json)
    $.ajax({
        url:"/user",
        type:"POST",
        data:a,
        success:function (data){
            if (data["info"]!=="创建成功"){
                alert(data["info"])
                return
            }else{
                document.getElementById("power").innerHTML = "2";
                alert(data["info"])
                var popBox = document.getElementById("popBox");
                var popLayer = document.getElementById("popLayer");
                popBox.style.display = "none";
                popLayer.style.display = "none";
                // 将表单置为空
                document.getElementById("in1").value = "";
                document.getElementById("in2").value = "";
                location.reload();
                return
            }
        },
        fail:function (a){
            console.log(a)
        }
    })
}

// 添加数据
function user_save1(){
    var datakey = document.getElementById('in1').value
    var value = document.getElementById('in2').value
    var tablename = document.getElementById('in3').value
    var json = {}
    json.key = "addDates"
    json.datakey = datakey
    json.value = value
    json.tablename = tablename
    var a = JSON.stringify(json)
    $.ajax({
        url:"/user",
        type:"POST",
        data:a,
        success:function (data){
            if (data["info"]!=="插入数据成功"){
                alert(data["info"])
                return
            }else{
                document.getElementById("power").innerHTML = "2";
                alert(data["info"])
                var popBox = document.getElementById("popBox");
                var popLayer = document.getElementById("popLayer");
                popBox.style.display = "none";
                popLayer.style.display = "none";
                // 将表单置为空
                document.getElementById("in1").value = "";
                document.getElementById("in2").value = "";
                document.getElementById("in3").value = "";
                location.reload();
                return
            }
        },
        fail:function (a){
            console.log(a)
        }
    })
}

// 修改数据
function user_save2(){
    var datakey = document.getElementById('in1').value
    var value = document.getElementById('in2').value
    var tablename = document.getElementById('in3').value
    // console.log(tablename)
    var json = {}
    json.key = "modifyDates"
    json.datakey = datakey
    json.value = value
    json.tablename = tablename
    var a = JSON.stringify(json)
    $.ajax({
        url:"/user",
        type:"POST",
        data:a,
        success:function (data){
            if (data["info"]!=="修改数据成功"){
                alert(data["info"])
                return
            }else{
                document.getElementById("power").innerHTML = "2";
                alert(data["info"])
                var popBox = document.getElementById("popBox");
                var popLayer = document.getElementById("popLayer");
                popBox.style.display = "none";
                popLayer.style.display = "none";
                // 将表单置为空
                document.getElementById("in1").value = "";
                document.getElementById("in2").value = "";
                document.getElementById("in3").value = "";
                location.reload();
                return
            }
        },
        fail:function (a){
            console.log(a)
        }
    })
}

// 用户登录
function user_save3(){
    var username = document.getElementById('in1').value
    var password = document.getElementById('in2').value
    var json = {}
    json.key = "login"
    json.username = username
    json.password = password
    var a = JSON.stringify(json)
    $.ajax({
        url:"/user",
        type:"POST",
        data:a,
        success:function (data){
            if (data["info"]!=="登录成功"){
                alert(data["info"])
                return
            }else{
                document.getElementById("power").innerHTML = "2";
                alert(data["info"])
                var popBox = document.getElementById("popBox");
                var popLayer = document.getElementById("popLayer");
                popBox.style.display = "none";
                popLayer.style.display = "none";
                // 将表单置为空
                document.getElementById("in1").value = "";
                document.getElementById("in2").value = "";
                location.reload();
                return
            }
        },
        fail:function (a){
            console.log(a)
        }
    })
}

// 连接服务器
function user_save4(){
    var ip = document.getElementById('in1').value
    var password = document.getElementById('in2').value
    if (ip === ""){
        alert("节点IP不能为空！");
        return;
    }
    // 判断ip的格式
    if (!solve(ip)){
        alert("节点IP格式不正确！")
        return;
    }
    if (password ===""){
        alert("密码不能为空！");
        return;
    }
    var json = {}
    json.key = "connectServer"
    json.ip = ip
    json.password = password
    var a = JSON.stringify(json)
    $.ajax({
        url:"/user",
        type:"POST",
        data:a,
        success:function (data){
            if (data["info"]!=="连接成功"){
                alert(data["info"])
                return
            }else{
                document.getElementById("power").innerHTML = "1";
                alert(data["info"])
                var popBox = document.getElementById("popBox");
                var popLayer = document.getElementById("popLayer");
                popBox.style.display = "none";
                popLayer.style.display = "none";
                // 将表单置为空
                document.getElementById("in1").value = "";
                document.getElementById("in2").value = "";
                location.reload();
                return
            }
        },
        fail:function (a){
            console.log(a)
        }
    })
}

function solve( IP ) {
    return checkIPv4(IP)?true:(checkIPv6(IP));
}
function checkIPv4(IP){
    let arr=IP.split(".");
    if(arr.length!==4) return false;
    for(let i of arr){
        if(Object.is(Number(i),NaN)||Number(i)>255||Number(i)<0||i[0]==='0'&&i.length!==1){
            return false;
        }
    }
    return true;
}
function checkIPv6(IP){
    let arr=IP.split(":");
    if(arr.length!==8) return false;
    for (let i of arr){
        if(i.length>4||Object.is(parseInt(i,16),NaN)){
            return false;
        }
    }
    return true;
}

// 修改共享表
function user_save5(){
    var tablename = document.getElementById('in1').value
    var userpower = document.getElementById('in2').value
    if (tablename === ""){
        alert("表名不能为空！");
        return;
    }
    if (userpower === ""){
        alert("权限不能为空！");
        return;
    }
    var json = {}
    json.key = "modifyTable"
    json.tablename = tablename
    json.userpower = userpower
    var a = JSON.stringify(json)
    $.ajax({
        url:"/user",
        type:"POST",
        data:a,
        success:function (data){
            if (data["info"]!=="修改成功"){
                alert(data["info"])
                return
            }else{
                document.getElementById("power").innerHTML = "2";
                alert(data["info"])
                var popBox = document.getElementById("popBox");
                var popLayer = document.getElementById("popLayer");
                popBox.style.display = "none";
                popLayer.style.display = "none";
                // 将表单置为空
                document.getElementById("in1").value = "";
                document.getElementById("in2").value = "";
                location.reload();
                return
            }
        },
        fail:function (a){
            console.log(a)
        }
    })
}