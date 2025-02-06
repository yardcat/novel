import  { useState } from 'react';
import { Cascader, Flex, Switch } from 'antd';
import axios from 'axios';

const options = [
  {
    value: 'user',
    label: 'user',
  },
];

const API_URL = 'http://127.0.0.1:8899'

const CallAPI = (path, params) => {
  var url = API_URL + '/' + path
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
  var params = ['value1']
  CallAPI(path, params)
};

const Choose = () => {
  return (
    <Flex vertical gap="small" align="flex-start">
      <Cascader.Panel options={options} onChange={handleChange} />
    </Flex>
  );
};

export default Choose;