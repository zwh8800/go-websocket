import 'babel-polyfill';
import $ from 'jquery';
import WS from './WS';

$(async function () {
	var ws = new WS('ws://localhost:9999/ws');
	var input = $('#input');
	var output = $('#output');
	$('#send').click(function () {
		ws.write({
			type: 'message',
			msg: input.val()
		});
		input.val('');
	});

	while (true) {
		let data = await ws.readAsync();
		output.append(data.msg);
	}
});
