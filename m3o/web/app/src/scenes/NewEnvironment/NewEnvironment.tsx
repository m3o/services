// Frameworks
import React from 'react';
import { connect } from 'react-redux';

// Utils
import { State as GlobalState } from '../../store';
import * as API from '../../api';

// Components
import PageLayout from '../../components/PageLayout';

// Styling
import './NewEnvironment.scss';
import { createEnvironment } from '../../store/Project';

interface Props {
  match: any;
  history: any;
  project?: API.Project;
  createEnvironment: (projectID: string, env: API.Environment) => void;
}

interface State {
  environment: API.Environment;
  error?: string;
  loading: boolean;
}

// regex to check for specical chars
var regex = /[^\w]|_/g

class NewEnvironment extends React.Component<Props, State> {
  readonly state: State = { environment: { name: '', description: '' }, loading: false };

  onInputChange(e: any) {
    // force the name to lowercase
    if(e.target.name === 'name') e.target.value = e.target.value.toLowerCase();

    // construct the new environment
    const env: API.Environment = { ...this.state.environment, [e.target.name]: e.target.value };

    // check for errors
    let error: string = undefined;
    if(env.name.length > 0 && regex.test(env.name)) {
      error = "Name cannot contain any special characters, must be URL safe.";
    } else {
      (this.props.project.environments || []).forEach(e => {
        if(e.name.toLowerCase() === env.name.toLowerCase()) {
          error = `${this.props.project.name}/${env.name} is already taken`;
        }
      });
    }

    // update the state
    this.setState({ error, environment: env });
  }

  onSubmit(e?: any): void {
    if(e) e.preventDefault();
    if(this.state.loading || this.state.error) return;

    const { name, description } = this.state.environment;
    const { project } = this.props;

    const params = {
      project_id: project.id,
      environment: { name: name, description: description },
    };
    
    API.Call("Projects/CreateEnvironment", params)
      .then((res) => {
        this.props.createEnvironment(project.id, res.data.environment);
        this.props.history.push(`/projects/${project.name}/${name}`);
      })
      .catch((err) => {
        this.setState({ error: (err.response ? err.response.data.detail : err.message) });
      })
  }

  render(): JSX.Element {
    const { project } = this.props;
    if(!project) return null;

    const { environment, error, loading } = this.state;
    const disabled = !!error || loading || environment.name.length === 0;

    return (
      <PageLayout className='NewEnvironment'>
        <div className='center'>
          <div className='header'>
            <h1>{project.name} / New Environment</h1>
          </div>

          <section>
            <h2>Environment Details</h2>
            <p>Set the name and description for your environment. You cannot change name once it is set.</p>

            <form onSubmit={this.onSubmit.bind(this)}>
              <div className='row'>
                <label>Name *</label>
                <input required type='text' value={environment.name} placeholder='production' name='name' onChange={this.onInputChange.bind(this)} />
              </div>
              
              <div className='row'>
                <label>Description</label>
                <input type='text' value={environment.description} placeholder={`The ${project.name} production environment`} name='description'  onChange={this.onInputChange.bind(this)} />
              </div>

              { error ? <p className='error'>{error}</p> : null }

              <button onClick={this.onSubmit.bind(this)} disabled={disabled} className='btn'>Create Environment</button>
            </form>
          </section>
        </div>
      </PageLayout>
    );
  }
}

function mapStateToProps(state: GlobalState, ownProps: Props): any {
  const { project } = ownProps.match.params;

  return({
    project: state.project.projects.find(p => p.name === project),
  });
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    createEnvironment: (projectID: string, env: API.Environment) => dispatch(createEnvironment(projectID, env)),
  });
}

export default connect(mapStateToProps, mapDispatchToProps)(NewEnvironment)