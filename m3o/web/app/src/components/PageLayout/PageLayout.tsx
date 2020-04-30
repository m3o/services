// Frameworks
import React from 'react';
import { NavLink } from 'react-router-dom';

// Styling
import Logo from './assets/logo.png';
import './PageLayout.scss';

interface Props {
  className?: string;
}

export default class PageLayout extends React.Component<Props> {
  render(): JSX.Element {
    return(
      <div className='PageLayout'>
        <div className='navbar'>
          <img src={Logo} alt='M3O Logo' className='logo' />

          <nav>
            <NavLink to='/'>
              <p>Dashboard</p>
            </NavLink>
            
            <NavLink exact to='/teams'>
              <p>Teams</p>
            </NavLink>

            <NavLink exact to='/billing'>
              <p>Billing</p>
            </NavLink>

            <NavLink exact to='/settings'>
              <p>Account</p>
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