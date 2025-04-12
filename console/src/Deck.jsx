import { useState, useEffect } from 'react';
import { Badge, Card, Button, Tag } from 'antd';
import { Config } from './Config';
import { CallAPI } from './Net';

class StartInfo {
  actor = {};
  enemy = {};
  handCard = [];
}

class RecvTurnInfo {
  handCard = [];
  discardCard = 0;
  deckCard = 0;
  health = 0;
  status = [];
}

class SendTurnInfo {
  useCard = [];
}

const Status = ({ name, count }) => {
  return (
    <Tag bordered={false} color="success">
      {name} : {count}
    </Tag>
  );
};

const MyCard = ({ name, isSelected, onClick }) => {
  return (
    <div
      style={{
        border: '1px solid black',
        padding: '50px',
        margin: '1px',
        backgroundColor: isSelected ? 'lightblue' : 'white',
        cursor: 'pointer',
      }}
      onClick={onClick}
    >
      {name}
    </div>
  );
};

const Deck = () => {
  const [turnInfo, setTurnInfo] = useState({
    handCard: ['1', '2'],
  });
  const [drawCount, setDrawCount] = useState(0);
  const [selectedCards, setSelectedCards] = useState([]);

  const toggleCardSelection = (card) => {
    if (selectedCards.includes(card)) {
      setSelectedCards(selectedCards.filter((c) => c !== card));
    } else {
      setSelectedCards([...selectedCards, card]);
    }
  };

  const startPlay = () => {
    const sendTurnInfo = new SendTurnInfo();
    sendTurnInfo.useCard = selectedCards;
    CallAPI('card/start', (startInfo) => {});
  };

  const sendTurnInfo = () => {
    const sendTurnInfo = new SendTurnInfo();
    sendTurnInfo.useCard = [1, 2, 3];
    CallAPI('card/send_turn', () => {});
  };

  const endPlay = () => {
    const sendTurnInfo = new SendTurnInfo();
    sendTurnInfo.useCard = [1, 2, 3];
    CallAPI('card/end_turn', () => {});
  };

  return (
    <Card title="Deck">
      <Card>
        <Status name="Health" count={1} />
        <Status name="dfs" count={2} />
      </Card>

      <div style={{ display: 'flex', flexDirection: 'row', gap: '20px' }}>
        <Card style={{ width: '10%' }}>
          draw
          <h3>{drawCount}</h3>
        </Card>

        <div style={{ display: 'flex', flexDirection: 'row', width: '80%' }}>
          {turnInfo.handCard &&
            turnInfo.handCard.map((card) => (
              <MyCard
                key={card}
                name={card}
                isSelected={selectedCards.includes(card)}
                onClick={() => toggleCardSelection(card)}
              />
            ))}
        </div>

        <Card>
          <h3>Discard Area</h3>
        </Card>
      </div>

      <Card>
        <Button onClick={startPlay}>Start</Button>
        <Button onClick={sendTurnInfo}>Send</Button>
        <Button onClick={endPlay}>End Turn</Button>
      </Card>
    </Card>
  );
};

export { Deck };
