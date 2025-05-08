import { useState, useEffect } from 'react';
import { Card, List } from 'antd';
import { socket } from './Socket';

const Destiny = () => {
  const [chats, setChat] = useState([]);
  const addChat = (chat) => {
    setChat((prevChats) => [...prevChats, chat]);
  };

  useEffect(() => {
    socket.onMsg('', (data) => {
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

export { Destiny };
