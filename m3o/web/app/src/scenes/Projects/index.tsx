import React from 'react';
import { connect } from 'react-redux';
import PageLayout from '../../components/PageLayout';
import * as API from '../../api';
import { State as GlobalState } from '../../store';
import { deleteProject } from '../../store/Project';
// import AddUser from './assets/add-user.png';

interface Props {
  history: any;
  projects: API.Project[];
  deleteProject: (user: API.Project) => void;
}

class ProjectsScene extends React.Component<Props> {
  render(): JSX.Element {
    return(
      <PageLayout className='Projects'>
        <header>
          <h1>Projects</h1>
          
          <button className='btn' onClick={() => this.props.history.push('/projects/new')}>
            <p>Create Project</p>
          </button>
        </header>

        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Actions</th>
            </tr>
          </thead>

          <tbody>
            { this.props.projects.map(p => <tr key={p.id}>
              <td>{p.name}</td>
              <td>
                <button className='warning' onClick={() => this.editProject(p)}>Edit</button>
                <button className='danger' onClick={() => this.deleteProject(p)}>Delete</button>
              </td>
            </tr>) }
          </tbody>
        </table>
      </PageLayout>
    )
  }

  editProject(project: API.Project): void {
    this.props.history.push(`/projects/${project.id}/edit`);
  }

  deleteProject(project: API.Project): void {
    // eslint-disable-next-line no-restricted-globals
    if (!confirm(`Are you sure you want to delete ${project.name}?`)) return;
    this.props.deleteProject(project);
  }
}

function mapStateToProps(state: GlobalState): any {
  return({
    projects: state.project.projects.sort(sortByName),
  });
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    deleteProject: (project: API.Project) => dispatch(deleteProject(project)),
  })
}

export default connect(mapStateToProps, mapDispatchToProps)(ProjectsScene);

function sortByName(a: API.Project, b: API.Project): number {
  const aName = a.name.toUpperCase();
  const bName = b.name.toUpperCase();
  if(aName > bName) return 1;
  if(aName < bName) return -1;
  return 0;
}