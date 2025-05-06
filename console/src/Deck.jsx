import { useState, useEffect } from 'react';
import { message, Select, Badge, Card, Button, Tag, Modal, Radio } from 'antd';
import { Panel } from './Panel';
import { Config } from './Config';
import { CallAPI } from './Net';
import { socket } from './Socket';

class StartInfo {
  difficuty = '';
}

class EndTurn {}

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

const PanelContainerStyle = {
  display: 'flex',
  flexDirection: 'row',
  gap: '20px',
  border: '1px',
};

const Deck = () => {
  const [handCards, setHandCards] = useState([]);
  const [drawCount, setDrawCount] = useState(0);
  const [discardCount, setDiscardCount] = useState(0);
  const [difficuty, setDifficuty] = useState('Easy');
  const [selectedCards, setSelectedCards] = useState([]);
  const [chooseType, setChooseType] = useState('');
  const [chooseEvents, setChooseEvents] = useState([]);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [selectedEvent, setSelectedEvent] = useState(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [actorPanelInfo, setActorPanelInfo] = useState([]);
  const [enemyPanelInfo, setEnemyPanelInfo] = useState([]);
  const [selectedEnemy, setSelectedEnemy] = useState(null);

  const toggleCardSelection = (card) => {
    if (selectedCards.includes(card)) {
      setSelectedCards(selectedCards.filter((c) => c !== card));
    } else {
      setSelectedCards([...selectedCards, card]);
    }
  };

  const togglePanelSelection = (index) => {
    setSelectedEnemy(selectedEnemy === index ? null : index);
  };

  useEffect(() => {
    const updateUI = (ev) => {
      setActorPanelInfo(ev.actorUI);
      setEnemyPanelInfo(ev.enemyUI);
      setDrawCount(ev.deckUI.drawCount);
      setDiscardCount(ev.deckUI.discardCount);
      setHandCards(ev.deckUI.handCards);
    };

    socket.onMsg('event.CardUpdateUIEvent', (ev) => {
      updateUI(ev);
    });

    socket.onMsg('event.CardCombatWin', (ev) => {
      showChooseModal('bonus', ev.bonus);
      setIsPlaying(false);
    });

    socket.onMsg('event.CardCombatLose', (ev) => {
      setIsPlaying(false);
      message.info('you lose');
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
      showChooseModal('welcome', reply.choices);
    });
  };

  const sendCards = (cards) => {
    if (selectedCards.length === 0) {
      message.info('no card selected');
      return;
    } else if (selectedEnemy == null) {
      message.info('no enemy selected');
      return;
    } else if (actorPanelInfo[0].energy <= 0) {
      message.info('no energy');
      return;
    }

    const params = {
      cards: cards.join(','),
      target: selectedEnemy.split('-')[1],
    };
    CallAPI('card/send_cards', params, (reply) => {
      setHandCards(
        handCards.filter((card, idx) => !selectedCards.includes(idx)),
      );
      setSelectedCards([]);
    });
  };

  const discardCards = (cards) => {
    let cards_param = cards.join(',');
    CallAPI('card/discard_cards', { cards: cards_param }, (reply) => {
      setDiscardCount(reply.discardCount);
    });
  };

  const endTurn = () => {
    const sendTurnInfo = new EndTurn();
    CallAPI('card/end_turn', {}, (reply) => {});
  };

  const showChooseModal = (type, choices) => {
    setChooseType(type);
    setChooseEvents(choices);
    setIsModalVisible(true);
  };

  const hideChooseModal = () => {
    setIsModalVisible(false);
  };

  const handleChooseEvent = (type) => {
    switch (chooseType) {
      case 'welcome':
        CallAPI('card/welcome', { event: selectedEvent }, (reply) => {
          setIsPlaying(true);
        });
        break;

      case 'bonus':
        CallAPI('card/choose_bonus', { event: selectedEvent }, (reply) => {
          console.log('bonus accept');
          showChooseModal('room', reply.rooms);
        });
        break;

      case 'room':
        CallAPI('card/enter_room', { event: selectedEvent }, (reply) => {
          setIsPlaying(true);
        });
        break;

      default:
        console.log('no type specified');
        break;
    }
    hideChooseModal();
  };

  const handleCancel = () => {
    setIsModalVisible(false);
  };

  return (
    <Card>
      {isPlaying && (
        <div style={PanelContainerStyle}>
          {actorPanelInfo.map((info, index) => (
            <Panel role="actor" info={info} key={index} />
          ))}
          {enemyPanelInfo.map((info, index) => (
            <Panel
              role="enemy"
              info={info}
              key={index}
              isSelected={selectedEnemy === `enemy-${index}`}
              onClick={() => togglePanelSelection(`enemy-${index}`)}
            />
          ))}
        </div>
      )}
      {isPlaying && (
        <div style={{ display: 'flex', flexDirection: 'row', gap: '20px' }}>
          <Card style={{ width: '6vw' }}>
            draw
            <h3>{drawCount}</h3>
          </Card>

          <div
            style={{
              display: 'flex',
              flexDirection: 'row',
              width: '90vw',
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

          <Card tyle={{ width: '5vw' }}>
            discard
            <h3>{discardCount}</h3>
          </Card>
        </div>
      )}
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
        title={chooseType}
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
