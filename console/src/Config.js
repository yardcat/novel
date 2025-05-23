import { CallAPI } from './Net';

class Config {
  static API_URL = 'http://127.0.0.1:8899';
  static SOCKET_URL = 'ws://localhost:8899/ws';
  static UPDATE_INTERVAL = 1000;
  static DEBUG = false;
  static collectable = [];
}

export function initConfig() {
  CallAPI('world/get_ui_info', {}, (response) => {
    Config.collectable = response.Collectable;
  });
}

export { Config };
