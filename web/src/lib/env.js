import axios from 'axios';

let _env = process.env;

export const loadEnv = async () => {
  if (!_env.REACT_APP_API_URL) {
    const res = await axios.get(`/env.json?v=${new Date().getTime()}`);
    _env = res.data;
  }
};

const getEnv = () => {
  return {
    API_ADDRESS: _env.REACT_APP_API_URL,
    WS_ADDRESS: _env.REACT_APP_WS_URL,
    NODE_URL: _env.REACT_APP_NODE_URL,
    HYDRO_PROXY_ADDRESS: _env.REACT_APP_HYDRO_PROXY_ADDRESS,
    HYDRO_TOKEN_ADDRESS: _env.REACT_APP_HYDRO_TOKEN_ADDRESS
  };
};

export default getEnv;
