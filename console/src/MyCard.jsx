import { useState, useEffect } from 'react';

import cardJson from '../../world/island/data/card/card.json';
import cardUpgradeJson from '../../world/island/data/card/card_upgrade.json';

const StringFormat = (str, values) => {
  return str.replace(/{(\w+)}/g, (match, key) => {
    return key in values ? values[key] : match;
  });
};

const MyCard = ({ name, isSelected, onClick }) => {
  let upgraded = name[name.length - 1] === '+';
  let json = upgraded ? cardUpgradeJson : cardJson;
  let description = StringFormat(json[name].description, json[name].values);

  return (
    <div
      style={{
        width: '120px',
        height: '150px',
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
        <div style={{ fontSize: 'small' }}>{description}</div>
        <div style={{ fontSize: 'small' }}>energy: {json[name].cost}</div>
      </div>
    </div>
  );
};

export { MyCard };
