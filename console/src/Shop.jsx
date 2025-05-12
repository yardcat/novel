import { useState, useEffect } from 'react';
import { Card, Button, Tooltip, Flex } from 'antd';
import { socket } from './Socket';

import cardJson from '../../world/island/data/card/card.json';
import relicJson from '../../world/island/data/card/relic.json';
import potionJson from '../../world/island/data/card/potion.json';

const Shop = () => {
  const [items, setItems] = useState({ cards: [], relics: [], potions: [] });

  const buyItem = (type, name) => {
    console.log(`Buying ${type}: ${name}`);
    // Add API call or logic to handle buying the item
  };

  useEffect(() => {
    socket.onMsg('event.CardUpdateShopUI', (data) => {
      const cards = data.cards.map((item) => {
        const v = cardJson[item.name];
        return {
          name: v.name,
          description: v.description,
          cost: v.cost,
        };
      });

      const relics = data.relics.map((item) => {
        const v = relicJson[item.name];
        return {
          name: v.name,
          description: v.description,
          cost: v.cost,
        };
      });

      const potions = data.potions.map((item) => {
        const v = potionJson[item.name];
        return {
          name: v.name,
          description: v.description,
          cost: v.cost,
        };
      });

      setItems({ cards, relics, potions });
    });
  }, []);

  const renderItems = (items, type) => (
    <Flex direction="column" gap="small">
      {items.map((item, index) => (
        <Flex key={index} align="center" gap="small">
          <Tooltip title={item.description}>
            <span>{item.name}</span>
          </Tooltip>
          <span>Cost: {item.cost}</span>
          <Button size="small" onClick={() => buyItem(type, item.name)}>
            Buy
          </Button>
        </Flex>
      ))}
    </Flex>
  );

  return (
    <Card style={{ height: '60vh', padding: '10px' }}>
      <div>
        <h3>Cards</h3>
        {renderItems(items.cards, 'card')}
      </div>
      <div>
        <h3>Relics</h3>
        {renderItems(items.relics, 'relic')}
      </div>
      <div>
        <h3>Potions</h3>
        {renderItems(items.potions, 'potion')}
      </div>
    </Card>
  );
};

export { Shop };
