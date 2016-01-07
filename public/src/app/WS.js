export default
class WS {
	constructor(url) {
		this._url = url;
		this._isOpen = false;
	}
	openAsync() {
		return new Promise((resolve, reject) => {
			this._ws = new WebSocket(this._url);
			this._ws.onopen = (e) => {
				console.log('ws open success');
				this._isOpen = true;
				resolve(e);
			};
			this._ws.onerror = (e) => {
				console.log('ws open fail');
				this._ws = null;
				this._isOpen = true;
				reject(e);
			};
		});
	}
	readAsync() {
		return new Promise((resolve, reject) => {
			if (!this._isOpen) {
				reject(new Error('WS not open'));
			}
			this._ws.onmessage = (e) => {
				var data;
				try {
					data = JSON.parse(e.data);
				} catch(error) {
					data = e.data;
				}
				resolve(data);
			};
			this._ws.onerror = (e) => {
				console.log('ws read fail');
				this._ws = null;
				this._isOpen = true;
				reject(e);
			};
		});
	}
	write(data) {
		if (!this._isOpen) {
			throw (new Error('WS not open'));
		}
		this._ws.send(JSON.stringify(data));
	}
	close() {
		if (!this._isOpen) {
			return;
		}
		this._ws.close();
	}
}
