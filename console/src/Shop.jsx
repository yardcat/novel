import { useState, useEffect } from 'react';
import { Card, List } from 'antd';
import { socket } from './Socket';

const Shop = () => {
  const [items, setItems] = useState([]);

  useEffect(() => {
    socket.onMsg('event.', (data) => {
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

export { Shop };
