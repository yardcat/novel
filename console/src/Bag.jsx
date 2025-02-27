import { React, useState } from 'react';
import { Card } from 'antd';


function GetBag(info, setInfo) {
  setInfo(info);
}

const Bag = ({ ApiRegister }) => {
  const [info, setInfo] = useState({});
  ApiRegister["player/get_bag"] = (response) => GetBag(response, setInfo);

  return (
      <Card title="Bag">
          {info.items && info.items.map((item, count) => (
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

export default Bag;