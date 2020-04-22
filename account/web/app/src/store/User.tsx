import { User } from "../api";

const SET_USER = 'SET_USET';

interface Action {
  type: string;
  user?: User;
}

interface State {
  user?: User;
}

export function setUser(user: User): Action {
  return { type: SET_USER, user };
}

const defaultState: State = {};
export default function(state = defaultState, action: Action): State {
  switch (action.type) {
    case SET_USER: 
      return { ...state, user: action.user! };
    default:
      return state;
  }
}