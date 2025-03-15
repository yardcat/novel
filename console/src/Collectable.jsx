import { useEffect, useState } from 'react';
import { InputNumber, List, Modal } from 'antd';
import CallAPI from './Net';
import Config from './Config';

const Collectable = ({ showCollect, setShowCollect, setAction }) => {
  const onSubmit = () => {
    const listItems = document.querySelectorAll('#collectable .ant-list-item');
    const collectableData = Array.from(listItems)
      .map((item) => {
        const label = item.querySelector('span').innerText;
        const number = item.querySelector('.ant-input-number-input').value;
        return { item: label, count: parseInt(number, 10) };
      })
      .filter((item) => item.count > 0);

    const jsonData = JSON.stringify({ items: collectableData });
    CallAPI('player/collect', { items: jsonData }, (response) => {
      if (response['items']) {
        response['action'] = 'collect';
        let log = '';
        response['items'].forEach((item) => {
          log += item['count'] + ' ' + item['item'] + ',';
        });
        response['log'] = log;
        setAction(response);
      }
    });
    setShowCollect(false);
  };

  return (
    <Modal
      id="collectable"
      title="Collectable"
      open={showCollect}
      onOk={onSubmit}
      onCancel={() => setShowCollect(false)}
    >
      <List
        itemLayout="horizontal"
        dataSource={Config.collectable}
        renderItem={(item, index) => (
          <List.Item>
            <span>{item}</span>
            <InputNumber type="number" defaultValue={0} />
          </List.Item>
        )}
      />
    </Modal>
  );
};

export default Collectable;
