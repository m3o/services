import React from 'react';
import { connect } from 'react-redux';
import PageLayout from '../../../../components/PageLayout';
import * as API from '../../../../api';
import { State as GlobalState } from '../../../../store';
import { updateUser } from '../../../../store/Team';

interface Props {
  user: API.User;
  updateUser: (user: API.User) => void;
  match: any;
  history: any;
}

interface State {
  user: API.User;
}

class EditTeamMemberScene extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { user: props.user };
  }

  render(): JSX.Element {
    const { user } = this.state;
    
    return(
      <PageLayout>
        <header>
          <h1>Edit {user.firstName} {user.lastName}</h1>

          <button className='btn danger' onClick={this.onCancel.bind(this)}>
            <p>Cancel</p>
          </button>

          <button className='btn' onClick={this.onSave.bind(this)}>
            <p>Save</p>
          </button>
        </header>

        <form onSubmit={(e: any) => {e.preventDefault(); this.onSave()}}>
          <label>First Name</label>
          <input
            required
            type='text' 
            name='firstName'
            value={user.firstName}
            onChange={this.onChange.bind(this)} />
            
          <label>Last Name</label>
          <input
            required
            type='text' 
            name='lastName'
            value={user.lastName}
            onChange={this.onChange.bind(this)} />
          
          <label>Email</label>
          <input
            required
            type='email'
            name='email' 
            value={user.email}
            onChange={this.onChange.bind(this)} />
        </form>
      </PageLayout>
    );
  }

  onChange(e: any): void {
    this.setState({
      user: {
        ...this.state.user,
        [e.target.name]: e.target.value,
      },
    });
  }

  onSave(): void {
    this.props.updateUser(this.state.user);
    this.props.history.push('/team');
  }

  onCancel(): void {
    // eslint-disable-next-line no-restricted-globals
    if (!confirm(`Are you sure you want to cancel? All your changes will be lost.`)) return;
    this.props.history.push('/team');
  }
}

function mapStateToProps(state: GlobalState, ownProps: Props): any {
  return({
    user: state.team.users.find(u => u.id === ownProps.match.params.id),
  });
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    updateUser: (user: API.User) => dispatch(updateUser(user)),
  })
}

export default connect(mapStateToProps, mapDispatchToProps)(EditTeamMemberScene);