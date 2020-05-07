// Frameworks
import React from 'react';
import { connect } from 'react-redux';

// Utils
import { State as GlobalState } from '../../store';
import { deleteEnvironment } from '../../store/Project';
import * as API from '../../api';

// Components
import PageLayout from '../../components/PageLayout';

// Styling
import EditIcon from './assets/edit.png';
import './Enviroment.scss';

interface Props {
  match: any;
  history: any;
  project?: API.Project;
  environment?: API.Environment;
  deleteEnvironment: (env: API.Environment) => void;
}

class Enviroment extends React.Component<Props> {
  render(): JSX.Element {
    const { project, environment } = this.props;
    const domain = `https://${environment?.namespace}.m3o.app`; 

    return <PageLayout className='Enviroment'>
      <div className='center'>
        <div className='header'>
          <h1><span>{project?.name}</span> / <span>{environment?.name}</span></h1>
          <img src={EditIcon} alt='Edit Name' />

          <button className='btn'><p>Launch Dashboard</p></button>
        </div>

        <section>
          <h2>Enviroment Details</h2>
          <p>These details are only visible to you and collaborators. All M3O projects are private.</p>

          <form>
            <div className='row'>
              <label>Name *</label>
              <input disabled required type='text' value={environment?.name} placeholder='My Awesome Project' name='name' />
            </div>
            
            <div className='row'>
              <label>Description</label>
              <input type='text' value={environment?.description} placeholder='Description' name='description' />
            </div>
          </form>
        </section>

        <section>
          <h2>DNS</h2>
          <p>Your default domain is <a href={domain} target='blank'>{domain}</a>. Your web domain is served at <a href={domain + '/web'}>/web</a> and your API is available at <a href={domain + '/api'}>/api</a>. To configure a custom domain, enter the domains below and then setup CNAME records for each domain pointing at <strong>m3o.app</strong>. For more information about custom domains, see <a href='/todo'>the docs</a>.</p>
          <form>
            <div className='row'>
              <label>Web Domain</label>
              <input disabled type='text' value='' placeholder='Coming Soon...' name='web_domain' />
            </div>
            
            <div className='row'>
              <label>API Domain</label>
              <input disabled type='text' value='' placeholder='Coming Soon...' name='api_domain' />
            </div>
          </form>
        </section>

        <section>
          <h2>CLI</h2>
          <p>Configure your CLI to use the {project?.name}/{environment?.name} enviroment. Firstly, all calls made to your enviroment are authenticated so if you aren't already, login using the following command and a token you can get at <a href='/todo'>this link</a>.</p>
          <p className='code'>
            micro login [token]
          </p>

          <p>Once you're logged in, add your enviroment and configure micro to use it with the following commands.</p>
          <p className='code'>
            micro env add {project?.name}/{environment?.name} {environment?.namespace}.proxy.m3o.app
            <br />
            micro env set {project?.name}/{environment?.name}
          </p>
        </section>

        <section>
          <h2>Settings</h2>
          <p><strong>Warning:</strong> Deleting your enviroment cannot be undone and all data will be lost.</p>
          <button onClick={this.onDeleteClicked.bind(this)} className='btn danger'>Delete {project?.name}/{environment?.name}</button>
        </section>
     </div>
    </PageLayout>
  }

  onDeleteClicked(): void {
    // eslint-disable-next-line
    if(!confirm("Are you sure you want to delete this environment?")) return

    API.Call("Projects/DeleteEnvironment", { id: this.props.environment.id })
    .then((res) => {
      this.props.deleteEnvironment(this.props.environment);
      this.props.history.push(`/projects/${this.props.project.name}`);
    })
    .catch(err => alert(err.response ? err.response.data.detail : err.message));
  }
}

function mapStateToProps(state: GlobalState, ownProps: Props): any {
  const { params } = ownProps.match;
  const project: API.Project = state.project.projects.find(p => p.name === params.project);
  const environment: API.Environment = project?.environments?.find(e => e.name === params.environment);

  return({ environment, project });
}

function mapDispatchToProps(dispatch: Function): any {
  return ({
    deleteEnvironment: (env: API.Environment) => dispatch(deleteEnvironment(env)),
  });
}

export default connect(mapStateToProps, mapDispatchToProps)(Enviroment);