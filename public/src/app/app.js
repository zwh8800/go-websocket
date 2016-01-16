import 'babel-polyfill';
import $ from 'jquery';
import WS from './WS';

//$(async function () {
//	var input = $('#input');
//	var output = $('#output');
//	var form = $('#sendForm');
//
//	var ws = new WS('ws://localhost:9999/ws');
//	try {
//		await ws.openAsync();
//	} catch (e) {
//		return;
//	}
//
//	form.submit(function (e) {
//		e.preventDefault();
//		ws.write({
//			type: 'message',
//			msg: input.val()
//		});
//		input.val('');
//	});
//
//	while (true) {
//		try {
//			let data = await ws.readAsync();
//			output.append(data.msg);
//		} catch (e) {
//			return;
//		}
//	}
//});

import ChatClient from './chat';

var output = $('#output');
var input = $('#input');
var send = $('#send');
var sendForm = $('#sendForm');

async function showMyChannels() {
	output.append('<p>' + new Date() + '</p>' +
		'<p>showMyChannels</p>');
	console.log('<p>' + new Date() + '</p>' +
		'<p>showMyChannels</p>');
	let myChannels = await ChatClient.getChannels();
	output.append('<p>' + new Date() + '</p>' +
		'<p>' + JSON.stringify(myChannels) + '</p>');
	console.log('<p>' + new Date() + '</p>' +
		'<p>' + JSON.stringify(myChannels) + '</p>');
}
async function showAllChannels() {
	output.append('<p>' + new Date() + '</p>' +
		'<p>showAllChannels</p>');
	console.log('<p>' + new Date() + '</p>' +
		'<p>showAllChannels</p>');
	let allChannels = await ChatClient.listAllChannels();
	output.append('<p>' + new Date() + '</p>' +
		'<p>' + JSON.stringify(allChannels) + '</p>');
	console.log('<p>' + new Date() + '</p>' +
		'<p>' + JSON.stringify(allChannels) + '</p>');
}

async function init() {
	await ChatClient.open();
	ChatClient.start();
	console.log('start');

	//showAllChannels();
	//showMyChannels();
	//showAllChannels();
}

init();

sendForm.submit(async function (e) {
	e.preventDefault();
	let channel = input.val();
	input.val('');
	send.attr('disabled', true);
	let joinData = await ChatClient.join(channel);
	send.attr('disabled', false);
	output.append('<p>' + JSON.stringify(joinData) + '</p>');
});
