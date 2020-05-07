import { createStore, combineReducers } from 'redux';
import AccountReducer, { State as AccountState } from './Account';
import ProjectReducer, { State as ProjectState } from './Project';

export default createStore(combineReducers({
  account: AccountReducer,
  project: ProjectReducer,
}), window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__())

export interface State {
  project: ProjectState;
  account: AccountState;
}