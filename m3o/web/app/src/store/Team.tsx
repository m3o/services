import * as API from '../api';

// Interfaces
export interface State {
  users: API.User[];
}

interface Action {
  type: string;
  user?: API.User;
}

// Action Types
const UPDATE_USER = 'team.user.update';
const DELETE_USER = 'team.user.delete';

// Actions
export function updateUser(user: API.User): Action {
  return { type: UPDATE_USER, user: user };
}

export function deleteUser(user: API.User): Action {
  return { type: DELETE_USER, user: user };
}

const defaultState: State = {
  users: [
    {id: "Asim", firstName: "Asim", lastName: "Aslam", email: "asim@micro.mu", roles: ["Admin", "Developer"], me: true},
    {id: "Jake", firstName: "Jake", lastName: "Sanders", email: "jake@micro.mu", roles: ["Developer"]},
    {id: "Ben", firstName: "Ben", lastName: "Toogood", email: "ben@micro.mu", roles: ["Developer"]},
    {id: "Janos", firstName: "Janos", lastName: "Dobronszki", email: "janos@micro.mu", roles: ["Developer"]},
    {id: "Vasiliy", firstName: "Vasiliy", lastName: "Tolstov", email: "vasiliy@micro.mu", roles: ["Developer"]},
  ],
};

// Reducer
export default function(state = defaultState, action: Action): State {
  switch(action.type) {
    case UPDATE_USER: {
      return {
        ...state, users: [
          ...state.users.filter(u => u.id !== action.user!.id), action.user,
        ],
      };
    }
    case DELETE_USER: {
      return {
        ...state, users: [
          ...state.users.filter(u => u.id !== action.user!.id),
        ],
      };
    }
  }
  return state;
}