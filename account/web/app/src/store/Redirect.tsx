const SET_REDIRECT = 'SET_REDIRECT';

interface Action {
  type: string;
  path: string;
}

interface State {
  path?: string;
}

export function setRedirect(path: string): Action {
  return { type: SET_REDIRECT, path };
}

const defaultState: State = {};
export default function(state = defaultState, action: Action): State {
  switch (action.type) {
    case SET_REDIRECT: 
      return { ...state, path: action.path! };
    default:
      return state;
  }
}