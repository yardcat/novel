import { useState, useEffect } from 'react';
import { Badge, Card, Button, Tag } from 'antd';
import { Config } from './Config';
import { CallAPI } from './Net';

class RecvTurnInfo {
  HandCard = [];
  DiscardCard = 0;
  DeckCard = 0;
  Health = 0;
  Dfs = 0;
  Status = [];
}

class SendTurnInfo {
  UseCard = [];
}

const Status = ({ name, count }) => {
  return (
    <Tag bordered={false} color="success">
      {name} : {count}
    </Tag>
  );
};

const MyCard = ({ name }) => {
  return <div>{name}</div>;
};

const Deck = () => {
  const [info, setInfo] = useState({});

  return (
    <Card title="Deck">
      <Card>
        <Status name="Health" count={1} />
        <Status name="dfs" count={2} />
      </Card>

      <div style={{ display: 'flex', flexDirection: 'row' }}>
        <Card>
          <h3>draw Area</h3>
        </Card>

        <Card style={{ width: '80%' }}>
          {info.cards && info.cards.map((card, index) => <MyCard key={index} name={card.name} />)}
        </Card>

        <Card>
          <h3>Discard Area</h3>
        </Card>
      </div>

      <Card>
        <Button onClick={startPlay()}>Start</Button>
        <Button onClick={sendTurnInfo}>Send</Button>
        <Button onClick={endPlay()}>End Turn</Button>
      </Card>
    </Card>
  );
};

function startPlay() {
  const sendTurnInfo = new SendTurnInfo();
  sendTurnInfo.UseCard = [1, 2, 3];
  CallAPI('card/start', () => {});
}

function sendTurnInfo() {
  const sendTurnInfo = new SendTurnInfo();
  sendTurnInfo.UseCard = [1, 2, 3];
  CallAPI('card/send_turn', () => {});
}

function endPlay() {
  const sendTurnInfo = new SendTurnInfo();
  sendTurnInfo.UseCard = [1, 2, 3];
  CallAPI('card/end_turn', () => {});
}

export { Deck };
