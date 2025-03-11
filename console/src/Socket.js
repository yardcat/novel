import Config from "./Config";
import io from 'socket.io-client';

const SOCKET_URL = Config.SOCKET_URL;

let socket;

function initSocket() {
  socket = io(SOCKET_URL);

  socket.on('connect', () => {
    console.log('Connected to WebSocket server');
  });

  socket.on('disconnect', (reason) => {
    console.log('Disconnected from WebSocket server:', reason);
  });

  socket.on('error', (error) => {
    console.error('WebSocket error:', error);
  });
}

function sendSocketMessage(event, data) {
  if (connected) {
    socket.emit(event, data);
  } else {
    console.log('WebSocket is not connected');
  }
}

function onSocketMessage(event, callback) {
  if (connected) {
    socket.on(event, callback);
  } else {
    console.log('WebSocket is not connected');
  }
}

export { initSocket, sendSocketMessage, onSocketMessage };