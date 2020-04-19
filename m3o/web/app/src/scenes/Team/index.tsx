import React from 'react';
import { connect } from 'react-redux';
import PageLayout from '../../components/PageLayout';
import AddUser from './assets/add-user.png';
import * as API from '../../api';
import { State as GlobalState } from '../../store';
import { deleteUser } from '../../store/Team';

interface Props {
  history: any;
  users: API.User[];
  deleteUser: (user: API.User) => void;
}

class TeamScene extends React.Component<Props> {
  render(): JSX.Element {
    return(
      <PageLayout className='Team'>
        <header>
          <h1>Team</h1>
          
          <button className='btn' onClick={() => this.props.history.push('/team/members/invite')}>
            <img src={AddUser} alt='Add User' />
            <p>Invite team members</p>
          </button>
        </header>

        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Roles</th>
              <th>Actions</th>
            </tr>
          </thead>

          <tbody>
            { this.props.users.map(u => <tr key={u.id}>
              <td>{u.firstName} {u.lastName}</td>
              <td>{u.email}</td>
              <td>{u.roles.join(', ')}</td>
              <td>
                <button className='warning' onClick={() => this.editUser(u)}>Edit</button>
                <button className='danger' onClick={() => this.deleteUser(u)}>Delete</button>
              </td>
            </tr>) }
          </tbody>
        </table>
      </PageLayout>
    )
  }

  editUser(user: API.User): void {
    this.props.history.push(`/team/members/${user.id}/edit`);
  }

  deleteUser(user: API.User): void {
    // eslint-disable-next-line no-restricted-globals
    if (!confirm(`Are you sure you want to delete ${user.firstName}?`)) return;
    this.props.deleteUser(user);
  }
}

function mapStateToProps(state: GlobalState): any {
  return({
    users: state.team.users.sort(sortByName),
  });
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    deleteUser: (user: API.User) => dispatch(deleteUser(user)),
  })
}

export default connect(mapStateToProps, mapDispatchToProps)(TeamScene);

function sortByName(a: API.User, b: API.User): number {
  const aName = (a.firstName + a.lastName).toUpperCase();
  const bName = (b.firstName + b.lastName).toUpperCase();
  if(aName > bName) return 1;
  if(aName < bName) return -1;
  return 0;
}