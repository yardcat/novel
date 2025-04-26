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

class PanelInfo {
  constructor({
    name = '',
    HP = 0,
    maxHP = 0,
    energy = 0,
    strength = 0,
    defense = 0,
    buffs = {},
  } = {}) {
    Object.assign(this, { name, HP, maxHP, energy, strength, defense, buffs });
  }

  update(info) {
    this.name = info['name'] != null ? info['name'] : this.name;
    this.HP = info['HP'] != null ? info['HP'] : this.HP;
    this.maxHP = info['maxHP'] != null ? info['maxHP'] : this.maxHP;
    this.energy = info['energy'] != null ? info['energy'] : this.energy;
    this.strength = info['strength'] != null ? info['strength'] : this.strength;
    this.defense = info['defense'] != null ? info['defense'] : this.defense;
    this.buffs = info['buffs'] != null ? info['buffs'] : this.buffs;
  }
}

const MyCard = ({ name, isSelected, onClick }) => {
  return (
    <div
      style={{
        width: '7vw',
        border: '1px solid black',
        margin: '1px',
        textAlign: 'center',
        lineHeight: '16vh',
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
  const [actorPanelInfo, setActorPanelInfo] = useState(new PanelInfo());
  const [enemyPanelInfo, setEnemyPanelInfo] = useState(new PanelInfo());
  const [intent, setIntent] = useState({});

  const toggleCardSelection = (card) => {
    if (selectedCards.includes(card)) {
      setSelectedCards(selectedCards.filter((c) => c !== card));
    } else {
      setSelectedCards([...selectedCards, card]);
    }
  };

  useEffect(() => {
    const updateUI = (ev) => {
      for (const actor of ev.actorUI) {
        let as = new PanelInfo(actor);
        setActorPanelInfo(as);
      }
      for (const enemy of ev.enemyUI) {
        let es = new PanelInfo(enemy);
        setEnemyPanelInfo(es);
      }
      setDrawCount(ev.deckUI.drawCount);
      setDiscardCount(ev.deckUI.discardCount);
      setHandCards(ev.deckUI.handCards);
    };

    socket.onMsg('event.CardUpdateUIEvent', (ev) => {
      updateUI(ev);
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
      setChooseEvents(reply.events);
      setIsModalVisible(true);
      setIsPlaying(true);
      setIntent({ action: reply.action, actionValue: reply.actionValue });
    });
  };

  const sendCards = (cards) => {
    if (cards.length === 0 || actorPanelInfo.energy <= 0) {
      return;
    }
    let cards_param = cards.join(',');
    CallAPI('world/send_cards', { cards: cards_param }, (reply) => {});
    setHandCards(handCards.filter((card, idx) => !selectedCards.includes(idx)));
    setSelectedCards([]);
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
      setIntent({ action: reply.action, actionValue: reply.actionValue });
    });
  };

  const handleChooseEvent = () => {
    CallAPI('world/card_choose_event', { event: selectedEvent }, (reply) => {});
    setIsModalVisible(false);
  };

  const handleCancel = () => {
    setIsModalVisible(false);
  };

  return (
    <Card>
      <div
        style={{
          display: 'flex',
          flexDirection: 'row',
          gap: '20px',
          border: '1px',
        }}
      >
        <Panel
          role="actor"
          info={actorPanelInfo}
          style={{ width: '45%' }}
        ></Panel>
        <Panel
          role="enemy"
          info={enemyPanelInfo}
          intent={intent}
          style={{ width: '45%' }}
        ></Panel>
      </div>

      <div style={{ display: 'flex', flexDirection: 'row', gap: '20px' }}>
        <Card style={{ width: '8%' }}>
          draw
          <h3>{drawCount}</h3>
        </Card>

        <div
          style={{
            display: 'flex',
            flexDirection: 'row',
            width: '80%',
            flexFlow: 'wrap',
          }}
        >
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

        <Card tyle={{ width: '8%' }}>
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
              }}
            >
              Send
            </Button>
            <Button
              onClick={() => {
                discardCards(selectedCards);
                setHandCards(
                  handCards.filter((card, idx) => !selectedCards.includes(idx)),
                );
                setSelectedCards([]);
              }}
            >
              Discard
            </Button>
            <Button onClick={endTurn}>End</Button>
          </>
        )}
      </Card>

      <Modal
        title="Choose Event"
        open={isModalVisible}
        onOk={handleChooseEvent}
        onCancel={handleCancel}
      >
        <Radio.Group
          onChange={(e) => setSelectedEvent(e.target.value)}
          value={selectedEvent}
        >
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
