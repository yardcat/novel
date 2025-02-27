import Navigator from './Navigator';
import PlayerInfo from './PlayerInfo';
import Bag from './Bag';
import CallAPI from './Net';

const apiHandlers = {};

const App = () => (
  <div className="App">
    <PlayerInfo ApiRegister={apiHandlers}></PlayerInfo>
    <Bag ApiRegister={apiHandlers}></Bag>
    <Navigator ApiRegister={apiHandlers}></Navigator>
  </div>
);

export default App;