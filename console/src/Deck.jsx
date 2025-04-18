import { useState, useEffect } from 'react';
import { message, Select, Badge, Card, Button, Tag, Modal, Radio } from 'antd';
import { Panel } from './Panel';
import { Config } from './Config';
import { CallAPI } from './Net';
import { socket } from './Socket';

class StartInfo {
  difficuty = '';
}

class SendCards {
  cards = [];
}

class SendCardsResult {}

class EndTurn {}

class TurnInfo {
  handCards = [];
  discardCard = 0;
  deckCard = 0;
  health = 0;
  status = [];
}

class Status {
  constructor({
    actorHP = 0,
    actorMaxHP = 0,
    enemyHP = 0,
    enemyMaxHP = 0,
    energy = 0,
    strength = 0,
    defense = 0,
    buffs = [],
  } = {}) {
    Object.assign(this, { actorHP, actorMaxHP, enemyHP, enemyMaxHP, energy, strength, defense, buffs });
  }
}

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
  const [turnInfo, setTurnInfo] = useState({});
  const [handCards, setHandCards] = useState([]);
  const [drawCount, setDrawCount] = useState(0);
  const [discardCount, setDiscardCount] = useState(0);
  const [difficuty, setDifficuty] = useState('Easy');
  const [selectedCards, setSelectedCards] = useState([]);
  const [chooseEvents, setChooseEvents] = useState([]);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [selectedEvent, setSelectedEvent] = useState(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [hp, setHP] = useState(new Status());

  const toggleCardSelection = (card) => {
    if (selectedCards.includes(card)) {
      setSelectedCards(selectedCards.filter((c) => c !== card));
    } else {
      setSelectedCards([...selectedCards, card]);
    }
  };

  useEffect(() => {
    socket.onMsg('event.CardUpdateHandEvent', (data) => {
      setHandCards(data.cards);
    });
  }, []);

  const startPlay = () => {
    if (difficuty === 'Difficuty') {
      message.error('Please select difficuty');
      return;
    }
    const sendTurnInfo = new StartInfo();
    sendTurnInfo.difficuty = difficuty;
    CallAPI('world/card_start', {}, (reply) => {
      setTurnInfo(reply);
      setHandCards(reply.handCards);
      setDrawCount(reply.deckCount);
      setDiscardCount(0);
      setChooseEvents(reply.events);
      setIsModalVisible(true);
      setIsPlaying(true);
      setHP(
        new Status({
          actorHP: reply.actorHP,
          actorMaxHP: reply.actorMaxHP,
          enemyHP: reply.enemyHP,
          enemyMaxHP: reply.enemyMaxHP,
        }),
      );
    });
  };

  const handleUI = (results) => {
    if (results != null) {
      let newHP = new Status({
        actorHP: results['actorHP'] != null ? results['actorHP'] : hp.actorHP,
        actorMaxHP: results['actorMaxHP'] != null ? results['actorMaxHP'] : hp.actorMaxHP,
        enemyHP: results['enemyHP'] != null ? results['enemyHP'] : hp.enemyHP,
        enemyMaxHP: results['enemyMaxHP'] != null ? results['enemyMaxHP'] : hp.enemyMaxHP,
      });
      setHP(newHP);
    }
  };

  const sendCards = (cards) => {
    let cards_param = cards.join(',');
    CallAPI('world/send_cards', { cards: cards_param }, (reply) => {
      let results = reply['Results'];
      handleUI(results);
    });
  };

  const endTurn = () => {
    const sendTurnInfo = new EndTurn();
    CallAPI('world/end_turn', (turnInfo) => {});
  };

  const handleOk = () => {
    CallAPI('world/card_choose_event', { event: selectedEvent }, (reply) => {
      let results = reply['Results'];
      handleUI(results);
    });
    setIsModalVisible(false);
  };

  const handleCancel = () => {
    setIsModalVisible(false);
  };

  return (
    <Card title="Deck">
      <div style={{ display: 'flex', flexDirection: 'row', gap: '20px', border: '1px' }}>
        <Panel info={turnInfo} style={{ width: '45%' }}></Panel>
        <Panel info={turnInfo} style={{ width: '45%' }}></Panel>
      </div>

      <div style={{ display: 'flex', flexDirection: 'row', gap: '20px' }}>
        <Card style={{ width: '10%' }}>
          draw
          <h3>{drawCount}</h3>
        </Card>

        <div style={{ display: 'flex', flexDirection: 'row', width: '80%' }}>
          {handCards &&
            handCards.map((name, idx) => (
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
        {!isPlaying && (
          <>
            <Select defaultValue="Easy" onChange={setDifficuty}>
              <Select.Option value="Easy">Easy</Select.Option>
              <Select.Option value="Normal">Normal</Select.Option>
              <Select.Option value="Hard">Hard</Select.Option>
            </Select>
            <Button onClick={startPlay}>Start</Button>
          </>
        )}
        {isPlaying && (
          <>
            <Button
              onClick={() => {
                sendCards(selectedCards);
                setHandCards(handCards.filter((card, idx) => !selectedCards.includes(idx)));
                setSelectedCards([]);
              }}
            >
              Send
            </Button>
            <Button onClick={endTurn}>End</Button>
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
