import React from 'react';
import GoogleLogo from '../../assets/images/google-logo.png';
import './Login.scss';

export default class Login extends React.Component {
  render(): JSX.Element {
    return(
      <div className='Login'>
        <div className='inner'>
          <h1>Welcome back!</h1>
          <p>To continue, log in with a Google or Micro account.</p>

          <div className='google-oauth' onClick={() => window.location.href = "/account/oauth/login"}>
            <img src={GoogleLogo} alt='Sign in with Google' />
            <p>Sign in with Google</p>
          </div>

          <form>
            <label>Email *</label>
            <input type='email' />

            <label>Password *</label>
            <input type='password' />

            <input type='submit' value='Log in to your account' />
          </form>

          <p className='signup'>Need an account? <span className='link'>Create your Micro account.</span></p>
        </div>
      </div>
    )
  }
}