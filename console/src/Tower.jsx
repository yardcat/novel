import { useState, useEffect } from 'react';
import { message, Select, Button } from 'antd';
import _ from 'lodash';
import { socket } from './Socket';
import { CallAPI } from './Net';
import { Deck } from './Deck';
import { Shop } from './Shop';
import { Rest } from './Rest';
import { Destiny } from './Destiny';
import { CardModal } from './Modal';
import { Potion } from './Potion';
import { Relic } from './Relic';

import cardJson from '../../world/island/data/card/card.json';

const Scene = {
  NONE: 0,
  COMBAT: 1,
  SHOP: 2,
  REST: 3,
  DESTINY: 4,
};

const toRoomString = (room) => {
  switch (room) {
    case 1:
      return 'fight';
    case 2:
      return 'shop';
    case 3:
      return 'rest';
    case 4:
      return 'destiny';
  }
};

const toRoomInt = (room) => {
  switch (room) {
    case 'fight':
      return 1;
    case 'shop':
      return 2;
    case 'rest':
      return 3;
    case 'destiny':
      return 4;
  }
};

const toBonusInt = (type) => {
  switch (type) {
    case 'cards':
      return 0;
    case 'potions':
      return 1;
    case 'relics':
      return 2;
  }
};

const Tower = () => {
  const [scene, setScene] = useState(null);
  const [difficuty, setDifficuty] = useState('Easy');
  const [modal, setModal] = useState(null);

  useEffect(() => {
    socket.onMsg('event.CardEnterRoomDone', (ev) => {
      ChangeScene(ev.type);
    });

    socket.onMsg('event.CardCombatWin', (ev) => {
      let rooms = _.map(ev.next_floor, (v) => {
        return toRoomString(v);
      });
      ev.bonus.rooms = rooms;
      cardModal.showCardModal('bonus', ev.bonus, (submit) => {
        const keys = _.keys(submit).map((v) => toBonusInt(v));
        const values = _.values(submit);
        const bonusArray = [];
        for (let i = 0; i < keys.length; i++) {
          if (!_.isNil(keys[i])) {
            bonusArray.push({
              type: keys[i],
              name: values[i],
            });
          }
        }
        CallAPI(
          'card/choose_bonus',
          { bonus: JSON.stringify(bonusArray) },
          (reply) => {
            console.log('bonus accept %s', reply);
            cardModal.hideCardModal();
          },
        );

        CallAPI(
          'card/next_floor',
          { room: toRoomInt(submit.rooms) },
          (reply) => {
            console.log('bonus accept %s', reply);
            cardModal.hideCardModal();
          },
        );
      });
    });

    socket.onMsg('event.CardCombatLose', (ev) => {
      message.info('you lose');
      ChangeScene(Scene.NONE);
    });
  }, []);

  const cardModal = {
    showCardModal: (type, choices, handler) => {
      setModal({
        type: type,
        choices: choices,
        visible: true,
        handler: handler,
      });
    },
    hideCardModal: () => {
      setModal({
        visible: false,
      });
    },
  };

  const startPlay = () => {
    if (difficuty === 'Difficuty') {
      message.error('Please select difficuty');
      return;
    }
    CallAPI('world/card_start', {}, (reply) => {
      cardModal.showCardModal(
        'welcome',
        { events: reply.choices },
        (submit) => {
          CallAPI('card/welcome', { event: submit['events'] }, (reply) => {
            console.log('welcome %s', reply);
            cardModal.hideCardModal();
            ChangeScene(Scene.COMBAT);
          });
        },
      );
    });
  };

  const ChangeScene = (type) => {
    switch (type) {
      case Scene.COMBAT:
        setScene(<Deck modal={cardModal} />);
        break;
      case Scene.SHOP:
        setScene(<Shop />);
        break;
      case Scene.REST:
        setScene(<Rest />);
        break;
      case Scene.DESTINY:
        setScene(<Destiny />);
        break;
    }
  };
  return (
    <>
      <Potion />
      <Relic />
      {scene}
      <Select defaultValue="Easy" onChange={setDifficuty}>
        <Select.Option value="Easy">Easy</Select.Option>
        <Select.Option value="Normal">Normal</Select.Option>
        <Select.Option value="Hard">Hard</Select.Option>
      </Select>
      <Button onClick={startPlay}>Start</Button>
      <CardModal modal={modal} />
    </>
  );
};

export { Tower };
