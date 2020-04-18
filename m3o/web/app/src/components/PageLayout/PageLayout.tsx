import React from 'react';
import { NavLink } from 'react-router-dom';
import Logo from './assets/logo.png';
import './PageLayout.scss';

interface Props {
  className?: string;
}

export default class PageLayout extends React.Component<Props> {
  render(): JSX.Element {
    return(
      <div className='PageLayout'>
        <div className='sidebar'>
          <img src={Logo} alt='M3O Logo' />

          <nav>
            <NavLink exact to='/'>
              <p>Getting Started</p>
            </NavLink>

            <NavLink exact to='/team'>
              <p>Team</p>
            </NavLink>

            <NavLink exact to='/services'>
              <p>Services</p>
            </NavLink>

            <NavLink exact to='/configuration'>
              <p>Configuration</p>
            </NavLink>
            
            <NavLink exact to='/billing'>
              <p>Billing</p>
            </NavLink>

            <NavLink exact to='/settings'>
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