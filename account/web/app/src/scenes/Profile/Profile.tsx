import React from 'react';
import Call, { User } from '../../api';
import './Profile.scss';

interface Props {
  user: User;
}

interface State {
  saving: boolean;
  user?: User;
  error: string;
}

export default class Profile extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { user: props.user, saving: false, error: '' };
  }
  
  componentDidUpdate(prevProps: Props)  {
    if(prevProps.user !== this.props.user) {
      this.setState({ user: this.props.user });
    }
  }

  onChange(e:any) {
    this.setState({ user: new User({
      ...this.state.user,
      [e.target.name]: e.target.value,
    })});
  };

  async onSubmit(e:any) {
    e.preventDefault();
    this.setState({ saving: true });

    const { user } = this.state;

    Call("UpdateUser", { user })
      .then(() => this.setState({ error: '' }))
      .catch(err => this.setState({ error: err.message }))
      .finally(() => this.setState({ saving: false }));
  }

  render(): JSX.Element {
    const { saving, user } = this.state;
    debugger

    return(
      <div className='Profile'>
        <form onSubmit={this.onSubmit.bind(this)}>
          <h3>Profile</h3>

          <label>First Name</label>
          <input
            type='text'
            name='firstName'
            value={user!.firstName} 
            disabled={this.state.saving}
            onChange={this.onChange.bind(this)} />
          
          <label>Last Name</label>
          <input
            type='text'
            name='lastName'
            value={user!.lastName} 
            disabled={this.state.saving}
            onChange={this.onChange.bind(this)} />
          
          <label>Email</label>
          <input
            name='email'
            type='email'
            value={user!.email}
            disabled={this.state.saving}
            onChange={this.onChange.bind(this)} />
          
          <label>Username</label>
          <input
            name='username'
            type='text'
            value={user!.username}
            disabled={this.state.saving}
            onChange={this.onChange.bind(this)} />

          <input disabled={this.state.saving} type='submit' value={ saving ? 'Saving' : 'Save Changes' } />
        </form>
      </div>
    );
  }
}