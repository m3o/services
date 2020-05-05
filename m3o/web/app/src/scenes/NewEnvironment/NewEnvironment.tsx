// Frameworks
import React from 'react';
import { connect } from 'react-redux';

// Utils
import { State as GlobalState } from '../../store';
import * as API from '../../api';

// Components
import PageLayout from '../../components/PageLayout';

interface Props {
  match: any;
  project?: API.Project;
}

interface State {
  environment: API.Environment;
}

class NewEnvironment extends React.Component<Props, State> {
  readonly state: State = { environment: { name: '', description: '' } };

  onInputChange(e: any) {
    this.setState({
      environment: { 
        ...this.state.environment,
        [e.target.name]: e.target.value,
      },
    });
  }

  render(): JSX.Element {
    const { project } = this.props;
    if(!project) return null;

    const { environment } = this.state;

    return (
      <PageLayout className='NewEnvironment'>
        <div className='center'>
          <div className='header'>
            <h1>{project.name} / New Environment</h1>
          </div>

          <section>
            <h2>Environment Details</h2>
            <p>Set the name and description for your environment. You cannot change name once it is set.</p>

            <form>
              <div className='row'>
                <label>Name *</label>
                <input required type='text' value={environment.name} placeholder='production' name='name' onChange={this.onInputChange.bind(this)} />
              </div>
              
              <div className='row'>
                <label>Description</label>
                <input type='text' value={environment.description} placeholder={`The ${project.name} production environment`} name='description'  onChange={this.onInputChange.bind(this)} />
              </div>

              <button className='btn'>Create Environment</button>
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

export default connect(mapStateToProps)(NewEnvironment)