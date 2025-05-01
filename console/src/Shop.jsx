import { useState, useEffect } from 'react';
import { Card, List } from 'antd';
import { socket } from './Socket';

class Item {
  constructor({ name, price = 0 } = {}) {
    Object.assign(this, { name, price });
  }
}

const Shop = () => {
  const [items, setItems] = useState([]);

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

export { Shop };
