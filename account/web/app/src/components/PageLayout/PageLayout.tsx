import React from 'react';
import { NavLink } from 'react-router-dom';
import BackArrow from '../../assets/images/back-arrow.png';
import ProfileActive from '../../assets/images/nav/profile-active.png';
import ProfileInctive from '../../assets/images/nav/profile-inactive.png';
import BillingActive from '../../assets/images/nav/billing-active.png';
import BillingInctive from '../../assets/images/nav/billing-inactive.png';
import SubscriptionsActive from '../../assets/images/nav/subscriptions-active.png';
import SubscriptionsInctive from '../../assets/images/nav/subscriptions-inactive.png';
import SettingsActive from '../../assets/images/nav/settings-active.png';
import SettingsInctive from '../../assets/images/nav/settings-inactive.png';
import './PageLayout.scss';

interface Props {
  className: string;
  match?: any;
}

export default class PageLayout extends React.Component<Props> {
  render():JSX.Element {
    const { className, match } = this.props;
    const path = match.path

    return(
      <div className='PageLayout'>
        <h1>Account Management</h1>

        <div className='page-return-link'>
          <img src={BackArrow} alt='Go back' />
          <p>Back to FooBar</p>
        </div>

        <div className='page-container'>
          <nav>
            <NavLink exact to='/account'>
              <img src={ path === '/account/' ? ProfileActive : ProfileInctive } alt='Profile' />
              <p>Profile</p>
            </NavLink>

            <NavLink exact to='/account/billing'>
              <img src={ path === '/account/billing/' ? BillingActive : BillingInctive } alt='Billing' />
              <p>Billing</p>
            </NavLink>

            <NavLink exact to='/account/subscriptions'>
              <img src={ path === '/account/subscriptions/' ? SubscriptionsActive : SubscriptionsInctive } alt='Subscriptions' />
              <p>Subscriptions</p>
            </NavLink>

            <NavLink exact to='/account/settings'>
              <img src={ path === '/account/settings/' ? SettingsActive : SettingsInctive } alt='Settings' />
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