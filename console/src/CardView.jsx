import { useState, useEffect } from 'react';
import { Modal } from 'antd';

const CardView = ({ cards, visible }) => {
  const [selectedCards, setSelectedCards] = useState([]);

  const toggleCardSelection = (card) => {
    if (selectedCards.includes(card)) {
      setSelectedCards(selectedCards.filter((c) => c !== card));
    } else {
      setSelectedCards([...selectedCards, card]);
    }
  };
  return (
    <Modal open={visible} onCancel={() => {}}>
      <div
        style={{
          display: 'flex',
          flexDirection: 'row',
          width: '90vw',
          flexFlow: 'wrap',
        }}
      >
        {cards &&
          cards.map((name, idx) => (
            <MyCard
              key={idx}
              name={name}
              isSelected={selectedCards.includes(idx)}
              onClick={() => toggleCardSelection(idx)}
            />
          ))}
      </div>
    </Modal>
  );
};

export { CardView };
