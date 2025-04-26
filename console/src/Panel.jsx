import { useState, useEffect } from 'react';
import { Card, Button, Tag } from 'antd';
import { Config } from './Config';
import { CallAPI } from './Net';
import { socket } from './Socket';

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
  width: '40vw',
};

const Panel = ({ role, info, intent }) => {
  return (
    <div style={panelStyle}>
      <p>name: {info.name}</p>
      <p>
        HP: {info.HP} / {info.maxHP}
      </p>
      <p> strength: {info.strength} </p>
      <p> defense: {info.defense} </p>
      {role == 'actor' && <p> energy: {info.energy} </p>}
      {role == 'enemy' && (
        <p>
          intent: {intent.action} {intent.actionValue}
        </p>
      )}
      <p>
        Buff:
        {info.buffs &&
          Object.values(info.buffs).map((buff) => (
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
