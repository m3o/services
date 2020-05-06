import * as API from '../api';

// Interfaces
export interface State {
  projects: API.Project[];
}

interface Action {
  type: string;
  projectID?: string;
  environment?: API.Environment;
  project?: API.Project;
  projects?: API.Project[];
}

// Action Types
const SET_PROJECTS = 'project.set';
const CREATE_PROJECT = 'project.create';
const UPDATE_PROJECT = 'project.update';
const DELETE_PROJECT = 'project.delete';
const CREATE_ENVIRONMENT = 'project.environment.create';
const UPDATE_ENVIRONMENT = 'project.environment.update';
const DELETE_ENVIRONMENT = 'project.environment.delete';

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

export function createEnvironment(projectID: string, environment: API.Environment): Action {
  return { type: CREATE_ENVIRONMENT, environment, projectID };
}

export function updateEnvironment(environment: API.Environment): Action {
  return { type: UPDATE_ENVIRONMENT, environment };
}

export function deleteEnvironment(environment: API.Environment): Action {
  return { type: DELETE_ENVIRONMENT, environment };
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
    case CREATE_ENVIRONMENT: {
      return {
        ...state, projects: state.projects.map((p: API.Project): API.Project => {
          let project = { ...p };
          if(p.id! === action.projectID!) {
            project.environments = [...(project.environments || []), action.project!];
          }
          return project;
        }),
      }
    }
    case UPDATE_ENVIRONMENT: {
      return {
        ...state, projects: state.projects.map((p: API.Project): API.Project => {
          const envs = p.environments?.map(e => {
            return e.id === action.environment!.id ? { ...e, ...action.environment! } : e;
          });

          return { ...p, environments: envs };
        }),
      };
    }
    case DELETE_ENVIRONMENT: {
      return {
        ...state, projects: state.projects.map((p: API.Project): API.Project => {
          return {...p, environments: p.environments?.filter(e => e.id !== action.environment!.id)};
        }),
      };
    }
  }
  return state;
}