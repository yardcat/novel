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

const Panel = ({ info }) => {
  return (
    <div style={{ bordered: '1px' }}>
      <p>name: {info.name}</p>
      <p>
        HP: {info.hp} / {info.maxHP}
      </p>
      <p> strength: {info.strength} </p>
      <p> defense: {info.defense} </p>
      {info.buffs && info.buffs.map((buff) => <Buff key={buff.name} name={buff.name} count={buff.value} />)}
    </div>
  );
};

export { Panel };
