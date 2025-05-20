import { useState, useEffect } from 'react';
import { message, Card, Button } from 'antd';
import { Panel } from './Panel';
import { Config } from './Config';
import { CallAPI } from './Net';
import { MyCard } from './MyCard';
import { CardView } from './CardView';
import { socket } from './Socket';

class StartInfo {
  difficuty = '';
}

class EndTurn {}

const PanelContainerStyle = {
  display: 'flex',
  flexDirection: 'row',
  gap: '20px',
  border: '1px',
};

const Deck = ({ modal }) => {
  const [handCards, setHandCards] = useState([]);
  const [drawCount, setDrawCount] = useState(0);
  const [discardCount, setDiscardCount] = useState(0);
  const [selectedCards, setSelectedCards] = useState([]);
  const [actorPanelInfo, setActorPanelInfo] = useState([]);
  const [enemyPanelInfo, setEnemyPanelInfo] = useState([]);
  const [selectedEnemy, setSelectedEnemy] = useState(null);
  const [showCardView, setShowCardView] = useState(false);
  const [cardsToShow, setCardsToShow] = useState([]);

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
  }, []);

  const sendCards = (cards) => {
    if (selectedCards.length == 0) {
      message.info('no card selected');
      return;
    } else if (actorPanelInfo[0].energy <= 0) {
      message.info('no energy');
      return;
    }

    const params = {
      cards: cards.join(','),
    };
    if (selectedEnemy == null && enemyPanelInfo.length == 1) {
      setSelectedEnemy('enemy-0');
      params.target = 0;
    } else {
      params.target = selectedEnemy.split('-')[1];
    }

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

  return (
    <Card>
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
      <div style={{ display: 'flex', flexDirection: 'row', gap: '20px' }}>
        <Card
          style={{ width: '6vw' }}
          onClick={() => {
            CallAPI('card/show_draw_cards', {}, (reply) => {
              setCardsToShow(reply.cards);
              setShowCardView(true);
            });
          }}
        >
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

        <Card
          tyle={{ width: '5vw' }}
          onClick={() => {
            CallAPI('card/show_discard_cards', {}, (reply) => {
              setCardsToShow(reply.cards);
              setShowCardView(true);
            });
          }}
        >
          discard
          <h3>{discardCount}</h3>
        </Card>
      </div>
      <Card>
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
      </Card>
      <CardView
        cards={cardsToShow}
        visible={showCardView}
        onCancel={() => setShowCardView(false)}
      />
    </Card>
  );
};

export { Deck };
