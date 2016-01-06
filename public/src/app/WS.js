export default
class WS {
	constructor(url) {
		this._ws = new WebSocket(url);
	}
	readAsync() {
		return new Promise((resolve, reject) => {
			this._ws.onmessage = function (e) {
				var data;
				try {
					data = JSON.parse(e.data);
				} catch(error) {
					data = e.data;
				}
				resolve(data);
			};
		});
	}
	write(data) {
		this._ws.send(JSON.stringify(data));
	}
}