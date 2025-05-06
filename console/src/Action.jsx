import { useEffect, useState, useRef } from 'react';
import { List, Card } from 'antd';
import { use } from 'react';
import { socket } from './Socket.js';

const Action = () => {
  const [actions, setAction] = useState([]);
  const listRef = useRef(null);

  useEffect(() => {
    socket.onMsg('event.ActionUpdateEvent', (ev) => {
      setAction((prevActions) => [...prevActions, ev.action]);
    });
  }, []);

  useEffect(() => {
    if (listRef.current) {
      listRef.current.scrollTop = listRef.current.scrollHeight;
    }
  }, [actions]);

  return (
    <List
      style={{ height: '80vh', overflowY: 'scroll' }}
      size="small"
      bordered
      header={<div>Action</div>}
      itemLayout="horizontal"
      dataSource={actions}
      renderItem={(item) => <List.Item>{item}</List.Item>}
      ref={listRef}
    />
  );
};

export { Action };
