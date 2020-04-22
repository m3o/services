import React from 'react';
import Cookies from 'universal-cookie';
import queryString from 'query-string';
import { connect } from 'react-redux';
import { withRouter } from "react-router-dom";
import Call, { User, Domain, Token } from '../../api';
import { setUser } from '../../store/User';
import GoogleLogo from '../../assets/images/google-logo.png';
import GitHubLogo from '../../assets/images/github-logo.png';
import './Login.scss';

interface Props {
  history: any;
  redirect?: string;
  setUser: (user: User) => void;
}

interface State {
  email: string;
  password: string;
  passwordConfirmation: string;
  loading: boolean;
  signup: boolean;
  error?: string;
  projectName?: string;
  inviteCode?: string;
}

interface Params {
  error?: string;
  inviteCode?: string;
  projectName?: string;
}

class Login extends React.Component<Props, State> {
  readonly state: State = { email: '', password: '', passwordConfirmation: '', loading: false, signup: false };

  componentDidMount() {
    const params: Params = queryString.parse(window.location.search);
    this.setState(params); // set inviteCode & projectName in the state
    }

  async onSubmit(event) {
    event.preventDefault();
    
    const { signup, email, password, passwordConfirmation, inviteCode, projectName } = this.state;
    if(signup && password !== passwordConfirmation) {
      this.setState({ error: 'Passwords must match' });
      return;
    } else if(password.length < 6) {
      this.setState({ error: 'Passwords must contain at least 6 characters' });
      return
    }
    
    this.setState({ loading: true, error: undefined });

    const params = signup ? { email, password, invite_code: inviteCode, project_invite: !!projectName } : { email, password };

    Call(signup ? 'Signup' : 'Login', params)
      .then((res) => {
        const user = new User(res.data.user);
        const token = new Token(res.data.token);
        this.props.setUser(user);
        
        const cookies = new Cookies();
        cookies.set('micro-token', token.access_token, { path: '/', domain: Domain, expires: token.expiry });                

        // check to see if the user needs onboarding
        if(user.requiresOnboarding()) {
          this.props.history.push('/signup');
          return
        }
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
      case 'passwordConfirmation':
        this.setState({ passwordConfirmation: e.target.value });
        return
      case 'inviteCode':
        this.setState({ inviteCode: e.target.value });
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

  redirectToOauth(name: string) {
    if(this.props.redirect) {
      const cookies = new Cookies();
      var expires = new Date();
      expires.setSeconds(expires.getSeconds() + 120);
      cookies.set('micro-account-redirect', this.props.redirect, { path: '/', domain: Domain, expires }); 
    }

    // pass the email to oauth if one was provided from the
    // email invite, this will auto-populate the oauth provider
    // and help to ensure the user signs up with the right email
    // address. The invite code will be stored in a cache so we
    // can retrieve it after the oauth flow is completed.
    if(!!this.state.projectName)  {
      const email = encodeURIComponent(this.state.email)
      const invite = encodeURIComponent(this.state.inviteCode)
      window.location.href = `/oauth/${name}/login?email=${email}&inviteCode=${invite}`;  
      return
    }

    window.location.href = `/oauth/${name}/login`;  
  }

  renderLogin(): JSX.Element {
    const { email, password, loading, error, projectName } = this.state;

    return(
      <div className='inner'>
        <h1>{ projectName? `Join the ${projectName} project` : 'Welcome back!'}</h1>
        <p className='subtitle'>To continue, log in with a Google or Micro account.</p>

        <div className='google oauth' onClick={() => this.redirectToOauth('google') }>
          <img src={GoogleLogo} alt='Sign in with Google' />
          <p>Sign in with Google</p>
        </div>

        <div className='github oauth' onClick={() => this.redirectToOauth('github') }>
          <img src={GitHubLogo} alt='Sign in with GitHub' />
          <p>Sign in with GitHub</p>
        </div>

        { error ? <p className='error'>Error: {error}</p> : null }

        <form onSubmit={this.onSubmit.bind(this)}>
          <label>Email *</label>
          <input
            type='email'
            name='email'
            value={email}
            autoFocus={!!projectName}
            disabled={loading || !!projectName}
            onChange={this.onChange.bind(this)} />

          <label>Password *</label>
          <input
            type='password'
            name='password'
            value={password}
            disabled={loading}
            autoFocus={!projectName}
            onChange={this.onChange.bind(this)} />
        
          <input
            type='submit'
            disabled={loading}
            value={loading ? 'Logging In' : 'Log in to your account'} />
        </form>

        <p className='signup'>Need an account? <span onClick={this.toggleSignup.bind(this)} className='link'>Create your Micro account.</span></p>
      </div>
    )
  }

  renderSignup(): JSX.Element {
    const { email, password, passwordConfirmation, loading, error, projectName, inviteCode } = this.state;

    return(
      <div className='inner'>
        <h1>{ projectName? `Join the ${projectName} project` : 'Signup'}</h1>
        <p className='subtitle'>Enter your email and password below to signup for a Micro account.</p>

        { error ? <p className='error'>Error: {error}</p> : null }

        <form onSubmit={this.onSubmit.bind(this)}>
          <label>Email *</label>
          <input
            type='email'
            name='email'
            value={email}
            disabled={loading || !!projectName}
            onChange={this.onChange.bind(this)} />

          <label>Password *</label>
          <input
            type='password'
            name='password'
            value={password} 
            disabled={loading}
            onChange={this.onChange.bind(this)} />

          <label>Password Confirmation *</label>
          <input
            type='password'
            disabled={loading}
            name='passwordConfirmation'
            value={passwordConfirmation}
            onChange={this.onChange.bind(this)} />

          { projectName ? null : <label>Invite Code *</label> }
          { projectName ? null : <input
                                required
                                type='text'
                                name='inviteCode'
                                disabled={loading}
                                value={inviteCode}
                                onChange={this.onChange.bind(this)} /> } 

          <input type='submit' value={loading ? 'Creating your account' : 'Create an account'} disabled={loading} />
        </form>

        <p className='signup'>Already have an account? <span onClick={this.toggleSignup.bind(this)} className='link'>Click here to login.</span></p>
      </div>
    )
  }
}

function mapStateToProps(state: any): any {
  return ({
    redirect: state.redirect.path,
  });
}

function mapDispatchToProps(dispatch: Function):any {
  return({
    setUser: (user: User) => dispatch(setUser(user)),
  });
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Login));
