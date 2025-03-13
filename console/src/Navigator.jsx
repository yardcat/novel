import { useEffect, useState } from 'react';
import { Cascader, Flex } from 'antd';
import CallAPI from './Net';
import Collectable from './Collectable';
import Mp from './Map';

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

const Navigator = ({ apiHandlers, addApiHandler, setAction }) => {
  const options = generateOptions(Object.keys(apiHandlers));
  const [showCollect, setShowCollect] = useState(false);
  const [showMap, setShowMap] = useState(false);

  const handleChange = (value, handlers) => {
    var path = value.join('/');
    if (value[value.length - 1] === 'collect') {
      setShowCollect(true);
    } else if (value[value.length - 1] === 'get_map') {
      setShowMap(true);
      // CallAPI(path, {}, handlers[path]);
    } else {
      setShowCollect(false);
      CallAPI(path, {}, handlers[path]);
    }
  };

  return (
    <Flex gap="small" align="flex-start">
      <Cascader.Panel options={options} onChange={(value) => handleChange(value, apiHandlers)} />
      <Collectable showCollect={showCollect} setShowCollect={setShowCollect} setAction={setAction} />
      <Mp addApiHandler={addApiHandler} showMap={showMap} setShowMap={setShowMap} />
    </Flex>
  );
};

export default Navigator;