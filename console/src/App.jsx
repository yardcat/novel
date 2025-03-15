import { useState, useEffect } from 'react';
import Navigator from './Navigator';
import PlayerInfo from './PlayerInfo';
import { Bag } from './Bag';
import { Action } from './Action';
import { Chat } from './Chat';
import { socket } from './Socket';
import { initConfig } from './Config';
import { Layout } from 'antd';

const { Header, Footer, Sider, Content } = Layout;

const App = () => {
  const [apiHandlers, setApiHandlers] = useState({});
  const [actions, setAction] = useState([]);

  const addApiHandler = (path, handler) => {
    setApiHandlers((prevHandlers) => ({
      ...prevHandlers,
      [path]: handler,
    }));
  };

  const addAction = (action) => {
    setAction((prevActions) => [...prevActions, action]);
  };

  useEffect(() => {
    initConfig();
    socket.initSocket();
  }, []);

  return (
    <Layout>
      <Header>
        <Navigator apiHandlers={apiHandlers} addApiHandler={addApiHandler} setAction={addAction}></Navigator>
      </Header>
      <Sider>
        <PlayerInfo addApiHandler={addApiHandler} autoUpdate={true}></PlayerInfo>
        <Bag addApiHandler={addApiHandler} autoUpdate={true}></Bag>
      </Sider>
      <Content>
        <Chat></Chat>
        <Action addApiHandler={addApiHandler} actions={actions}></Action>
      </Content>
    </Layout>
  );
};

export { App };
