import { useState, useEffect } from 'react';
import { Card } from 'antd';
import Config from './Config';
import CallAPI from './Net';
const API_PATH = 'player/get_bag';

function GetBag(info, setInfo) {
  if (info.items) {
    info.items.sort((a, b) => a.name.localeCompare(b.name));
  }
  setInfo(info);
}

const Bag = ({ addApiHandler, autoUpdate }) => {
  const [info, setInfo] = useState({});

  useEffect(() => {
    addApiHandler(API_PATH, (response) => GetBag(response, setInfo));
    if (autoUpdate) {
      setInterval(() => {
        CallAPI(API_PATH, {}, (response) => GetBag(response, setInfo));
      }, Config.UPDATE_INTERVAL);
    }
  }, [autoUpdate]);

  return (
    <Card title="Bag">
      {info.items &&
        info.items.map((item, count) => (
          <div key={count}>
            {Object.entries(item).map(([key, value]) => (
              <p key={key}>
                <strong>{key}:</strong> {value}
              </p>
            ))}
          </div>
        ))}
    </Card>
  );
};

export { Bag };
