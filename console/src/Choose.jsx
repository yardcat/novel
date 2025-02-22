import  { Children, useState } from 'react';
import { Cascader, Flex, Switch } from 'antd';
import axios from 'axios';
import qs from 'qs'

const options = [
  {
    value: 'user',
    label: 'user',
    children: [
      {
        value: 'get_user_info',
        label: 'UserInfo',
      },
      {
        value: 'get_bag',
        label: 'Bag',
      },
    ],
  },
];

const API_URL = 'http://127.0.0.1:8899'

const CallAPI = (path, params) => {
  var url = API_URL + '/' + path
  console.log('request:', url, params);
  axios.post(url, params)
    .then(response => {
      console.log('Success:', response.data);
    })
    .catch(error => {
      console.error('Error:', error);
    });
}

const handleChange = (value) => {
  var path = value.join('/')
  var params = {"key1": "value1", "key2": "value2"}
  CallAPI(path, qs.stringify(params))
};

const Choose = () => {
  return (
    <Flex vertical gap="small" align="flex-start">
      <Cascader.Panel  options={options} onChange={handleChange} />
    </Flex>
  );
};

export default Choose;