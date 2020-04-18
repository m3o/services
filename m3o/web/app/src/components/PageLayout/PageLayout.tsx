import React from 'react';
import { NavLink } from 'react-router-dom';
import Logo from './assets/logo.png';
import NavGettingStarted from './assets/nav-getting-started.png';
import NavTeam from './assets/nav-team.png';
import NavServices from './assets/nav-services.png';
import NavConfiguration from './assets/nav-configuration.png';
import NavBilling from './assets/nav-billing.png';
import NavSettings from './assets/nav-settings.png';
import './PageLayout.scss';

interface Props {
  className?: string;
}

export default class PageLayout extends React.Component<Props> {
  render(): JSX.Element {
    return(
      <div className='PageLayout'>
        <div className='sidebar'>
          <img src={Logo} alt='M3O Logo' className='logo' />

          <nav>
            <NavLink exact to='/'>
              <img src={NavGettingStarted} alt='Getting Started' />
              <p>Getting Started</p>
            </NavLink>

            <NavLink exact to='/team'>
              <img src={NavTeam} alt='Team' />
              <p>Team</p>
            </NavLink>

            <NavLink exact to='/services'>
              <img src={NavServices} alt='Services' />
              <p>Services</p>
            </NavLink>

            <NavLink exact to='/configuration'>
              <img src={NavConfiguration} alt='Configuration' />
              <p>Configuration</p>
            </NavLink>
            
            <NavLink exact to='/billing'>
              <img src={NavBilling} alt='Billing' />
              <p>Billing</p>
            </NavLink>

            <NavLink exact to='/settings'>
              <img src={NavSettings} alt='Settings' />
              <p>Settings</p>
            </NavLink>
          </nav>
        </div>

        <div className={`main ${this.props.className}`}>
          { this.props.children }
        </div>
      </div>
    );
  }
}