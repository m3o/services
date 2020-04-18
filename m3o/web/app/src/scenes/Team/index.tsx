import React from 'react';
import PageLayout from '../../components/PageLayout';
import AddUser from './assets/add-user.png';

interface User {
  firstName: string;
  lastName: string;
  email: string;
  roles: string[];
  me?: boolean;
}

const users: User[] = [
  {firstName: "Asim", lastName: "Aslam", email: "asim@micro.mu", roles: ["Admin", "Developer"], me: true},
  {firstName: "Jake", lastName: "Sanders", email: "jake@micro.mu", roles: ["Developer"]},
  {firstName: "Ben", lastName: "Toogood", email: "ben@micro.mu", roles: ["Developer"]},
  {firstName: "Janos", lastName: "Dobronszki", email: "janos@micro.mu", roles: ["Developer"]},
  {firstName: "Vasiliy", lastName: "Tolstov", email: "vasiliy@micro.mu", roles: ["Developer"]},
];

export default class TeamScene extends React.Component {
  render(): JSX.Element {
    return(
      <PageLayout className='Team'>
        <header>
          <h1>Team</h1>
          
          <button className='btn'>
            <img src={AddUser} alt='Add User' />
            <p>Invite a team member</p>
          </button>
        </header>

        <table>
          <thead>
            <tr>
              <th>First Name</th>
              <th>Last Name</th>
              <th>Email</th>
              <th>Roles</th>
              <th>Actions</th>
            </tr>
          </thead>

          <tbody>
            { users.map(u => <tr key={u.email}>
              <td>{u.firstName}</td>
              <td>{u.lastName}</td>
              <td>{u.email}</td>
              <td>{u.roles.join(', ')}</td>
              <td>
                <button className='warning'>Edit</button>
                <button className='danger'>Delete</button>
              </td>
            </tr>) }
          </tbody>
        </table>
      </PageLayout>
    )
  }
}