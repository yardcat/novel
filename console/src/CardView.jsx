/* eslint-disable react/prop-types */
import { useState } from 'react';
import { Modal, Button } from 'antd';
import { MyCard } from './MyCard';

// eslint-disable-next-line react/prop-types
const CardView = ({ cards, onOk, onCancel }) => {
  const [selectedCards, setSelectedCards] = useState([]);
  const visible = cards && cards.length > 0;

  const toggleCardSelection = (card) => {
    if (selectedCards.includes(card)) {
      setSelectedCards(selectedCards.filter((c) => c !== card));
    } else {
      setSelectedCards([...selectedCards, card]);
    }
  };

  return (
    <Modal
      width={'60vw'}
      open={visible}
      onCancel={onCancel}
      afterClose={() => {
        setSelectedCards([]);
      }}
      footer={[
        <Button key="ok" type="primary" onClick={() => onOk(selectedCards)}>
          OK
        </Button>,
        <Button key="back" onClick={onCancel}>
          Cancel
        </Button>,
      ]}
    >
      <div
        style={{
          display: 'flex',
          flexDirection: 'row',
          width: '90vw',
          height: '50vh',
          flexFlow: 'wrap',
        }}
      >
        {cards &&
          cards.map(
            (name, idx) =>
              name != null && (
                <MyCard
                  key={idx}
                  name={name}
                  isSelected={selectedCards.includes(idx)}
                  onClick={() => toggleCardSelection(idx)}
                />
              ),
          )}
      </div>
    </Modal>
  );
};

export { CardView };
