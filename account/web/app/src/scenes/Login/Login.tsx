import React from 'react';
import Cookies from 'universal-cookie';
import { connect } from 'react-redux';
import GoogleLogo from '../../assets/images/google-logo.png';
import './Login.scss';
import Call, { User, Domain } from '../../api';
import { setUser } from '../../store/User';

interface Props {
  setUser: (user: User) => void;
}

interface State {
  email: string;
  password: string;
  loading: boolean;
  signup: boolean;
  error?: string;
}

class Login extends React.Component<Props, State> {
  readonly state: State = { email: '', password: '', loading: false, signup: false };

  async onSubmit(event) {
    event.preventDefault();
    this.setState({ loading: true, error: undefined });

    const { email, password, signup } = this.state;
    const path = signup ? 'EmailSignup' : 'EmailLogin';

    Call(path, { email, password })
      .then((res) => {
        const cookies = new Cookies();
        cookies.set('micro-token', res.data.token, { path: '/', domain: Domain });        
        console.log(res.data.token);
        this.props.setUser(res.data.user);
      })
      .catch((err: any) => {
        const error = err.response ? err.response.data.detail : err.message;
        this.setState({ error, loading: false });
      });
  }

  onChange(e: any) {
    switch(e.target.name) {
      case 'email':
        this.setState({ email: e.target.value })
        return 
      case 'password':
        this.setState({ password: e.target.value })
        return
    }
  }

  toggleSignup() {
    this.setState({ signup: !this.state.signup });
  }

  render(): JSX.Element {
    const { signup } = this.state;

    return(
      <div className='Login'>
        { signup ? this.renderSignup() : this.renderLogin() }
      </div>
    )
  }

  renderLogin(): JSX.Element {
    const { email, password, loading, error } = this.state;

    return(
      <div className='inner'>
        <h1>Welcome back!</h1>
        <p className='subtitle'>To continue, log in with a Google or Micro account.</p>

        <div className='google-oauth' onClick={() => window.location.href = "/account/oauth/login"}>
          <img src={GoogleLogo} alt='Sign in with Google' />
          <p>Sign in with Google</p>
        </div>

        { error ? <p className='error'>Error: {error}</p> : null }

        <form onSubmit={this.onSubmit.bind(this)}>
          <label>Email *</label>
          <input type='email' name='email' value={email} disabled={loading} onChange={this.onChange.bind(this)} />

          <label>Password *</label>
          <input type='password' name='password' value={password} disabled={loading} onChange={this.onChange.bind(this)} />

          <input type='submit' value={loading ? 'Logging In' : 'Log in to your account'} disabled={loading} />
        </form>

        <p className='signup'>Need an account? <span onClick={this.toggleSignup.bind(this)} className='link'>Create your Micro account.</span></p>
      </div>
    )
  }

  renderSignup(): JSX.Element {
    const { email, password, loading, error } = this.state;

    return(
      <div className='inner'>
        <h1>Signup</h1>
        <p className='subtitle'>Enter your email and password below to signup for a Micro account.</p>

        { error ? <p className='error'>Error: {error}</p> : null }

        <form onSubmit={this.onSubmit.bind(this)}>
          <label>Email *</label>
          <input type='email' name='email' value={email} disabled={loading} onChange={this.onChange.bind(this)} />

          <label>Password *</label>
          <input type='password' name='password' value={password} disabled={loading} onChange={this.onChange.bind(this)} />

          <input type='submit' value={loading ? 'Logging In' : 'Log in to your account'} disabled={loading} />
        </form>

        <p className='signup'>Already have an account? <span onClick={this.toggleSignup.bind(this)} className='link'>Click here to login.</span></p>
      </div>
    )
  }
}

function mapDispatchToProps(dispatch: Function):any {
  return({
    setUser: (user: User) => dispatch(setUser(user)),
  });
}

export default connect(null, mapDispatchToProps)(Login);