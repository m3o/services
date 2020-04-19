import * as API from '../api';

// Interfaces
export interface State {
  envVars: API.EnvVar[];
}

interface Action {
  type: string;
  envVar?: API.EnvVar;
}

// Action Types
const ADD_ENV_VAR = 'configuration.envvar.add';
const UPDATE_ENV_VAR = 'configuration.envvar.update';
const DELETE_ENV_VAR = 'configuration.envvar.delete';

// Actions
export function addEnvVar(ev: API.EnvVar): Action {
  return { type: ADD_ENV_VAR, envVar: ev };
}

export function updateEnvVar(ev: API.EnvVar): Action {
  return { type: UPDATE_ENV_VAR, envVar: ev };
}

export function deleteEnvVar(ev: API.EnvVar): Action {
  return { type: DELETE_ENV_VAR, envVar: ev };
}

const defaultState: State = {
  envVars: [
    { id: '0', service: 'go.micro.service.payments', key: 'STRIPE_API_KEY', value: 'FOOBARFOOBARFOOBARFOOBARFOOBAR', secret: true },
    { id: '1', service: '*', key: 'MICRO_LOG_LEVEL', value: 'info', secret: false },
  ],
};

// Reducer
export default function(state = defaultState, action: Action): State {
  switch(action.type) {
    case ADD_ENV_VAR: {
      return {
        ...state, envVars: [
          ...state.envVars, {...action.envVar, id: state.envVars.length.toString()},
        ],
      };
    }
    case UPDATE_ENV_VAR: {
      return {
        ...state, envVars: [
          ...state.envVars.filter(v => v.id !== action.envVar.id), action.envVar,
        ],
      };
    }
    case DELETE_ENV_VAR: {
      return {
        ...state, envVars: [
          ...state.envVars.filter(v => v.id !== action.envVar.id),
        ],
      };
    }
  }
  return state;
}