import { useState, useEffect } from 'react';
import { message, Card, Button } from 'antd';
import { Panel } from './Panel';
import { CallAPI } from './Net';
import { MyCard } from './MyCard';
import { CardView } from './CardView';
import { socket } from './Socket';

import cardJson from '../../world/island/data/card/card.json';
import cardUpgradeJson from '../../world/island/data/card/card_upgrade.json';

const ChooseFrom = {
  HAND: 0,
  DRAW: 1,
  DISCARD: 2,
  EXHAUST: 3,
};

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

const Deck = () => {
  const [handCards, setHandCards] = useState([]);
  const [drawCount, setDrawCount] = useState(0);
  const [discardCount, setDiscardCount] = useState(0);
  const [exhaustCount, setExhaustCount] = useState(0);
  const [selectedCards, setSelectedCards] = useState([]);
  const [actorPanelInfo, setActorPanelInfo] = useState([]);
  const [enemyPanelInfo, setEnemyPanelInfo] = useState([]);
  const [selectedEnemy, setSelectedEnemy] = useState(null);
  const [cardViewConfig, setCardViewConfig] = useState({
    visible: false,
    cards: [],
    onOk: null,
  });

  const toggleCardSelection = (card) => {
    if (selectedCards.includes(card)) {
      setSelectedCards([]);
    } else {
      setSelectedCards([card]);
    }
  };

  const togglePanelSelection = (index) => {
    setSelectedEnemy(selectedEnemy === index ? null : index);
  };

  const getDrawCards = (finishFunc) => {
    CallAPI('card/show_draw_cards', {}, (reply) => {
      finishFunc(reply);
    });
  };

  const getDiscardCards = (finishFunc) => {
    CallAPI('card/show_discard_cards', {}, (reply) => {
      finishFunc(reply);
    });
  };

  const getExhaustCards = (finishFunc) => {
    CallAPI('card/show_exhaust_cards', {}, (reply) => {
      finishFunc(reply);
    });
  };

  useEffect(() => {
    const updateUI = (ev) => {
      setActorPanelInfo(ev.actorUI);
      setEnemyPanelInfo(ev.enemyUI);
      setDrawCount(ev.deckUI.drawCount);
      setDiscardCount(ev.deckUI.discardCount);
      setExhaustCount(ev.deckUI.exhaustCount);
      setHandCards(ev.deckUI.handCards);
    };

    socket.onMsg('event.CardUpdateUIEvent', (ev) => {
      updateUI(ev);
    });
  }, []);

  const sendCardWithChoosen = (card, choosen) => {
    const params = {
      card: card,
      choosen: choosen.join(','),
    };
    if (selectedEnemy == null && enemyPanelInfo.length == 1) {
      setSelectedEnemy('enemy-0');
      params.target = 0;
    } else {
      params.target = selectedEnemy.split('-')[1];
    }

    CallAPI('card/send_cards', params, (reply) => {
      if (reply.result != 'ok') {
        message.info("card can't be used");
        return;
      }
      setHandCards(handCards.filter((c, idx) => idx != card));
      setSelectedCards([]);
    });
  };

  const sendCard = (cards) => {
    if (selectedCards.length == 0) {
      message.info('no card selected');
      return;
    } else if (actorPanelInfo[0].energy <= 0) {
      setSelectedCards([]);
      message.info('no energy');
      return;
    }

    let cardIdx = cards[0];
    let cardId = handCards[cards[0]];
    let upgraded = cardId[cardId.length - 1] === '+';
    let cardInfo = upgraded ? cardUpgradeJson[cardId] : cardJson[cardId];
    if (cardInfo.cost > actorPanelInfo[0].energy) {
      message.info('energy not enough');
      setSelectedCards([]);
      return;
    }

    if (cardInfo.values['choose_count'] != null) {
      let count = cardInfo.values['choose_count'];
      let from = Number(cardInfo.values['choose_from']);
      message.info(`select ${count} cards`);
      switch (from) {
        case ChooseFrom.HAND:
          let chooseCards = handCards.map((_, idx) => {
            return idx == cardIdx ? null : handCards[idx];
          });
          showCardView(chooseCards, (selected) => {
            sendCardWithChoosen(cardIdx, selected);
          });
          break;
        case ChooseFrom.DRAW:
          getDrawCards((reply) => {
            showCardView(reply.cards, (selected) => {
              sendCardWithChoosen(cardIdx, selected);
            });
          });
          break;
        case ChooseFrom.DISCARD:
          getDiscardCards((reply) => {
            showCardView(reply.cards, (selected) => {
              sendCardWithChoosen(cardIdx, selected);
            });
          });
          break;
        case ChooseFrom.EXHAUST:
          getExhaustCards((reply) => {
            showCardView(reply.cards, (selected) => {
              sendCardWithChoosen(cardIdx, selected);
            });
          });
          break;
      }
    } else {
      sendCardWithChoosen(cardIdx, []);
    }
  };

  const discardCards = (cards) => {
    let cards_param = cards.join(',');
    CallAPI('card/discard_cards', { cards: cards_param }, (reply) => {
      setDiscardCount(reply.discardCount);
    });
  };

  const endTurn = () => {
    CallAPI('card/end_turn', {}, (reply) => {});
  };

  const showCardView = (cards, onOk) => {
    setCardViewConfig({ visible: true, cards, onOk });
  };

  return (
    <Card>
      {/* enemy region */}
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
        {/* draw cards region */}
        <Card
          style={{ width: '6vw' }}
          onClick={() => {
            getDrawCards((reply) => {
              showCardView(reply.cards, (selected) => {
                console.log('selected draw cards', selected);
                setCardViewConfig({ ...cardViewConfig, visible: false });
              });
            });
          }}
        >
          draw
          <h3>{drawCount}</h3>
        </Card>

        {/* cards region */}
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

        {/* discard and exhaust region */}
        <div style={{ display: 'flex', flexDirection: 'column' }}>
          <Card
            tyle={{ width: '5vw' }}
            size="small"
            onClick={() => {
              getDiscardCards((reply) => {
                showCardView(reply.cards, (selected) => {
                  console.log('selected discard cards', selected);
                  setCardViewConfig({ ...cardViewConfig, visible: false });
                });
              });
            }}
          >
            discard
            <h3>{discardCount}</h3>
          </Card>
          <Card
            tyle={{ width: '5vw' }}
            size="small"
            onClick={() => {
              getExhaustCards((reply) => {
                showCardView(reply.cards, (selected) => {
                  console.log('selected exhaust cards', selected);
                  setCardViewConfig({ ...cardViewConfig, visible: false });
                });
              });
            }}
          >
            exhaust
            <h3>{exhaustCount}</h3>
          </Card>
        </div>
      </div>

      {/* action region */}
      <Card>
        <Button
          onClick={() => {
            sendCard(selectedCards);
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

      {cardViewConfig.visible && (
        <CardView
          cards={cardViewConfig.cards}
          visible={cardViewConfig.visible}
          onCancel={() =>
            setCardViewConfig({ ...cardViewConfig, visible: false })
          }
          onOk={(selected) => {
            if (cardViewConfig.onOk) cardViewConfig.onOk(selected);
            setCardViewConfig({ ...cardViewConfig, visible: false });
          }}
        />
      )}
    </Card>
  );
};

export { Deck };
