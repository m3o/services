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
import ValidatedInput from '../../components/ValidatedInput';

interface Props {
  match: any;
  history: any;
  project?: API.Project;
  createEnvironment: (projectID: string, env: API.Environment) => void;
}

interface State {
  environment: API.Environment;
  loading: boolean;
  nameValid: boolean;
}

// regex to check for specical chars
var regex = /[^\w]|_/g

class NewEnvironment extends React.Component<Props, State> {
  readonly state: State = { 
    loading: false,
    nameValid: false,
    environment: { name: '', description: '' },
  };


  onSubmit(e?: any): void {
    if(e) e.preventDefault();
    if(this.state.loading || !this.state.nameValid) return;

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
      .catch((err) => alert(err.response ? err.response.data.detail : err.message));
  }

  render(): JSX.Element {
    const { project } = this.props;
    if(!project) return null;

    const { environment, loading, nameValid } = this.state

    const validateName = async (name: string): Promise<string> => {      
      return new Promise(async (resolve: Function, reject: Function) => {
        if(name.length < 3) {
          reject("Name must be at least 3 characters long");
          return
        }

        if(regex.test(name)) {
          reject("Name cannot contain any special characters");
          return;
        }

        API.Call('Projects/ValidateEnvironmentName', { name, project_id: project.id })
          .then(() => resolve())
          .catch(err => reject(err.response ? err.response.data.detail : err.message));
      });
    }

    const setNameValid = () => this.setState({ nameValid: true });
    const setNameInvalid = () => this.setState({ nameValid: false });
    
    const onChange = (key: string, value: string) => {
      if(key === 'name') value = value.toLowerCase();
      this.setState({ environment: { ...this.state.environment, [key]: value } });
    };

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

                <ValidatedInput
                  name='name'
                  onChange={onChange}
                  value={environment.name}
                  placeholder='production'
                  validate={validateName}
                  onValid={setNameValid}
                  onInvalid={setNameInvalid} />
              </div>
              
              <div className='row'>
                <label>Description</label>

                <ValidatedInput
                  name='description'
                  onChange={onChange}
                  value={environment.description}
                  placeholder={`The ${project.name} production environment`} />
              </div>

              <button onClick={this.onSubmit.bind(this)} disabled={loading || !nameValid} className='btn'>Create Environment</button>
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