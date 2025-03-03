import { React, useState, useEffect } from 'react';
import Navigator from './Navigator';
import PlayerInfo from './PlayerInfo';
import Bag from './Bag';
import Action from './Action';
import {initConfig} from './Config';

const App = () => {
  const [apiHandlers, setApiHandlers] = useState({});

  const addApiHandler = (path, handler) => {
    setApiHandlers(prevHandlers => ({
      ...prevHandlers,
      [path]: handler
    }));
  };

  useEffect(() => {
    initConfig();
  }, []);

  return (
    <div className="App">
      <Navigator apiHandlers={apiHandlers}></Navigator>
      <PlayerInfo addApiHandler={addApiHandler} autoUpdate={true}></PlayerInfo>
      <Bag addApiHandler={addApiHandler} autoUpdate={true}></Bag>
      <Action addApiHandler={addApiHandler}></Action>
    </div>
  );
};

export default App;