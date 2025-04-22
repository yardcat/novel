import { useState, useEffect } from 'react';
import { Card, Button, Tag } from 'antd';
import { Config } from './Config';
import { CallAPI } from './Net';
import { socket } from './Socket';

const Buff = ({ name, value }) => {
  return (
    <Tag bordered={false} color="success">
      {name} : {value}
    </Tag>
  );
};

const panelStyle = {
  border: '1px solid gray',
  borderRadius: '10px',
  padding: '0px 10px',
  margin: '5px 0px',
  width: '40vw',
};

const Panel = ({ info }) => {
  return (
    <div style={panelStyle}>
      <p>name: {info.name}</p>
      <p>
        HP: {info.HP} / {info.maxHP}
      </p>
      <p> strength: {info.strength} </p>
      <p> defense: {info.defense} </p>
      <p> energy: {info.energy} </p>
      <p>Buff: {info.buffs && info.buffs.map((buff) => <Buff key={buff.name} name={buff.name} count={buff.turn} />)}</p>
    </div>
  );
};

export { Panel };
