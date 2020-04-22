import * as API from '../api';

// Interfaces
export interface State {
  projects: API.Project[];
  currentProjectID?: string;
}

interface Action {
  type: string;
  project?: API.Project;
}

// Action Types
const SWITCH_PROJECT = 'project.switch';
const CREATE_PROJECT = 'project.create';
const UPDATE_PROJECT = 'project.update';
const DELETE_PROJECT = 'project.delete';

// Actions
export function switchProject(project: API.Project): Action {
  return { type: SWITCH_PROJECT, project };
}

export function createProject(project: API.Project): Action {
  return { type: CREATE_PROJECT, project };
}

export function updateProject(project: API.Project): Action {
  return { type: UPDATE_PROJECT, project };
}

export function deleteProject(project: API.Project): Action {
  return { type: DELETE_PROJECT, project };
}

const defaultState: State = {
  projects: [
    { id: "one", name: "Kytra", namespace: "kytra-production" },
    { id: "two", name: "Kytra / Staging", namespace: "kytra-staging" },
    { id: "three", name: "Kytra / Ben Development", namespace: "kytra-development-ben" },
  ],
  currentProjectID: "one",
};

// Reducer
export default function(state = defaultState, action: Action): State {
  switch(action.type) {
    case SWITCH_PROJECT: {
      return { ...state, currentProjectID: action.project!.id };
    }
    case CREATE_PROJECT: {
      return {
        ...state, projects: [
          ...state.projects, action.project!,
        ],
      };
    }
    case UPDATE_PROJECT: {
      return {
        ...state, projects: [
          ...state.projects.filter(u => u.id !== action.project!.id), action.project!,
        ],
      };
    }
    case DELETE_PROJECT: {
      return {
        ...state, projects: [
          ...state.projects.filter(u => u.id !== action.project!.id),
        ],
      };
    }
  }
  return state;
}