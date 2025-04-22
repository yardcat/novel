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
  constructor({ name = '', HP = 0, maxHP = 0, energy = 0, strength = 0, defense = 0, buffs = [] } = {}) {
    Object.assign(this, { name, HP, maxHP, energy, strength, defense, buffs });
  }

  update(status) {
    this.name = status['name'] != null ? status['name'] : this.name;
    this.HP = status['HP'] != null ? status['HP'] : this.HP;
    this.maxHP = status['maxHP'] != null ? status['maxHP'] : this.maxHP;
    this.energy = status['energy'] != null ? status['energy'] : this.energy;
    this.strength = status['strength'] != null ? status['strength'] : this.strength;
    this.defense = status['defense'] != null ? status['defense'] : this.defense;
    this.buffs = status['buffs'] != null ? status['buffs'] : this.buffs;
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
  const [actorStatus, setActorStatus] = useState(new Status());
  const [enemyStatus, setEnemyStatus] = useState(new Status());

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
    });
  };

  const updateStatus = (results) => {
    if (results != null) {
      let as = new Status();
      Object.assign(as, actorStatus);
      as.update(results['actorStatus']);
      setActorStatus(as);
      let es = new Status();
      Object.assign(es, enemyStatus);
      es.update(results['enemyStatus']);
      setEnemyStatus(es);
    }
  };

  const sendCards = (cards) => {
    let cards_param = cards.join(',');
    CallAPI('world/send_cards', { cards: cards_param }, (reply) => {
      updateStatus(reply);
      setDrawCount(reply.drawCount);
    });
  };

  const discardCards = (cards) => {
    let cards_param = cards.join(',');
    CallAPI('world/discard_cards', { cards: cards_param }, (reply) => {
      setDiscardCount(reply.discardCount);
    });
  };

  const endTurn = () => {
    const sendTurnInfo = new EndTurn();
    CallAPI('world/end_turn', {}, (reply) => {
      updateStatus(reply);
    });
  };

  const handleChooseEvent = () => {
    CallAPI('world/card_choose_event', { event: selectedEvent }, (reply) => {
      updateStatus(reply);
    });
    setIsModalVisible(false);
  };

  const handleCancel = () => {
    setIsModalVisible(false);
  };

  return (
    <Card>
      <div style={{ display: 'flex', flexDirection: 'row', gap: '20px', border: '1px' }}>
        <Panel info={actorStatus} style={{ width: '45%' }}></Panel>
        <Panel info={enemyStatus} style={{ width: '45%' }}></Panel>
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
            <Button
              onClick={() => {
                discardCards(selectedCards);
                setHandCards(handCards.filter((card, idx) => !selectedCards.includes(idx)));
                setSelectedCards([]);
              }}
            >
              Discard
            </Button>
            <Button onClick={endTurn}>End</Button>
          </>
        )}
      </Card>

      <Modal title="Choose Event" open={isModalVisible} onOk={handleChooseEvent} onCancel={handleCancel}>
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
