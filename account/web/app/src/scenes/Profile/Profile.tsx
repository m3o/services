import React from 'react';
import { connect } from 'react-redux';
import Call, { User } from '../../api';
import PageLayout from '../../components/PageLayout';
import { setUser } from '../../store/User';
import './Profile.scss';

interface Props {
  user: User;
  updateUser: (user: User) => void;
}

interface State {
  saving: boolean;
  error: string;
}

class Profile extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { saving: false, error: '' };
  }
  
  onChange(e:any) {
    this.props.updateUser(new User({
      ...this.props.user,
      [e.target.name]: e.target.value,
    }));
  };

  async onSubmit(e:any) {
    e.preventDefault();
    this.setState({ saving: true });

    const { user } = this.props;

    Call("UpdateUser", { user })
      .then(() => this.setState({ error: '' }))
      .catch(err => this.setState({ error: err.message }))
      .finally(() => setTimeout(() => this.setState({ saving: false }), 500));
  }

  render(): JSX.Element {
    const { saving } = this.state;
    const { user } = this.props;

    return(
      <PageLayout className='Profile' {...this.props}>
        <form onSubmit={this.onSubmit.bind(this)}>
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
      </PageLayout>
    );
  }
}

function mapStateToProps(state: any): any {
  return({
    user: state.user.user,
  });
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    updateUser: (user: User) => dispatch(setUser(user)),
  });
}

export default connect(mapStateToProps, mapDispatchToProps)(Profile);