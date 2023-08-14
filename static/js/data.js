
function dataQuery(){
    var power = document.getElementById("power").innerHTML
    if (power === "0"){
        alert("服务器未连接，请先连接服务器")
        return;
    }else if (power === "1") {
        alert("用户未登录，请先登录")
        return;
    }else if (power === "2"){
        var datakey = document.getElementById('dataKey').value;
        var tablename = document.getElementById('tableName').value;
        if (datakey === ""){
            alert("数据的key不能为空！")
            return;
        }
        if (tablename === ""){
            alert("数据所在的表名不能为空！")
            return;
        }
        var json = {}
        json.key = "dataQuery"
        json.datakey = datakey
        json.tablename = tablename
        var a = JSON.stringify(json)
        $.ajax({
            url:"/data",
            type:"POST",
            data:a,
            success:function (data){
                if (data["info"]!=="查询成功") {
                    alert(data["info"])
                    return
                }else {
                    document.getElementById("power").innerHTML = "2";
                    // 获取值
                    var tableName = data["tablename"]
                    var Key = data["key"]
                    var Value = data["value"]
                    var User = data["user"]
                    var time = data["time"]
                    var getTime = data["gettime"]
                    // 添加表格
                    addTable(tableName,Key,Value,User,time,getTime)
                }
            },
            fail:function (a){
                console.log(a)
            }
        })
    }
}

function addTable(tableName,Key,Value,User,time,getTime){
    // 改变原本得到的数据格式为textNode格式
    tableName = document.createTextNode(tableName)
    Key = document.createTextNode(Key);
    Value = document.createTextNode(Value);
    User = document.createTextNode(User);
    time = document.createTextNode(time);
    getTime = document.createTextNode(getTime)
    // 获得table
    var table = document.getElementById('listTable');
    // 创建tr
    var tr = document.createElement("tr");
    // 创建 td 并赋于class和值（创建单元格，并输入值）
    // 编号
    var td1 = document.createElement("td");
    td1.className = "col1";
    td1.appendChild(tableName);
    // 昵称
    var td2 = document.createElement("td");
    td2.className = "col2";
    td2.appendChild(Key)
    // 手机号
    var td3 = document.createElement("td");
    td3.className = "col3";
    td3.appendChild(Value)
    // 邮箱
    var td4 = document.createElement("td");
    td4.className = "col4";
    td4.appendChild(User)
    // 管理的业务员
    var td5 = document.createElement("td");
    td5.className = "col5";
    td5.appendChild(time)
    var td6 = document.createElement("td");
    td6.className = "col6";
    td6.appendChild(getTime);

    // 将tr加入table中
    table.appendChild(tr);
    // 将td依次加入tr中
    tr.appendChild(td1);
    tr.appendChild(td2);
    tr.appendChild(td3);
    tr.appendChild(td4);
    tr.appendChild(td5);
    tr.appendChild(td6);
}