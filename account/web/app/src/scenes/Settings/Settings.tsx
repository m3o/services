import React from 'react';
import Cookies from 'universal-cookie';
import PageLayout from '../../components/PageLayout';
import './Settings.scss';
import Call from '../../api';

export default class Settings extends React.Component {
  render(): JSX.Element {
    return(
      <PageLayout className='Settings' {...this.props}>
        <div className='section'>
          <h3>Logout</h3>
          <p>Press the button below to logout. You will need to log back in to access your account.</p>
          <button onClick={this.logout}>Logout</button>
        </div>

        <div className='section'>
          <h3>Delete Account</h3>
          <p>Press the button below to permentantly delete your account. This action is permenant.</p>
          <button className='danger' onClick={this.deleteAccount}>Delete account</button>
        </div>        
      </PageLayout>
    )
  }

  logout() {
    // eslint-disable-next-line no-restricted-globals
    if(!confirm("Are you sure you want to logout?")) return;

    // remove cookies
    const cookies = new Cookies();
    cookies.remove("micro-token", {path: "/", domain: "micro.mu"});

    // reload so micro web will redirect to login
    window.location.href = '/account';
  }

  deleteAccount() {
    // eslint-disable-next-line no-restricted-globals
    if(!confirm("Are you sure you want to delete your account?")) return;

    Call('DeleteUser').then(() => {
      // remove cookies
      const cookies = new Cookies();
      cookies.remove("micro-token", {path: "/", domain: "micro.mu"});

      // reload so micro web will redirect to login
      window.location.href = '/account';
    }).catch(err => {
      // eslint-disable-next-line no-restricted-globals
      alert("There was a problem ddeleting your account: " + err)
    })

  }
}