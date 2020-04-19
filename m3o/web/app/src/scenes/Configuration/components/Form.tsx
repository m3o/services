import React from 'react';
import * as API from '../../../api';

interface Props {
  envVar: API.EnvVar;
  onUpdate: (envVar: API.EnvVar) => void;
  onSubmit: () => void;
}

export default class Form extends React.Component<Props> {
  constructor(props: Props) {
    super(props);
    this.state = { envVar: props.envVar };
  }

  render(): JSX.Element {    
    const { envVar } = this.props;

    return(
      <form onSubmit={(e: any) => {e.preventDefault(); this.props.onSubmit()}}>
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
    );
  }

  onChange(e: any): void {
    this.props.onUpdate({...this.props.envVar, [e.target.name]: e.target.value});
  }


  onSecretChange(e): void {
    this.props.onUpdate({...this.props.envVar, secret: e.target.value === 'yes'});
  }

  onServiceChange(e): void {
    this.props.onUpdate({...this.props.envVar, service: e.target.value});
  }
}