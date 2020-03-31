import React from 'react';
import Cookies from 'universal-cookie';
import PageLayout from '../../components/PageLayout';
import './Settings.scss';
import Call from '../../api';

interface State {
  copied: boolean;
}

export default class Settings extends React.Component {
  state: State = { copied: false };

  render(): JSX.Element {
    const cookies = new Cookies();
    const command = `micro login --platform --token=${cookies.get("micro-token")}`;

    return(
      <PageLayout className='Settings' {...this.props}>
        <div className='section'>
          <h3>Login to CLI</h3>
          <p>Copy the command below to login with the Micro CLI</p>
          <input id='command' value={command} className='code'/>

          <button className='fixed-width' onClick={this.copyLogin.bind(this)}>
            { this.state.copied ? 'Done ✅' : 'Copy Login Command' }
          </button>
        </div>

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

  copyLogin() {
    /* Get the text field */
    var copyText = document.getElementById("command") as HTMLInputElement;

    /* Select the text field */
    copyText.select();
    copyText.setSelectionRange(0, 99999); /*For mobile devices*/

    /* Copy the text inside the text field */
    document.execCommand("copy");

    /* Update the UI */
    this.setState({ copied: true });
  }

  logout() {
    // eslint-disable-next-line no-restricted-globals
    if(!confirm("Are you sure you want to logout?")) return;

    // remove cookies
    const cookies = new Cookies();
    cookies.remove("micro-token", {path: "/", domain: "micro.mu"});

    // reload so micro web will redirect to login
    window.location.href = '/';
  }

  deleteAccount() {
    // eslint-disable-next-line no-restricted-globals
    if(!confirm("Are you sure you want to delete your account?")) return;

    Call('DeleteUser').then(() => {
      // remove cookies
      const cookies = new Cookies();
      cookies.remove("micro-token", {path: "/", domain: "micro.mu"});

      // reload so micro web will redirect to login
      window.location.href = '/';
    }).catch(err => {
      // eslint-disable-next-line no-restricted-globals
      alert("There was a problem ddeleting your account: " + err)
    })

  }
}