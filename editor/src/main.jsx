import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { App } from './App.jsx';
import { ConfigProvider } from 'antd';

const themes = {
  components: {
    Layout: {
      siderBg: '#fff',
      headerBg: '#fff',
    },
  },
};

createRoot(document.getElementById('root')).render(
  // <StrictMode>
  <ConfigProvider theme={themes}>
    <App />,
  </ConfigProvider>,
  // </StrictMode>,
);
