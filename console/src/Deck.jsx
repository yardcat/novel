import { useState, useEffect } from 'react';
import { message, Select, Badge, Card, Button, Tag, Modal, Radio } from 'antd';
import { Config } from './Config';
import { CallAPI } from './Net';

class StartInfo {
  difficuty = '';
}

class RecvTurnInfo {
  handCards = [];
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
    handCards: [],
  });
  const [drawCount, setDrawCount] = useState(0);
  const [discardCount, setDiscardCount] = useState(0);
  const [difficuty, setDifficuty] = useState('Difficuty');
  const [selectedCards, setSelectedCards] = useState([]);
  const [chooseEvents, setChooseEvents] = useState([]);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [selectedEvent, setSelectedEvent] = useState(null);
  const [isPlaying, setIsPlaying] = useState(false);

  const toggleCardSelection = (card) => {
    if (selectedCards.includes(card)) {
      setSelectedCards(selectedCards.filter((c) => c !== card));
    } else {
      setSelectedCards([...selectedCards, card]);
    }
  };

  const startPlay = () => {
    if (difficuty === 'Difficuty') {
      message.error('Please select difficuty');
      return;
    }
    const sendTurnInfo = new StartInfo();
    sendTurnInfo.difficuty = difficuty;
    CallAPI('world/card_start', {}, (reply) => {
      setTurnInfo(reply);
      setDrawCount(reply.deckCount);
      setDiscardCount(0);
      setChooseEvents(reply.events);
      setIsModalVisible(true);
      setIsPlaying(true); // 设置为正在游戏中
    });
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

  const handleOk = () => {
    console.log('Selected Event:', selectedEvent);
    CallAPI('world/card_choose_event', { event: selectedEvent });
    setIsModalVisible(false);
  };

  const handleCancel = () => {
    setIsModalVisible(false);
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
          {turnInfo.handCards &&
            turnInfo.handCards.map((name, idx) => (
              <MyCard
                key={idx}
                name={name}
                isSelected={selectedCards.includes(idx)}
                onClick={() => toggleCardSelection(idx)}
              />
            ))}
        </div>

        <Card>
          discard
          <h3>{discardCount}</h3>
        </Card>
      </div>

      <Card>
        <Select defaultValue="Difficuty" onChange={setDifficuty}>
          <Select.Option value="Easy">Easy</Select.Option>
          <Select.Option value="Normal">Normal</Select.Option>
          <Select.Option value="Hard">Hard</Select.Option>
        </Select>
        {!isPlaying && <Button onClick={startPlay}>Start</Button>}
        {isPlaying && (
          <>
            <Button onClick={sendTurnInfo}>Send</Button>
            <Button onClick={endPlay}>End Turn</Button>
          </>
        )}
      </Card>

      <Modal title="Choose Event" open={isModalVisible} onOk={handleOk} onCancel={handleCancel}>
        <Radio.Group onChange={(e) => setSelectedEvent(e.target.value)} value={selectedEvent}>
          {chooseEvents.map((event, index) => (
            <Radio key={index} value={event}>
              {event}
            </Radio>
          ))}
        </Radio.Group>
      </Modal>
    </Card>
  );
};

export { Deck };
