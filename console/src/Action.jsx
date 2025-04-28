import { useEffect, useState } from 'react';
import { List, Card } from 'antd';
import { use } from 'react';
import { socket } from './Socket.js';

const Action = () => {
  const [actions, setAction] = useState([]);

  useEffect(() => {
    socket.onMsg('event.ActionUpdateEvent', (ev) => {
      setAction((prevActions) => {
        [...prevActions, ev.action];
      });
    });
  }, []);

  return (
    <List
      style={{ height: '80vh', overflowY: 'scroll' }}
      size="small"
      bordered
      header={<div>Action</div>}
      itemLayout="horizontal"
      dataSource={actions}
      renderItem={(item) => <List.Item>{item}</List.Item>}
    />
  );
};

export { Action };
