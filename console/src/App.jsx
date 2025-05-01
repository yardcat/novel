import { useState, useEffect } from 'react';
import { Navigator } from './Navigator';
import { PlayerInfo } from './PlayerInfo';
import { Bag } from './Bag';
import { Action } from './Action';
import { Chat } from './Chat';
import { Deck } from './Deck';
import { socket } from './Socket';
import { initConfig } from './Config';
import { Layout, Tabs } from 'antd';

const { Header, Footer, Sider, Content } = Layout;

const App = () => {
  const [apiHandlers, setApiHandlers] = useState({});

  const addApiHandler = (path, handler) => {
    setApiHandlers((prevHandlers) => ({
      ...prevHandlers,
      [path]: handler,
    }));
  };

  useEffect(() => {
    initConfig();
    socket.initSocket();
  }, []);

  const items = [
    {
      key: '1',
      label: 'Deck',
      children: <Deck></Deck>,
    },
    {
      key: '2',
      label: 'Chat',
      children: <Chat></Chat>,
    },
  ];

  return (
    <Layout>
      <Header>
        <Navigator
          apiHandlers={apiHandlers}
          addApiHandler={addApiHandler}
        ></Navigator>
      </Header>
      <Layout>
        <Sider width="20%">
          <PlayerInfo
            addApiHandler={addApiHandler}
            autoUpdate={true}
          ></PlayerInfo>
          <Bag addApiHandler={addApiHandler} autoUpdate={true}></Bag>
        </Sider>
        <Content>
          <Tabs defaultActiveKey="1" items={items} type="card"></Tabs>
        </Content>
        <Sider width="20%">
          <Action></Action>
        </Sider>
      </Layout>
    </Layout>
  );
};

export { App };
