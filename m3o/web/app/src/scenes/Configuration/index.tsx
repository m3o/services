import React from 'react';
import { connect } from 'react-redux';
import PageLayout from '../../components/PageLayout';
// import AddUser from './assets/add-user.png';
import * as API from '../../api';
import { State as GlobalState } from '../../store';
import { deleteEnvVar } from '../../store/Configuration';
import './style.scss';

interface Props {
  history: any;
  envVars: API.EnvVar[];
  deleteEnvVar: (user: API.EnvVar) => void;
}

interface State {
  revealed: API.EnvVar[];
}

class ConfigurationScene extends React.Component<Props> {
  readonly state: State = { revealed: [] };

  render(): JSX.Element {
    return(
      <PageLayout className='Configuration'>
        <header>
          <h1>Configuration</h1>
          
          <button className='btn' onClick={() => this.props.history.push('/configuration/add')}>
            {/* <img src={AddUser} alt='Add Configuraation' /> */}
            <p>Add Configuration</p>
          </button>
        </header>

        <table>
          <thead>
            <tr>
              <th>Service</th>
              <th>Key</th>
              <th>Value</th>
              <th>Actions</th>
            </tr>
          </thead>

          <tbody>
            { this.props.envVars.map(u => {
              const hidden = u.secret && !this.state.revealed.includes(u);

              let secretClassName = '';
              if(u.secret) secretClassName += 'secret';
              if(hidden) secretClassName += ' hidden';

              const toggleHidden = () => {
                debugger
                if(!u.secret) return;

                if(hidden) {
                  this.setState({ revealed: [...this.state.revealed, u] });
                } else {
                  this.setState({ revealed: this.state.revealed.filter(r => r !== u) });
                }
              }

              return (<tr key={u.service+u.key}>
                <td>{u.service}</td>
                <td>{u.key}</td>
                <td onClick={toggleHidden.bind(this)} className={secretClassName}>{hidden ? 'Hidden' : u.value}</td>
                <td>
                  <button className='warning' onClick={() => this.editEnvVar(u)}>Edit</button>
                  <button className='danger' onClick={() => this.deleteEnvVar(u)}>Delete</button>
                </td>
              </tr>);
            }) }
          </tbody>
        </table>
      </PageLayout>
    )
  }

  editEnvVar(envVar: API.EnvVar): void {
    this.props.history.push(`/configuration/${envVar.service}/${envVar.key}/edit`);
  }

  deleteEnvVar(envVar: API.EnvVar): void {
    // eslint-disable-next-line no-restricted-globals
    if (!confirm(`Are you sure you want to delete ${envVar.service} / ${envVar.key}?`)) return;
    this.props.deleteEnvVar(envVar);
  }
}

function mapStateToProps(state: GlobalState): any {
  return({
    envVars: state.configuration.envVars.sort(sortByServiceName),
  });
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    deleteEnvVar: (envVar: API.EnvVar) => dispatch(deleteEnvVar(envVar)),
  })
}

export default connect(mapStateToProps, mapDispatchToProps)(ConfigurationScene);

function sortByServiceName(a: API.EnvVar, b: API.EnvVar): number {
  const aName = (a.service + a.key).toUpperCase();
  const bName = (b.service + b.key).toUpperCase();
  if(aName > bName) return 1;
  if(aName < bName) return -1;
  return 0;
}