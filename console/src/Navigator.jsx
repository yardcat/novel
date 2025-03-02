import { useState } from 'react';
import { Cascader, Flex, InputNumber, Button, List } from 'antd';
import CallAPI from './Net';

class Node {
  value;
  label;
  children;
  constructor(value, label, children) {
    this.value = value;
    this.label = label;
    this.children = children;
  }
}

function value2Label(input) {
  const words = input.split('_');
  const capitalizedWords = words.map(word => {
      return word.charAt(0).toUpperCase() + word.slice(1);
  });
  return capitalizedWords.join('');
}

function helper(node, keys, i) {
  if (i == keys.length) {
    return;
  }
  var find = -1;
  for (var j = 0; j < node.children.length; j++) {
    if (node.children[j].value == keys[i]) {
      find = j;
      break;
    }
  }
  if (find == -1) {
    node.children.push(new Node(keys[i], value2Label(keys[i]), []));
    helper(node.children[node.children.length - 1], keys, i + 1);
  } else {
    helper(node.children[find], keys, i + 1);
  }
}

function generateOptions(api_list) {
  var root = new Node('root', 'root', []);
  for (var i = 0; i < api_list.length; i++) {
    const keys = api_list[i].split('/');
    helper(root, keys, 0);
  }
  return root.children;
}

const CollectableComponent = ({ visible, data, onSubmit }) => {
  if (!visible) return null;

  return (
    <div id="collectable">
      <List
        itemLayout="horizontal"
        dataSource={data}
        renderItem={(item, index) => (
          <List.Item>
            <span>{item}</span>
            <InputNumber
              type="number"
              defaultValue={0}
            />
          </List.Item>
        )}
      />
      <Button type="primary" onClick={() => onSubmit(data)}>提交</Button>
    </div>
  );
};

const Navigator = ({ apiHandlers }) => {
  const options = generateOptions(Object.keys(apiHandlers));
  const [collectableVisible, setCollectableVisible] = useState(false);
  const [collectableData, setCollectableData] = useState(["a", "b", "c"]);

  const handleChange = (value, handlers) => {
    var path = value.join('/');
    if (value[value.length - 1] === 'collect') {
      setCollectableVisible(true);
      // const labels = value.slice(0, -1).map(value2Label);
      // setCollectableData(labels.map(label => ({ label, value: 0 })));
    } else {
      setCollectableVisible(false);
      CallAPI(path, {}, handlers[path]);
    }
  };

  const handleCollectableSubmit = (data) => {
    setCollectableVisible(false);
    console.log('Selected data:', data);
  };

  return (
    <Flex gap="small" align="flex-start">
      <Cascader.Panel options={options} onChange={(value) => handleChange(value, apiHandlers)} />
      <CollectableComponent
        visible={collectableVisible}
        data={collectableData}
        onSubmit={handleCollectableSubmit}
      />
    </Flex>
  );
};

export default Navigator;