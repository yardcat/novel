import { useState, useEffect } from 'react';
import { Card, List } from 'antd';
import { socket } from './Socket';

const Chat = () => {
  const [chats, setChat] = useState([]);
  const addChat = (chat) => {
    setChat((prevChats) => [...prevChats, chat]);
  };

  useEffect(() => {
    socket.onMsg('world.CombatWinEvent', (data) => {
      addChat(data.Result);
    });
  }, []);

  return (
    <Card title="Chat" style={{ height: '40vh' }}>
      <List dataSource={chats} renderItem={(item) => <List.Item>{item}</List.Item>} />
    </Card>
  );
};

export { Chat };
