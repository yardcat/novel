import Navigator from './Navigator';
import PlayerInfo from './PlayerInfo';
import CallAPI from './Net';

const apiHandlers = {};

const App = () => (
  <div className="App">
    <PlayerInfo ApiRegister={apiHandlers}></PlayerInfo>
    <Navigator ApiRegister={apiHandlers}></Navigator>
  </div>
);

export default App;