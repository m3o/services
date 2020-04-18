import { createStore, combineReducers } from 'redux';
import TeamReducer, { State as TeamState } from './Team';

export default createStore(combineReducers({
  team: TeamReducer,
}))

export interface State {
  team: TeamState;
}