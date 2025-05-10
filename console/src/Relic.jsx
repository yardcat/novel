import { useState, useEffect } from 'react';
import { Button, Tooltip, Flex } from 'antd';
import { socket } from './Socket';

class RelicModel {
  constructor(name, description) {
    Object.assign(this, { name, description });
  }
}

const Relic = () => {
  const [items, setItems] = useState([]);

  useEffect(() => {
    socket.onMsg('event.CardUpdateRelic', (data) => {
      setItems(data);
    });
  }, []);

  return (
    <Flex
      gap="small"
      wrap
      style={{ paddingLeft: '10px', paddingRight: '10px' }}
    >
      Relics:
      {items.map((item, index) => (
        <Tooltip title={item.description} key={index}>
          <Button value={item.name}>{item.name}</Button>
        </Tooltip>
      ))}
    </Flex>
  );
};

export { Relic };
