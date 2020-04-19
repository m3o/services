import React from 'react';
import { connect } from 'react-redux';
import PageLayout from '../../../../components/PageLayout';
import * as API from '../../../../api';
import { State as GlobalState } from '../../../../store';
import { updateEnvVar } from '../../../../store/Configuration';

interface Props {
  envVar: API.EnvVar;
  updateEnvVar: (envVar: API.EnvVar) => void;
  match: any;
  history: any;
}

interface State {
  envVar: API.EnvVar;
}

class EditConfigurationService extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { envVar: props.envVar };
  }

  render(): JSX.Element {
    const { envVar } = this.state;
    
    return(
      <PageLayout>
        <header>
          <h1>Edit {envVar.key}</h1>

          <button className='btn danger' onClick={this.onCancel.bind(this)}>
            <p>Cancel</p>
          </button>

          <button className='btn' onClick={this.onSave.bind(this)}>
            <p>Save</p>
          </button>
        </header>

        <form onSubmit={(e: any) => {e.preventDefault(); this.onSave()}}>
          <label>Service</label>
          <select value={envVar.service} onChange={this.onServiceChange.bind(this)}>
            <option value='*'>All (*)</option>
            <option value='go.micro.service.payments'>go.micro.service.payments</option>
            <option value='go.micro.service.users'>go.micro.service.users</option>
            <option value='go.micro.service.foo'>go.micro.service.foo</option>
            <option value='go.micro.service.bar'>go.micro.service.bar</option>
          </select>

          <label>Key</label>
          <input
            required
            type='text' 
            name='key'
            value={envVar.key}
            onChange={this.onChange.bind(this)} />
            
          <label>Value</label>
          <input
            required
            type='string' 
            name='value'
            value={envVar.value}
            onChange={this.onChange.bind(this)} />
          
          <label>Secret</label>
          <select value={envVar.secret ? 'yes' : 'no'} onChange={this.onSecretChange.bind(this)}>
            <option value='yes'>Yes</option>
            <option value='no'>No</option>
          </select>
        </form>
      </PageLayout>
    );
  }

  onChange(e: any): void {
    this.setState({
      envVar: {
        ...this.state.envVar,
        [e.target.name]: e.target.value,
      },
    });
  }


  onSecretChange(e): void {
    this.setState({
      envVar: {
        ...this.state.envVar,
        secret: e.target.value === 'yes',
      },
    });
  }

  onServiceChange(e): void {
    this.setState({
      envVar: {
        ...this.state.envVar,
        service: e.target.value,
      },
    });
  }

  onSave(): void {
    this.props.updateEnvVar(this.state.envVar);
    this.props.history.push('/configuration');
  }

  onCancel(): void {
    // eslint-disable-next-line no-restricted-globals
    if (!confirm(`Are you sure you want to cancel? All your changes will be lost.`)) return;
    this.props.history.push('/configuration');
  }
}

function mapStateToProps(state: GlobalState, ownProps: Props): any {
  const { params } = ownProps.match;
  return({
    envVar: state.configuration.envVars.find(e => e.service === params.service && e.key === params.key),
  });
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    updateEnvVar: (envVar: API.EnvVar) => dispatch(updateEnvVar(envVar)),
  })
}

export default connect(mapStateToProps, mapDispatchToProps)(EditConfigurationService);