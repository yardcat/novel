import axios from 'axios';
import qs from 'qs';

import { Config } from './Config';

let connected = true;

function extractResponse(response) {
  return JSON.parse(response.data);
}

const CallAPI = (path, params, callback) => {
  if (!connected) {
    return false;
  }
  var url = Config.API_URL + '/' + path;
  if (Config.DEBUG) {
    console.log('request:', url, params);
  }
  axios
    .post(url, qs.stringify(params))
    .then((response) => {
      callback(extractResponse(response));
      connected = true;
    })
    .catch((error) => {
      if (Config.DEBUG) {
        console.error('Error:', error);
      }
      if (error.code === 'ERR_NETWORK') {
        connected = false;
      }
      console.log('net disconnected, %s', error);
    });
};

export { CallAPI };
