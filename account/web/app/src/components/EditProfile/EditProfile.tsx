import React from 'react';
import { connect } from 'react-redux';
import Call, { User } from '../../api';
import { setUser } from '../../store/User';
import './EditProfile.scss';

interface Props {
  user: User;
  updateUser: (user: User) => void;
  onSave?: () => void;
  buttonText?: string;
}

interface State {
  saving: boolean;
  error: string;
}

class EditProfile extends React.Component<Props, State> {
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
      .then((res) => {
        this.setState({ error: '' });
        this.props.updateUser(new User(res.data.user));
        if(this.props.onSave) this.props.onSave();
      })
      .catch(err => this.setState({ error: err.response.data.detail }))
      .finally(() => setTimeout(() => {
        this.setState({ saving: false });
      }, 500));
  }

  render(): JSX.Element {
    const { saving } = this.state;
    const { user } = this.props;

    return(
      <form className='EditProfile' onSubmit={this.onSubmit.bind(this)}>
        { this.state.error.length > 0 ? <p className='error'>{this.state.error}</p> : null }

        { user.invite_verified ? null : <label>Invite Code *</label> }
        { user.invite_verified ? null : <input
          required
          type='text'
          name='invite_code'
          value={user!.invite_code} 
          disabled={this.state.saving}
          onChange={this.onChange.bind(this)} /> }

        <label>First Name *</label>
        <input
          required
          type='text'
          name='first_name'
          value={user!.first_name} 
          disabled={this.state.saving}
          onChange={this.onChange.bind(this)} />
        
        <label>Last Name *</label>
        <input
          required
          type='text'
          name='last_name'
          value={user!.last_name} 
          disabled={this.state.saving}
          onChange={this.onChange.bind(this)} />
        
        <label>Email *</label>
        <input
          required
          name='email'
          type='email'
          value={user!.email}
          disabled={true} />
        
        <input disabled={this.state.saving} type='submit' value={ saving ? 'Saving' : (this.props.buttonText || 'Save Changes') } />
      </form>
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

export default connect(mapStateToProps, mapDispatchToProps)(EditProfile);