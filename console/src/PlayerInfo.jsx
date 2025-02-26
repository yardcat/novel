import { React, useState } from 'react';
import { Card } from 'antd';


function GetPlayerInfo(info, setInfo) {
  setInfo(info);
}

const PlayerInfoCard = ({ ApiRegister }) => {
  const [info, setInfo] = useState({});
  ApiRegister["player/get_player_info"] = (response) => GetPlayerInfo(response, setInfo);

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

export default PlayerInfoCard;