import { createStore, combineReducers } from 'redux';
import TeamReducer, { State as TeamState } from './Team';
import AccountReducer, { State as AccountState } from './Account';
import ProjectReducer, { State as ProjectState } from './Project';
import ConfigurationReducer, { State as ConfigurationState } from './Configuration';

export default createStore(combineReducers({
  team: TeamReducer,
  account: AccountReducer,
  project: ProjectReducer,
  configuration: ConfigurationReducer,
}), window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__())

export interface State {
  team: TeamState;
  project: ProjectState;
  account: AccountState;
  configuration: ConfigurationState;
}