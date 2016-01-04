import 'babel-polyfill';
import $ from 'jquery';
import WS from './WS';

$(async function () {
	var ws = new WS('ws://localhost:9999/ws');
	$('#send').click(function () {
		ws.write($('#input').val());
		$('#input').val('');
	});

	while (true) {
		let data = await ws.readAsync();
		$('#output').append(data);
	}
});
