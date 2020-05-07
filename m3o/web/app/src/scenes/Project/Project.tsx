// Frameworks
import React from 'react';
import { connect } from 'react-redux';

// Components
import PageLayout from '../../components/PageLayout';
import ValidatedInput from '../../components/ValidatedInput';

// Utils
import { State as GlobalState } from '../../store';
import * as API from '../../api'; 

// Styling
import './Project.scss';
import RefreshIcon from './assets/refresh.png';
import { updateProject } from '../../store/Project';

interface Props {
  match: any;
  history: any;
  
  user: API.User;
  project?: API.Project;
  updateProject: (project: API.Project) => void;
}

interface State {
  clientID?: string;
  clientSecret?: string;
  credsLoading: boolean;

  inviteName?: string;
  inviteEmail?: string;
  inviteLoading: boolean;
}

class Project extends React.Component<Props, State> {
  readonly state: State = { credsLoading: false, inviteLoading: false };

  render(): JSX.Element {
    const { project } = this.props;
    
    if(!project) {
      // this.props.history.push('/not-found');
      return null
    }

    return <PageLayout className='Project'>
      <div className='center'>
        <div className='header'>
          <h1>{project.name}</h1>
        </div>

        { project.environments ? null : this.renderFirstEnv() }
        { this.renderDetails() }
        { this.renderGithub() }
        { this.renderCollaborators() }
     </div>
    </PageLayout>
  }

  renderFirstEnv(): JSX.Element {
    const onClick = () => this.props.history.push(`/new/environment/${this.props.project.name}`);

    return(
      <div onClick={onClick.bind(this)} className='first-env'>
        <h5>Create your first enviroment</h5>
        <p>You don't have any enviroments setup for {this.props.project.name}. Click here to create your first one.</p>
      </div>
    );
  }

  renderDetails(): JSX.Element {
    const { project } = this.props;

    const onChange = (key: string, value: string): void => {
      const env = { ...this.props.project, [key]:value };
      this.props.updateProject(env);
    }

    const onSave = (): Promise<void> => {
      return new Promise((resolve: Function, reject: Function) => {
        API.Call("Projects/UpdateProject", { project: this.props.project })
          .then(() => resolve())
          .catch(err => reject(err.response ? err.response.data.detail : err.message));
      });
    }

    return(
      <section>
        <h2>Project Details</h2>
        <p>These details are only visible to you and collaborators. All M3O projects are private.</p>
        <form>
          <div className='row'>
            <label>Name *</label>
            <ValidatedInput disabled value={project?.name} />
          </div>
          
          <div className='row'>
            <label>Description</label>
            <ValidatedInput
              name='description'
              validate={onSave} 
              onChange={onChange}
              validateDelay={1000}
              placeholder='Description'
              value={project?.description || ''} />
          </div>
        </form>
      </section>
    );
  }

  renderGithub(): JSX.Element {
    const { credsLoading, clientID, clientSecret } = this.state;

    const refreshCreds = () => {
      if(this.state.credsLoading) return;
      this.setState({ credsLoading: true });

      API.Call("Projects/WebhookAPIKey", { project_id: this.props.project.id })
        .then((res) => {
          this.setState({
            credsLoading: false,
            clientID: res.data.client_id,
            clientSecret: res.data.client_secret
          });
        })
        .catch((err) => {
          alert(err.response ? err.response.data.detail : err.message);
          this.setState({ credsLoading: false });
        });
    }

    return(
      <section>
        <h2>GitHub</h2>
        <p>M3O connects to GitHub and builds your services in your repo, keeping your source and builds firmly in your control. The <a href='https://github.com/micro/actions' target='blank'>micro/actions</a> GitHub action automatically builds your services when any changes are detected and triggers a release. Find our more at our <a href='/todo'>docs</a>.</p>

        <form>
          <div className='row'>
            <label>Repository</label>
            <input disabled type='text' value={this.props.project.repository} name='repository' />
          </div>
          <div className='row'>
            <label>Client ID</label>

            <div className='refresh-input'>
              <input
                disabled
                type='text'
                value={clientID}
                placeholder='************' />

              <img
                src={RefreshIcon}
                onClick={refreshCreds}
                alt='Refresh Credentials'
                className={credsLoading ? 'loading' : ''} />
            </div>
          </div>
          <div className='row'>
            <label>Client Secret</label>

            <div className='refresh-input'>
              <input
                disabled
                type='text'
                value={clientSecret}
                placeholder='************************' />

              <img
                src={RefreshIcon}
                onClick={refreshCreds}
                alt='Refresh Credentials'
                className={credsLoading ? 'loading' : ''} />
            </div>
          </div>
        </form>
      </section>
    );
  }

  renderCollaborators(): JSX.Element {
    // only render collaborators if the current user is the owner
    // of the project.
    var isOwner: Boolean;
    this.props.project?.members?.forEach(u => {
      if(u.role.toLowerCase() !== 'owner') return;
      if(u.id !== this.props.user.id) return;
      isOwner = true;
    });
    if(!isOwner) return;

    const onChange = (e: any) => {
      if(e.target.name === 'name') {
        this.setState({ inviteName: e.target.value });
      } else {
        this.setState({ inviteEmail: e.target.value });
      }
    }

    const onSubmit = () => {
      const { inviteName, inviteEmail, inviteLoading } = this.state;
      if(!inviteEmail || !inviteName || inviteLoading) return;
      if(inviteName.length === 0) return;
      if(inviteEmail.length === 0) return;

      this.setState({ inviteLoading: true });
      
      const params = {
        name: inviteName,
        email: inviteEmail,
        project_id: this.props.project.id,
      };

      API.Call("Projects/Invite", params)
        .then(() => {
          alert("Invite sent to " + inviteName);
          this.setState({ inviteLoading: false, inviteEmail: '', inviteName: '' });
        })
        .catch(err => {
          alert(err.response ? err.response.data.detail : err.message);
          this.setState({ inviteLoading: false });
        });
    }

    const { inviteName, inviteEmail, inviteLoading } = this.state;

    return(
      <section>
        <h2>Collaborators</h2>
        <p>Collaborators have full access to all enviroments, but only the owner (you) can invite additional collaborators.</p>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Role</th>
              <th>Actions</th>
            </tr>
          </thead>

          <tbody>
            { this.props.project?.members?.map(u => {
              const currentUser = u.id === this.props.user.id;

              return(
                <tr key={'asim'}>
                  <td>{u.first_name} {u.last_name} {currentUser ? '(me)' : ''}</td>
                  <td>{u.email}</td>
                  <td>{u.role}</td>
                  <td>
                    <button className='danger'>Remove</button>
                  </td>
                </tr>
              );
            })}

            <tr>
              <td>
                <input required value={inviteName} onChange={onChange} type='text' placeholder='John Doe' name='name' />
              </td>

              <td>
                <input required value={inviteEmail} onChange={onChange} type='email' placeholder='john@doe.com' name='email' />
              </td>

              <td>Collaborator</td>
              <td>
                <button onClick={onSubmit} disabled={inviteLoading}>Invite</button>
              </td>
            </tr>
          </tbody>
        </table>
      </section>
    );
  }
}

function mapStateToProps(state: GlobalState, ownProps: Props): any {
  const { project } = ownProps.match.params;

  return({
    user: state.account.user,
    project: state.project.projects.find(p => p.name === project),
  });
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    updateProject: (project: API.Project) => dispatch(updateProject(project)),
  });
}

export default connect(mapStateToProps, mapDispatchToProps)(Project)