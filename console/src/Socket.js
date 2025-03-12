import Config from "./Config";

class Socket {
  constructor() {
    this.ws = null;
    this.eventHandlers = {};
  }

  initSocket() {
    this.ws = new WebSocket(Config.SOCKET_URL);

    this.ws.onopen = () => {
      console.log('WebSocket connection opened');
    };

    this.ws.onclose = () => {
      console.log('WebSocket connection closed');
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    this.ws.onmessage = (message) => {
      try {
        const msg = JSON.parse(message.data);
        const { event, data } = msg;
        if (this.eventHandlers[event]) {
          try {
            this.eventHandlers[event](data);
          } catch (callbackError) {
            console.error(`Callback for event "${event}" failed:`, callbackError);
          }
        } else {
          console.warn(`No handler registered for event "${event}"`);
        }
      } catch (parseError) {
        console.error('Failed to parse message:', parseError);
      }
    };
  }

  sendMsg(event, data) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.warn('WebSocket is not open. Message not sent.');
      return;
    }
    try {
      this.ws.send(JSON.stringify({ event, data }));
    } catch (error) {
      console.error('Failed to send message:', error);
    }
  }

  onMsg(event, callback) {
    if (typeof callback !== 'function') {
      console.error('Invalid callback provided for event:', event);
      return;
    }
    this.eventHandlers[event] = callback;
  }
}

const socket = new Socket();

export { socket };