import { useState, useEffect } from 'react';

import cardJson from '../../world/island/data/card/card.json';

const MyCard = ({ name, isSelected, onClick }) => {
  return (
    <div
      style={{
        width: '7vw',
        border: '1px solid black',
        margin: '1px',
        backgroundColor: isSelected ? 'lightblue' : 'white',
        cursor: 'pointer',
        textAlign: 'center', // Center content horizontally
      }}
      onClick={onClick}
    >
      <div style={{ marginTop: '10px ' }}>{name}</div>
      <div style={{ marginTop: '30px ' }}>
        <div style={{ fontSize: 'small' }}>{cardJson[name].description}</div>
        <div style={{ fontSize: 'small' }}>energy: {cardJson[name].cost}</div>
      </div>
    </div>
  );
};

export { MyCard };
