import React, { createRef } from 'react';
import Popup from "reactjs-popup";
import './style.scss';

interface Props {

}

interface State {
  scene: 'list' | 'create';
  project: Project;
  projects: Project[];
}

interface Project {
  name: string;
  namespace: string;
}

export default class ProjectSwitcher extends React.Component<Props, State> {
  popup: React.RefObject<any> = createRef();

  readonly state: State = {
    scene: 'create',
    projects: [
      { name: "Kytra", namespace: "kytra-production" },
      { name: "Kytra / Staging", namespace: "kytra-staging" },
      { name: "Kytra / Ben Development", namespace: "kytra-development-ben" },
    ],
    project: {
      name: '',
      namespace: ''
    },
  };

  render(): JSX.Element {
    const button = (
      <div className='project'>
        <p>Kytra</p>
        <p className='descriptor'>Current Project</p>
    </div>
    );

    let inner: JSX.Element;
    switch(this.state.scene) {
      case 'list': {
        inner = this.renderList();
        break;
      }
      case 'create': {
        inner = this.renderCreate();
        break;
      }
    }

    return(<Popup ref={this.popup} trigger={button} modal={true} defaultOpen={true} onClose={() => this.setState({ scene: 'list' })}>
      { inner }
    </Popup>);
  }

  renderCreate(): JSX.Element {
    return(
      <div className='CreateProject'>
        <div className='upper'>
          <h3>Create Project</h3>

          <button className='btn btn-small danger' onClick={() => this.setState({ scene: 'list' })}>
            <p>Cancel</p>
          </button>

          <button className='btn btn-small' onClick={() => this.setState({ scene: 'list' })}>
            <p>Save</p>
          </button>
        </div>

        <form>
          <label>ID *</label>
          <input
            type='text'
            name='namespace'
            placeholder='my-first-project-839292'
            value={this.state.project.namespace}
            onChange={this.onFormChange.bind(this)} />

          <label>Name *</label>
          <input
            type='text'
            name='name'
            placeholder='My first project'
            value={this.state.project.name}
            onChange={this.onFormChange.bind(this)} />
        </form>
      </div>
    )
  }

  onFormChange(e): void {
    this.setState({
      project: {
        ...this.state.project,
        [e.target.name]: e.target.value,
      },
    });
  }

  renderList(): JSX.Element {
    return(
      <div className='ListProjects'>
        <div className='upper'>
          <h3>Projects</h3>

          <button className='btn btn-small' onClick={() => this.setState({ scene: 'create' })}>
            <p>Create Project</p>
          </button>
        </div>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>ID</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            { this.state.projects.map(p => <tr>
              <td>{p.name}</td>
              <td>{p.namespace}</td>
              <td>
                <button onClick={() => this.popup.current.closePopup()}>Switch</button>
              </td>
            </tr>)}
          </tbody>
        </table>
      </div>
    );
  }
}