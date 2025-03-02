import { React, useState, useEffect } from 'react';
import { List } from 'antd';

function updateAction(response, setAction) {
  setAction(response);
}

const Action = ({ addApiHandler }) => {
  const [actions, setAction] = useState([]);

  useEffect(() => {
    addApiHandler('player/collect', (response) => { updateAction(response, setAction) });
  }, []);

  return (
    <div style={{ height: '300px', overflowY: 'auto', border: '1px solid #ccc', padding: '10px' }}>
      <List
        itemLayout="horizontal"
        dataSource={actions}
        renderItem={(item) => (
          <List.Item>
            <List.Item.Meta
              title={`${item.endpoint} - ${item.timestamp}`}
              description={
                <div>
                  <p><strong>Params:</strong> {JSON.stringify(item.params)}</p>
                </div>
              }
            />
          </List.Item>
        )}
      />
    </div>
  );
};

export default Action;