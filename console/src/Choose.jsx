import { Children, useState } from 'react';
import { Cascader, Flex, Card } from 'antd';
import axios from 'axios';
import qs from 'qs'

const options = [
  {
    value: 'player',
    label: 'player',
    children: [
      {
        value: 'get_player_info',
        label: 'PlayerInfo',
      },
      {
        value: 'get_bag',
        label: 'Bag',
      },
    ],
  },
];

const API_URL = 'http://127.0.0.1:8899'

const CallAPI = (path, params, success_callback) => {
  var url = API_URL + '/' + path
  console.log('request:', url, params);
  axios.post(url, qs.stringify(params))
    .then(success_callback)
    .catch(error => {
      console.error('Error:', error);
    });
}

const handleGetPlayerInfo = (setResponse) => {
  const params = {
    id: "0"
  };
  CallAPI('player/get_player_info', params, response => {
    console.log('Success:', response.data);
    setResponse(JSON.parse(response.data));
  });
};

const handleGetBag = () => {
  const params = {
    id: "0"
  };
  CallAPI('player/get_bag', params);
};

const routeHandlers = {
  'player/get_player_info': handleGetPlayerInfo,
  'player/get_bag': handleGetBag,
};

const handleChange = (value, setResponse) => {
  var path = value.join('/')
  const handler = routeHandlers[path];
  if (handler) {
    handler(setResponse);
  } else {
    console.error('No handler found for path:', path);
  }
};

const Choose = () => {
  const [response, setResponse] = useState(null);

  return (
    <Flex vertical gap="small" align="flex-start">
      <Cascader.Panel options={options} onChange={(value) => handleChange(value, setResponse)} />
      {response && (
        <Card title="Player Info">
          {Object.entries(response).map(([key, value]) => (
            <p key={key}>
              <strong>{key}:</strong> {value}
            </p>
          ))}
        </Card>
      )}
    </Flex>
  );
};

export default Choose;