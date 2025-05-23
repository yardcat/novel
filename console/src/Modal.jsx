import { useState, useEffect } from 'react';
import { Card, Modal, Radio } from 'antd';

const CardModal = ({ modal }) => {
  const [choice, setChoice] = useState({});

  const handleTypeChange = (type) => (e) => {
    setChoice({
      ...choice,
      [type]: e.target.value,
    });
  };

  return (
    modal &&
    modal.visible && (
      <Modal
        title={modal.type}
        open={modal.visible}
        onOk={() => {
          modal.handler(choice);
          setChoice({});
        }}
        onCancel={() => {
          modal.hideCardModal();
          setChoice({});
        }}
      >
        {Object.keys(modal.choices).map((key) => {
          let group = modal.choices[key];
          if (Array.isArray(group) && group.length > 0) {
            return (
              <div key={key}>
                <div>{key}</div>
                <Radio.Group onChange={handleTypeChange(key)} key={key}>
                  {group.map((item, index) => (
                    <Radio key={`${key}-${index}`} value={item}>
                      {item}
                    </Radio>
                  ))}
                </Radio.Group>
              </div>
            );
          }
          return null;
        })}
      </Modal>
    )
  );
};

export { CardModal };
