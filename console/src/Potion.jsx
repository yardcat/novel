import { useState, useEffect } from 'react';
import { Button, Tooltip, Flex } from 'antd';
import { socket } from './Socket';
import { CallAPI } from './Net';

const Potion = () => {
  const [items, setItems] = useState([]);

  const usePotion = (potion) => {
    CallAPI('card/use_potion', { name: potion }, (reply) => {
      console.log('potion use %s', reply);
    });
  };

  useEffect(() => {
    socket.onMsg('event.CardUpdatePotion', (data) => {
      setItems(data.potions);
    });
  }, []);

  return (
    <Flex
      gap="small"
      wrap
      style={{ paddingLeft: '10px', paddingRight: '10px' }}
    >
      Potions:
      {items.length > 0 &&
        items.map((item, index) => (
          <Tooltip title={item.description} key={index}>
            <Button
              value={item.name}
              color="cyan"
              variant="solid"
              size="small"
              onClick={() => usePotion(item.name)} // Pass item.name correctly
            >
              {item.name}
            </Button>
          </Tooltip>
        ))}
    </Flex>
  );
};

export { Potion };
