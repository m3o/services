import React from 'react';
import { connect } from 'react-redux';
import PageLayout from '../../../../components/PageLayout';
import * as API from '../../../../api';
import { State as GlobalState } from '../../../../store';
import { addEnvVar } from '../../../../store/Configuration';
import Form from '../../components/Form';

interface Props {
  envVar: API.EnvVar;
  addEnvVar: (envVar: API.EnvVar) => void;
  match: any;
  history: any;
}

interface State {
  envVar: API.EnvVar;
}

class EditConfigurationService extends React.Component<Props, State> {
  readonly state: State = { envVar: {key: '', value: '', service: ''} };

  render(): JSX.Element {
    const { envVar } = this.state;
    
    return(
      <PageLayout>
        <header>
          <h1>Add Configuration</h1>

          <button className='btn danger' onClick={this.onCancel.bind(this)}>
            <p>Cancel</p>
          </button>

          <button className='btn' onClick={this.onSave.bind(this)}>
            <p>Save</p>
          </button>
        </header>

        <Form envVar={envVar} onSubmit={this.onSave.bind(this)} onUpdate={this.onUpdate.bind(this)} />
      </PageLayout>
    );
  }

  onUpdate(envVar: API.EnvVar): void {
    this.setState({ envVar });
  }

  onSave(): void {
    this.props.addEnvVar(this.state.envVar);
    this.props.history.push('/configuration');
  }

  onCancel(): void {
    // eslint-disable-next-line no-restricted-globals
    if (!confirm(`Are you sure you want to cancel? All your changes will be lost.`)) return;
    this.props.history.push('/configuration');
  }
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    addEnvVar: (envVar: API.EnvVar) => dispatch(addEnvVar(envVar)),
  })
}

export default connect(null, mapDispatchToProps)(EditConfigurationService);