import { useState, useEffect } from 'react';
import { Card, Button, Tag } from 'antd';
import { Config } from './Config';
import { CallAPI } from './Net';
import { socket } from './Socket';
import { use } from 'react';

const Types = ['vulnerable', 'weak', 'strength', 'armor'];

const Buff = ({ type, value, turn }) => {
  let name = Types[type];
  return (
    <>
      <Tag bordered={false} color="success">
        {name} : {value} ({turn})
      </Tag>
    </>
  );
};

const panelStyle = {
  border: '1px solid gray',
  borderRadius: '10px',
  padding: '0px 10px',
  margin: '5px 0px',
  width: '15vw',
};

const Panel = ({ role, info, isSelected, onClick }) => {
  let enemyStyle = {
    border: '1px solid blue',
    borderRadius: '10px',
    padding: '0px 10px',
    margin: '5px 0px',
    width: '10vw',
    backgroundColor: 'white',
    cusor: 'pointer',
    userSelect: 'none',
  };
  enemyStyle.backgroundColor = isSelected ? 'lightgreen' : 'white';

  return (
    <div style={role == 'actor' ? panelStyle : enemyStyle} onClick={onClick}>
      <p>name: {info.name}</p>
      <p>
        HP: {info.HP} / {info.maxHP}
      </p>
      <p> strength: {info.strength} </p>
      <p> defense: {info.defense} </p>
      {role == 'actor' && <p> energy: {info.energy} </p>}
      {role == 'enemy' && (
        <strong>
          {info.intent.action} {info.intent.value}
          {info.intent.target}
        </strong>
      )}
      <p>
        Buff:
        {info.buffs &&
          info.buffs.map((buff) => (
            <Buff
              key={buff.type}
              type={buff.type}
              value={buff.value}
              turn={buff.turn}
            />
          ))}
      </p>
    </div>
  );
};

export { Panel };
