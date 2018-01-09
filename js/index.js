$(function(){ 
	var name = getUrlParam('name');
	if (name != null && name !='') {
		$('#input_loginname').val(name);
	}
	$('#modal_login').modal('show');
}); 

var servaddr = window.location.host;
var wsaddr = 'ws://'+servaddr+'/login'
var ws;
var lasttime = '';

$('#button_login').click(function (){
	var name = $('#input_loginname').val();
	if (name == '' ) {
		alert('名字不能为空!');
		return false;
	} else if ( name.length > 10) {
		alert('名字不能超过10个字符');
		return false;
	}

	if (ws) {
		return false;
	}

	ws = new WebSocket(wsaddr+'?name='+name);
	ws.onopen = function(evt) {
		console.log("CONNECT !")
		$('#modal_login').modal('hide');
	}

	ws.onmessage = function(evt) {
		newmsg = dealMsg(evt.data)
		$('.all-msg').append(newmsg);
		$('.all-msg').scrollTop( $('.all-msg')[0].scrollHeight );

	}

	ws.onclose = function(evt) {
		console.log("CONNECT CLOSE")
		var msg = '<div class="row time-style" >连接已断开,请刷新</div >';
		$('.all-msg').append(msg);
		$('.all-msg').scrollTop( $('.all-msg')[0].scrollHeight );
		ws = null;
	}

	ws.onerror = function(evt) {
		console.log("ERROR: " + evt.data);
	}
});

$('#button_send').click(function (){
	var msg = $('#input_msg').val();
	if (msg == '') {
		alert('消息不能为空!');
		return ;
	}
	if (!ws) {
		return false;
	}
	ws.send(msg);
	$('#input_msg').val('');
});


$('#input_msg').keyup(function (e){
	var ev = window.event||e;
	if (ev.keyCode == 13) {
		$('#button_send').click();
	}

});

$('#input_loginname').keyup(function (e){
    var ev = window.event || e;
    if (ev.keyCode == 13) {
        $('#button_login').click()
    }
});

function dealMsg(recvdata) {
	var newmsg = '';
	var obj = JSON.parse(recvdata);

	var msg = obj['msg'];
	var time = obj['time'];
	var name = obj['name'];
	msg = msg.replace(/\n/g,'<br>');


	var needtime = false;
	if (lasttime == '') {
		needtime = true;
	} else if (time.substr(0,15) > lasttime.substr(0,15) ) {
		needtime = true;
	} else if (parseInt(time[15]) >= parseInt(lasttime[15]) + 5) {
		needtime = true;
	}

	if (needtime) {
		newmsg += '<div class="row time-style" > '+ time.substr(11,5) +' </div >';
		lasttime = time;
	}

	newmsg += '<div class="row single-msg" >';
        newmsg += '<div class="user-style">'+name+' >: </dev></br>';
	newmsg += '	<div class="msg-style">'+msg+'</div>';
	newmsg += '</div>';
	return newmsg;
}

function getUrlParam(name) {
	var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)"); //构造一个含有目标参数的正则表达式对象
	var r = window.location.search.substr(1).match(reg);  //匹配目标参数
	//if (r != null) return unescape(r[2]); return null; //返回参数值
	if (r != null) return decodeURI(r[2]); return null; //返回参数值
}