// Frameworks
import React from 'react';
import { connect } from 'react-redux';

// Components
import PageLayout from '../../components/PageLayout';

// Utils
import { State as GlobalState } from '../../store';
import * as API from '../../api'; 

// Styling
import './Project.scss';

interface Props {
  match: any;
  history: any;
  project?: API.Project;
}

class Project extends React.Component<Props> {
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

    return(
      <section>
        <h2>Project Details</h2>
        <p>These details are only visible to you and collaborators. All M3O projects are private.</p>

        <form>
          <div className='row'>
            <label>Name *</label>
            <input disabled required type='text' value={project.name} placeholder='My Awesome Project' name='name' />
          </div>
          
          <div className='row'>
            <label>Description</label>
            <input type='text' value={project.description} placeholder='Description' name='description' />
          </div>
        </form>
      </section>
    );
  }

  renderGithub(): JSX.Element {
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
            <input disabled type='text' value='************' name='repository' />
          </div>
          <div className='row'>
            <label>Client Secret</label>
            <input disabled type='text' value='************************' name='repository' />
          </div>
        </form>

        {/* <button className='btn warning'>Regenerate Credentials</button> */}
      </section>
    );
  }

  renderCollaborators(): JSX.Element {
    return(
      <section>
        <h2>Collaborators</h2>
        <p>Collaborators have full access to all enviroments, but only the owner (you) can invite additional collaborators or delete enviroments.</p>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Actions</th>
            </tr>
          </thead>

          <tbody>
            <tr key={'asim'}>
              <td>Asim Aslam</td>
              <td>asim@micro.mu</td>
              <td>
                <button className='danger'>Remove</button>
              </td>
            </tr>
            <tr key={'jake'}>
              <td>Jake Sanders</td>
              <td>jake@micro.mu</td>
              <td>
                <button className='danger'>Remove</button>
              </td>
            </tr>
            <tr key={'ben'}>
              <td>Ben Toogood</td>
              <td>ben@micro.mu</td>
              <td>
                <button className='danger'>Remove</button>
              </td>
            </tr>
          </tbody>
        </table>

        <p>Invite users to your project. Collaborators will recieve an email invite which is valid for 24 hours. </p>
        <form>
          <div className='row'>
            <label>Name</label>
            <input required type='text' placeholder='John Doe' name='name' />
          </div>

          <div className='row'>
            <label>Email</label>
            <input required type='email' placeholder='john@doe.com' name='email' />
          </div>

          <button className='btn'>Send Invite</button>
        </form>
      </section>
    );
  }
}

function mapStateToProps(state: GlobalState, ownProps: Props): any {
  const { project } = ownProps.match.params;

  return({
    project: state.project.projects.find(p => p.name === project),
  });
}

export default connect(mapStateToProps)(Project)