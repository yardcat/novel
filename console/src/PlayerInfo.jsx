import { React, useState, useEffect } from 'react';
import { Card } from 'antd';
import { Config } from './Config';
import { CallAPI } from './Net';

const API_PATH = 'player/get_player_info';

function GetPlayerInfo(info, setInfo) {
  setInfo(info);
}

const PlayerInfo = ({ addApiHandler }) => {
  const [info, setInfo] = useState({});

  useEffect(() => {
    addApiHandler(API_PATH, (response) => GetPlayerInfo(response, setInfo));
  }, []);

  return (
    <Card title="Player Info">
      {Object.entries(info).map(([key, value]) => (
        <p key={key}>
          <strong>{key}:</strong> {value}
        </p>
      ))}
    </Card>
  );
};

export { PlayerInfo };
