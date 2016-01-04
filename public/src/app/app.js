import $ from 'jquery';
$(function () {
	var ws = new WebSocket("ws://localhost:9999/ws");
	ws.onmessage = function (e) {
		$('#output').append(e.data + '\r\n');
	};
	$('#send').click(function () {
		ws.send($('#input').val());
		$('#input').val('');
	});



});
