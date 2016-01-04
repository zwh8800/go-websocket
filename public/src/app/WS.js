export default
class WS {
	constructor(url) {
		this._ws = new WebSocket(url);
	}
	readAsync() {
		return new Promise((resolve, reject) => {
			this._ws.onmessage = function (e) {
				resolve(e.data);
			};
		});
	}
	write(data) {
		this._ws.send(data);
	}
}