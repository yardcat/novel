import axios from 'axios';
import qs from 'qs';
import Config from "./Config";
import io from 'socket.io-client';

const API_URL = Config.API_URL;
const SOCKET_URL = Config.SOCKET_URL; d

let connected = true;
let socket; // 定义socket变量

function extractResponse(response) {
  return JSON.parse(response.data);
}

const CallAPI = (path, params, callback) => {
  if (!connected) {
    return false;
  }
  var url = API_URL + '/' + path;
  if (Config.DEBUG) {
    console.log('request:', url, params);
  }
  axios.post(url, qs.stringify(params))
    .then((response) => {
      callback(extractResponse(response));
      connected = true;
    })
    .catch(error => {
      if (Config.DEBUG) {
        console.error('Error:', error);
      }
      if (error.code === "ERR_NETWORK") {
        connected = false;
      }
      console.log('net disconnected, %s', error);
    });
}

function initSocket() {
  socket = io(SOCKET_URL, {
    reconnection: true,
    reconnectionDelay: 1000,
    reconnectionAttempts: 5
  });

  socket.on('connect', () => {
    console.log('Connected to WebSocket server');
    connected = true;
  });

  socket.on('disconnect', (reason) => {
    console.log('Disconnected from WebSocket server:', reason);
    connected = false;
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

export default {
  CallAPI,
  initSocket,
  sendSocketMessage,
  onSocketMessage
}