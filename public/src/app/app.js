import 'babel-polyfill';
import $ from 'jquery';
import WS from './WS';

$(async function () {
	var input = $('#input');
	var output = $('#output');
	var button = $('#send');

	var ws = new WS('ws://localhost:9999/ws');
	try {
		await ws.open();
	} catch (e) {
		return;
	}

	button.click(function () {
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
