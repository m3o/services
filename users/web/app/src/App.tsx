import React from 'react';
import Cookies from 'js-cookie';
import Call, { User } from './api';
import Spinner from './assets/images/spinner.gif'; 
import './App.scss';

interface Props {}

interface State {
  token?: string;
  error?: string;
  user?: User;
  saving: boolean;
}

export default class App extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);

    const token = Cookies.get('micro_token');
    this.state = { token, saving: false };
  }

  componentDidMount() {
    if(!this.state.token) return;

    Call("Read", this.state.token)
      .then(res => this.setState({ user: res.data.user }))
      .catch(err => this.setState({ error: err.message, token: undefined }))
  }

  onChange(e:any) {
    this.setState({ user: new User({
      ...this.state.user,
      [e.target.name]: e.target.value,
    })});
  };

  onSubmit(e:any) {
    e.preventDefault();
    this.setState({ saving: true });

    const { token, user } = this.state;
    Call("Update", token!, { user })
      .then(() => this.setState({ error: '' }))
      .catch(err => this.setState({ error: err.message }))
      .finally(() => this.setState({ saving: false }))
  }

  render(): JSX.Element {
    const { error, token, user, saving } = this.state;
    if(!token) return this.renderNoToken();
    if(!user) return this.renderLoading();

    return (
      <div className="App">
        <h1>Your Profile</h1>
        <p className='error'>{error}</p>

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
          
          <input disabled={this.state.saving} type='submit' value={ saving ? 'Saving' : 'Save Changes' } />
        </form>
      </div>
    );
  }

  renderNoToken(): JSX.Element {
    return(
      <div className="App">
        <h1>Not Logged In</h1>
        <p>You cannot edit your profile as you're not logged in.</p>
      </div>
    );
  }

  renderLoading(): JSX.Element {
    return(
      <div className="App">
        <img className='spinner' src={Spinner} alt='Loading' />
      </div>
    );
  }
}