import * as API from '../api';

// Interfaces
export interface State {
  projects: API.Project[];
}

interface Action {
  type: string;
  project?: API.Project;
  projects?: API.Project[];
}

// Action Types
const SET_PROJECTS = 'project.set';
const CREATE_PROJECT = 'project.create';
const UPDATE_PROJECT = 'project.update';
const DELETE_PROJECT = 'project.delete';

// Actions
export function setProjects(projects: API.Project[]): Action {
  return { type: SET_PROJECTS, projects };
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

// Reducer
const defaultState: State = { projects: [] };
export default function(state = defaultState, action: Action): State {
  switch(action.type) {
    case SET_PROJECTS: {
      return { ...state, projects: action.projects! };
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