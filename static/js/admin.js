// 添加用户
var Root = "http://localhost"
function adduser() {
    // 元素的可见
    document.getElementById('in1_for').style.display = "inline-table";
    document.getElementById('in2_for').style.display = "inline-table";
    document.getElementById('in3_for').style.display = "none";
    document.getElementById('in1').style.display = "inline-table";
    document.getElementById('in2').style.display = "inline-table";
    document.getElementById('in3').style.display = "none";
    // 元素赋值
    document.getElementById('in1_for').innerText = "用户名:"
    document.getElementById('in2_for').innerText = "密 码:"
    var popBox = document.getElementById("popBox");
    var popLayer = document.getElementById("popLayer");
    popBox.style.display = "block";
    popLayer.style.display = "block";

    // 显示提交1按钮
    document.getElementById("up").style.display = "block";
    document.getElementById("up1").style.display = "none";
    document.getElementById("up2").style.display = "none";
    document.getElementById("up3").style.display = "none";
}

// 创建集群
function creatCluster(){
    // 元素的可见
    document.getElementById('in1_for').style.display = "inline-table";
    document.getElementById('in2_for').style.display = "inline-table";
    document.getElementById('in3_for').style.display = "none";
    document.getElementById('in1').style.display = "inline-table";
    document.getElementById('in2').style.display = "inline-table";
    document.getElementById('in3').style.display = "none";

    document.getElementById('in1_for').innerText = "集群名:"
    document.getElementById('in2_for').innerText = "密  码:"
    var popBox = document.getElementById("popBox");
    var popLayer = document.getElementById("popLayer");
    popBox.style.display = "block";
    popLayer.style.display = "block";

    // 显示提交1按钮
    document.getElementById("up").style.display = "none";
    document.getElementById("up1").style.display = "block";
    document.getElementById("up2").style.display = "none";
    document.getElementById("up3").style.display = "none";
    // 将表单置为空

}
// 加入集群
function joinCluster(){
    document.getElementById('in1_for').style.display = "inline-table";
    document.getElementById('in2_for').style.display = "inline-table";
    document.getElementById('in3_for').style.display = "none";
    document.getElementById('in1').style.display = "inline-table";
    document.getElementById('in2').style.display = "inline-table";
    document.getElementById('in3').style.display = "none";

    document.getElementById('in1_for').innerText = "IP:"
    document.getElementById('in2_for').innerText = "密码:"
    // document.getElementById('in3_for').innerText = "密码:"

    var popBox = document.getElementById("popBox");
    var popLayer = document.getElementById("popLayer");
    popBox.style.display = "block";
    popLayer.style.display = "block";

    // 显示提交1按钮
    document.getElementById("up").style.display = "none";
    document.getElementById("up1").style.display = "none";
    document.getElementById("up2").style.display = "block";
    document.getElementById("up3").style.display = "none";
    // 将表单置为空
}

//设置打包数
function setNum(){
    document.getElementById('in1_for').style.display = "inline-table";
    document.getElementById('in2_for').style.display = "none";
    document.getElementById('in3_for').style.display = "none";
    document.getElementById('in1').style.display = "inline-table";
    document.getElementById('in2').style.display = "none";
    document.getElementById('in3').style.display = "none";

    document.getElementById('in1_for').innerText = "打包数："
    // document.getElementById('in3_for').innerText = "密码:"

    var popBox = document.getElementById("popBox");
    var popLayer = document.getElementById("popLayer");
    popBox.style.display = "block";
    popLayer.style.display = "block";

    // 显示提交1按钮
    document.getElementById("up").style.display = "none";
    document.getElementById("up1").style.display = "none";
    document.getElementById("up2").style.display = "none";
    document.getElementById("up3").style.display = "block";
}
function index_save(){
    // 判断输入框是否有值
    var username = document.getElementById('in1').value
    var password = document.getElementById('in2').value
    if (username === ""){
        alert("用户名不能为空！");
        return;
    }
    if (password ===""){
        alert("密码不能为空！");
        return;
    }
    var json = {}
    json.key = "adduser"
    json.name = username
    json.password = password
    var a = JSON.stringify(json)
    $.ajax({
        url:"/admin",
        type:"POST",
        data:a,
        success:function (data){
            if (data["info"]!=="添加成功"){
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
            }else{
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
function index_save1(){
    var username = document.getElementById('in1').value
    var password = document.getElementById('in2').value
    if (username === ""){
        alert("集群名不能为空！");
        return;
    }
    if (password ===""){
        alert("密码不能为空！");
        return;
    }
    var json = {}
    json.key = "creatCluster"
    json.password = password
    var a = JSON.stringify(json)
    $.ajax({
        url:"/admin",
        type:"POST",
        data:a,
        success:function (data){
            if (data["info"]!=="集群创建成功"){
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
            }else{
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
function index_save2(){
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
    json.key = "joinCluster"
    json.ip = ip
    json.password = password
    var a = JSON.stringify(json)
    $.ajax({
        url:"/admin",
        type:"POST",
        data:a,
        success:function (data){
            if (data["info"]!=="加入集群成功"){
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
            }else{
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

function index_save3(){
    var num = document.getElementById('in1').value
    if (num ===""){
        alert("打包数不能为空！")
        return;
    }
    if (num <= 1){
        alert("打包数最小为1！")
        return;
    }
    var json = {}
    json.key = "set"
    json.num = num
    var a = JSON.stringify(json)
    $.ajax({
        url:"/admin",
        type:"POST",
        data:a,
        success:function (data){
            if (data["info"]!=="设置成功"){
                alert(data["info"])
                var popBox = document.getElementById("popBox");
                var popLayer = document.getElementById("popLayer");
                popBox.style.display = "none";
                popLayer.style.display = "none";
                // 将表单置为空
                document.getElementById("in1").value = "";
                location.reload();
                return
            }else{
                alert(data["info"])
                var popBox = document.getElementById("popBox");
                var popLayer = document.getElementById("popLayer");
                popBox.style.display = "none";
                popLayer.style.display = "none";
                // 将表单置为空
                document.getElementById("in1").value = "";
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