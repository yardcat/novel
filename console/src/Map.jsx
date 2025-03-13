import { React, useState, useEffect } from 'react';
import { Modal } from 'antd';
import CallAPI from "./Net";

const API_PATH = "world/get_map";

function OnMapGet(response, setInfo) {
}

const Mp = ({ addApiHandler, showMap, setShowMap }) => {
  const [info, setMap] = useState(122);

  useEffect(() => {
    addApiHandler(API_PATH, (response) => OnMapGet(response, setMap));
  }, []);

  return (
    <Modal title="Map" open={showMap}>
      <p>{info}</p>
    </Modal>
  );
};

export default Mp;