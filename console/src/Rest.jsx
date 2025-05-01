import { useState, useEffect } from 'react';
import { Card, List } from 'antd';
import { socket } from './Socket';

const Rest = () => {
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
    <Card style={{ height: '60vh' }}>
      <List
        dataSource={chats}
        renderItem={(item) => <List.Item>{item}</List.Item>}
      />
    </Card>
  );
};

export { Rest };
