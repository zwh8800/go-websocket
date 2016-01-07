import 'babel-polyfill';
import $ from 'jquery';
import WS from './WS';

$(async function () {
	var input = $('#input');
	var output = $('#output');
	var form = $('#sendForm');

	var ws = new WS('ws://localhost:9999/ws');
	try {
		await ws.openAsync();
	} catch (e) {
		return;
	}

	form.submit(function (e) {
		e.preventDefault();
		ws.write({
			type: 'message',
			msg: input.val()
		});
		input.val('');
	});

	while (true) {
		try {
			let data = await ws.readAsync();
			output.append(data.msg);
		} catch (e) {
			return;
		}
	}
});
