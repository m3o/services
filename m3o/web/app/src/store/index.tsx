import { createStore, combineReducers } from 'redux';
import TeamReducer, { State as TeamState } from './Team';
import ConfigurationReducer, { State as ConfigurationState } from './Configuration';

export default createStore(combineReducers({
  team: TeamReducer,
  configuration: ConfigurationReducer,
}))

export interface State {
  team: TeamState;
  configuration: ConfigurationState;
}