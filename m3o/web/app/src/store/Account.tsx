import * as API from '../api';

// Interfaces
export interface State {
  user?: API.User;
}

interface Action {
  type: string;
  user?: API.User;
}

// Action Types
const SET_USER = 'account.user.set';

// Actions
export function setUser(user: API.User): Action {
  return { type: SET_USER, user };
}

const defaultState: State = {};

// Reducer
export default function(state = defaultState, action: Action): State {
  switch(action.type) {
    case SET_USER: {
      return { ...state, user: action.user! };
    }
  }
  return state;
}