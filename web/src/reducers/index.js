import { combineReducers } from 'redux';
import { reducer as formReducer } from 'redux-form';
import market from './market';
import account from './account';
import config from './config';
import { WalletReducer } from '@gongddex/hydro-sdk-wallet';

const rootReducer = combineReducers({
  market,
  account,
  config,
  form: !!formReducer ? formReducer : {},
  WalletReducer
});
export default rootReducer;
