// 刷新信息
function tables(){
    var power = document.getElementById("power").innerHTML
    if (power === "0"){
        alert("服务器未连接，请先连接服务器")
        return;
    }else if (power === "1") {
        alert("用户未登录，请先登录")
        return;
    }else if (power === "2") {
        var json = {}
        json.key = "getMyTables"
        var a = JSON.stringify(json)
        $.ajax({
            url: "/table",
            type: "POST",
            data: a,
            success: function (data) {
                if (data["info"]!=="查询成功"){
                    alert(data["info"])
                    return
                }else{
                    $("#Table  tr:not(:first)").remove("");
                    // 创建表格
                    for(var i = 0;i<data["sum"];i++){
                        // console.log(data["infos"][i])
                        const myArray = data["infos"][i].split(",");
                        addTable(myArray[0],myArray[1],myArray[2]);
                    }
                }
            },
            fail: function (a) {
                console.log(a)
            }
        })
    }
}

function addTable(tableName,powerL,user) {
    var power = ""
    if (powerL === "1"){
        power = document.createTextNode("只读");
    }
    if (powerL === "2"){
        power = document.createTextNode("读写");
    }
    if (powerL === "3"){
        power = document.createTextNode("覆盖写");
    }
    if (powerL >= "4"){
        power = document.createTextNode("表修改");
    }
    // 改变原本得到的数据格式为textNode格式
    tableName = document.createTextNode(tableName);
    powerL = document.createTextNode(powerL);
    user = document.createTextNode(user);

    // 获得table
    var table = document.getElementById('listTable1');

    // 创建tr
    var tr = document.createElement("tr");
    // 创建 td 并赋于class和值（创建单元格，并输入值）
    // 编号
    var td1 = document.createElement("td");
    td1.className = "col1";
    td1.appendChild(tableName)
    // 昵称
    var td2 = document.createElement("td");
    td2.className = "col2";
    td2.appendChild(powerL)
    // 手机号
    var td3 = document.createElement("td");
    td3.className = "col3";
    td3.appendChild(power)
    // 邮箱
    var td4 = document.createElement("td");
    td4.className = "col4";
    td4.appendChild(user)

    // 将tr加入table中
    table.appendChild(tr);
    // 将td依次加入tr中
    tr.appendChild(td1);
    tr.appendChild(td2);
    tr.appendChild(td3);
    tr.appendChild(td4);
}

