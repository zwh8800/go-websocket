import WS from './WS';

class Channel {

}

class Client {
	constructor(username, ws) {
		this.token = 0;
		this.username = username;
		this.ws = ws;
		this.promiseHandlers = {
			'GetChannelsResponse': {},
			'ListChannelsResponse': {},
			'JoinResponse': {},
			'QuitResponse': {},
			'OnlineResponse': {},
			'OfflineResponse': {}
		};
	}
	open() {
		return this.ws.openAsync();
	}
	async start() {
		while (true) {
			try {
				let response = await this.ws.readAsync();
				let promiseArray = this.promiseHandlers[response.type];
				let promiseHandler = promiseArray[response.token];
				promiseHandler.resolve(response);
				delete promiseArray[response.token];
			} catch (error) {

			}
		}
	}
	getChannels() {
		const token = this.token++;
		const message = {
			type: 'GetChannels',
			token: token,
			data: null
		};
		this.ws.write(message);
		return new Promise((resolve, reject) => {
			this.promiseHandlers['GetChannelsResponse'][token] = {
				resolve, reject
			};
		});
	}
	get currentChannel() {

	}
	listAllChannels() {
		const token = this.token++;
		const message = {
			type: 'ListChannels',
			token: token,
			data: null
		};
		this.ws.write(message);
		return new Promise((resolve, reject) => {
			this.promiseHandlers['ListChannelsResponse'][token] = {
				resolve, reject
			};
		});
	}
	join(channel) {

	}
	quit(channel) {

	}
	online() {

	}
	offline() {

	}
}

var client = new Client('zwh8800', new WS('ws://localhost:9999/ws'));

export default client;
