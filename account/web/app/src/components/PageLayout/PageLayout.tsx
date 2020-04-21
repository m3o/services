import React from 'react';
import { connect } from 'react-redux';
import { NavLink } from 'react-router-dom';
import BackArrow from '../../assets/images/back-arrow.png';
import ProfileActive from '../../assets/images/nav/profile-active.png';
import ProfileInctive from '../../assets/images/nav/profile-inactive.png';
import BillingActive from '../../assets/images/nav/billing-active.png';
import BillingInctive from '../../assets/images/nav/billing-inactive.png';
// import SubscriptionsActive from '../../assets/images/nav/subscriptions-active.png';
// import SubscriptionsInctive from '../../assets/images/nav/subscriptions-inactive.png';
import SettingsActive from '../../assets/images/nav/settings-active.png';
import SettingsInctive from '../../assets/images/nav/settings-inactive.png';
import './PageLayout.scss';

interface Props {
  className: string;
  match?: any;
  redirect?: string;
}

class PageLayout extends React.Component<Props> {
  render():JSX.Element {
    const { className, match, redirect } = this.props;
    const path = match.path

    let redirectUI: JSX.Element;
    if(redirect) { 
      redirectUI = (
        <a href={redirect} className='page-return-link'>
          <img src={BackArrow} alt='Return' />
          <p>Go back</p>
        </a>
      );
    } else {
      redirectUI = (
        <a href='https://web.micro.mu' className='page-return-link'>
          <img src={BackArrow} alt='Go Home' />
          <p>Home</p>
        </a>
      );
    }

    return(
      <div className='PageLayout'>
        <h1>Account Management</h1>
        { redirectUI }

        <div className='page-container'>
          <nav>
            <NavLink exact to=''>
              <img src={ path === '/' ? ProfileActive : ProfileInctive } alt='Profile' />
              <p>Profile</p>
            </NavLink>

            <NavLink exact to='/settings'>
              <img src={ path === '/settings' ? SettingsActive : SettingsInctive } alt='Settings' />
              <p>Settings</p>
            </NavLink>
          </nav>

          <div className={`page-content ${className}`}>
            { this.props.children }
          </div>
        </div>
      </div>
    )
  }
}

function mapStateToProps(state: any): any {
  return({
    redirect: state.redirect.path,
  });
}

export default connect(mapStateToProps)(PageLayout);