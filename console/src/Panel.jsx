import { Tag } from 'antd';

const Buff = ({ name, value, turn }) => {
  return (
    <>
      <Tag bordered={false} color="success">
        {name} : {value == 0 ? '' : value} ({turn})
      </Tag>
    </>
  );
};

const panelStyle = {
  border: '1px solid gray',
  borderRadius: '10px',
  padding: '0px 10px',
  margin: '5px 0px',
  width: '15vw',
};

const Panel = ({ role, info, isSelected, onClick }) => {
  let enemyStyle = {
    border: '1px solid blue',
    borderRadius: '10px',
    padding: '0px 10px',
    margin: '5px 0px',
    width: '10vw',
    backgroundColor: 'white',
    cusor: 'pointer',
    userSelect: 'none',
  };
  enemyStyle.backgroundColor = isSelected ? 'lightgreen' : 'white';

  return (
    <div style={role == 'actor' ? panelStyle : enemyStyle} onClick={onClick}>
      <p>name: {info.name}</p>
      <p>
        HP: {info.HP} / {info.maxHP}
      </p>
      <p> strength: {info.strength} </p>
      {role == 'actor' && <p> energy: {info.energy} </p>}
      {role == 'enemy' && (
        <strong>
          {info.intent.action} {info.intent.value}
        </strong>
      )}
      <p>
        Buff:
        {info.buffs &&
          info.buffs.map((buff) => (
            <Buff
              key={buff.name}
              name={buff.name}
              value={buff.value}
              turn={buff.turn}
            />
          ))}
      </p>
    </div>
  );
};

export { Panel };
