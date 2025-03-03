import { React, useState, useEffect } from 'react';
import { List, Card } from 'antd';

const Action = ({ addApiHandler,actions }) => {
  useEffect(() => {
    addApiHandler('player/collect', null);
  }, []);

  return (
    <Card title="Action" style={{ width: '1000px' }}>
      <List
        itemLayout="horizontal"
        dataSource={actions}
        renderItem={(item) => (
          <List.Item>
            <List.Item.Meta
              title={item.action}
              description={
                <div> {item.log} </div>
              }
            />
          </List.Item>
        )}
      />
    </Card>
  );
};

export default Action;