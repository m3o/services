import React, { createRef } from 'react';
import Popup from "reactjs-popup";
import { connect } from 'react-redux';
import { State as GlobalState } from '../../../../store';
import { createProject, switchProject } from '../../../../store/Project';
import * as API from '../../../../api';
import './style.scss';

interface Props {
  projects: API.Project[];
  currentProject?: API.Project;
  createProject: (p: API.Project) => void;
  switchProject: (p: API.Project) => void;
}

interface State {
  scene: 'list' | 'create';
  project: API.Project;
}

class ProjectSwitcher extends React.Component<Props, State> {
  popup: React.RefObject<any> = createRef();

  readonly state: State = {
    scene: 'list',
    project: {
      name: '',
      namespace: ''
    },
  };

  render(): JSX.Element {
    const button = (
      <div className='project'>
        <p>{this.props.currentProject?.name}</p>
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

    return(<Popup ref={this.popup} trigger={button} modal={true} onClose={() => this.setState({ scene: 'list' })}>
      { inner }
    </Popup>);
  }

  renderCreate(): JSX.Element {
    const cancel = () => {
      this.setState({ scene: 'list', project: { name: '', namespace: '' }});
    }

    const save = () => {
      // temp hack to set ID whilst not connected to api
      this.setState({ project: {...this.state.project, id: this.props.projects.length.toString() }});

      this.props.createProject(this.state.project);
      this.setState({ scene: 'list', project: { name: '', namespace: '' }});
      
      // not sure if we want to auto-select and close popup at this point?
      // this.props.switchProject(this.state.project);
      // this.popup.current.closePopup();
    }

    return(
      <div className='CreateProject'>
        <div className='upper'>
          <h3>Create Project</h3>

          <button className='btn btn-small danger' onClick={cancel.bind(this)}>
            <p>Cancel</p>
          </button>

          <button className='btn btn-small' onClick={save.bind(this)}>
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
            { this.props.projects.map(p => <tr key={p.id!}>
              <td>{p.name}</td>
              <td>{p.namespace}</td>
              <td>
                <button onClick={() => this.switchTo(p)}>Switch</button>
              </td>
            </tr>)}
          </tbody>
        </table>
      </div>
    );
  }

  switchTo(project: API.Project): void {
    this.props.switchProject(project);
    this.popup.current.closePopup();
  }
}

function mapStateToProps(state: GlobalState): any {
  const { projects, currentProjectID } = state.project;
  const currentProject = projects.find(p => p.id === currentProjectID);
  return({ projects, currentProject });
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    createProject: (p: API.Project) => dispatch(createProject(p)),
    switchProject: (p: API.Project) => dispatch(switchProject(p)),
  });
}

export default connect(mapStateToProps, mapDispatchToProps)(ProjectSwitcher);