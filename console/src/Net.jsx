import axios from 'axios';
import qs from 'qs'

const API_URL = 'http://127.0.0.1:8899'

function extractResponse(response) {
  return JSON.parse(response.data)
}

const CallAPI = (path, params, callback) => {
    var url = API_URL + '/' + path
    console.log('request:', url, params);
    axios.post(url, qs.stringify(params))
      .then((response) => {
        callback(extractResponse(response));
      })
      .catch(error => {
        console.error('Error:', error);
      });
  }

  export default CallAPI