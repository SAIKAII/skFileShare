layui.use('element', function () {
    var element = layui.element;

    element.tabChange("im_body", layid);
})

var ws;
var server_host;
window.onload = function () {
    this.server_host = window.location.host;
    var divAll = document.getElementById("div-all");
    divAll.style.height = document.documentElement.clientHeight + "px";
    ws = new WebSocket("ws://localhost:8080/ws");
    ws.onmessage = function (evt) {
        var content = document.getElementById("dialog");
        var str = JSON.parse(evt.data);

        var div1 = document.createElement('div');
        div1.style.margin = '15px 0px';
        var blueDot = document.createElement('span');
        blueDot.className = 'layui-badge-dot layui-bg-blue';
        blueDot.style.margin = '0px 15px';
        var commRecord = document.createElement('button');
        commRecord.className = 'layui-btn layui-btn-radius layui-btn-primary';
        commRecord.innerText = str.ip + ':' + str.body;

        content.appendChild(div1);
        div1.appendChild(blueDot);
        div1.appendChild(commRecord);
        content.scrollTop = content.scrollHeight;
    }
    ws.onopen = function (evt) {
        console.log("Connect succeed!");
    }
    ws.onclose = function () {
        ws.close();
    }
};

function sendMsg(id) {
    var inputData = document.getElementById(id);
    ws.send(inputData.value);
    inputData.value = '';
}

function getFiles() {
    $('#files-list').empty();
    $.ajax({
        type: "GET",
        url: "http://" + server_host + "/getfiles",
        success: function(data) {
            data = JSON.parse(data);
            for (var p in data){
                addItem(data[p]);
            }
        }
    })
}

function addItem(entry) {
    var item = $("<div></div>").addClass('item');
    $('#files-list').append(item);
    var hyper = $("<a></a>");
    hyper.attr('href', 'http://' + server_host + '/getfile/' + entry.name).attr('download', entry.name);
    hyper.attr('target', '_blank');
    hyper.appendTo(item);
    var file_icon = $("<p></p>");
    file_icon.appendTo(hyper);
    if(entry.ext == ".jpg"){
        file_icon.addClass('file_icon layui-icon layui-icon-picture');
    }else{
        file_icon.addClass('file_icon layui-icon layui-icon-app');
    }
    
    $("<p></p>").css('text-overflow', 'ellipsis').css('overflow', 'hidden').css('white-space', 'nowrap').text(entry.name).appendTo(item);
}