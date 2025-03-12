import { React, useState, useEffect } from 'react';
import Navigator from './Navigator';
import PlayerInfo from './PlayerInfo';
import Bag from './Bag';
import Action from './Action';
import {socket} from './Socket';
import {initConfig} from './Config';
import { Flex } from 'antd';

const App = () => {
  const [apiHandlers, setApiHandlers] = useState({});
  const [actions, setAction] = useState([]);

  const addApiHandler = (path, handler) => {
    setApiHandlers(prevHandlers => ({
      ...prevHandlers,
      [path]: handler
    }));
  };

  const addAction = (action) => {
    setAction(prevActions => [...prevActions, action]);
  };

  useEffect(() => {
    initConfig();
    socket.initSocket();
  }, []);

  return (
    <Flex gap="middle">
      <Flex vertical gap="middle">
        <Navigator apiHandlers={apiHandlers} setAction={addAction}></Navigator>
        <PlayerInfo addApiHandler={addApiHandler} autoUpdate={true}></PlayerInfo>
        <Bag addApiHandler={addApiHandler} autoUpdate={true}></Bag>
      </Flex>
      <Flex>
        <Action addApiHandler={addApiHandler} actions={actions}></Action>
      </Flex>
    </Flex>
  );
};

export default App;