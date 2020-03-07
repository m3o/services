import React from 'react';
import { NavLink } from 'react-router-dom';
import Logo from '../../assets/images/logo.png';
import HomeActive from '../../assets/images/nav/home-active.png';
import HomeInactive from '../../assets/images/nav/home-inactive.png';
import NotesActive from '../../assets/images/nav/notes-active.png';
import NotesInactive from '../../assets/images/nav/notes-inactive.png';
import SprintsActive from '../../assets/images/nav/sprints-active.png';
import SprintsInactive from '../../assets/images/nav/sprints-inactive.png';
import './PageLayout.scss';

interface Props {
  className: string;
  match?: any;
}

export default class PageLayout extends React.Component<Props> {
  render(): JSX.Element {
    const { className, match } = this.props;
    const path = match.path
    
    return(
      <div className='PageLayout'>
        <div className='sidebar'>
          <div className='upper'>
            <NavLink exact to='/distributed'>
              <img src={Logo} alt='logo'/>
            </NavLink>
          </div>

          <nav>
            <NavLink exact to='/distributed'>
              <img src={ path === '/distributed/' ? HomeActive : HomeInactive } alt='Home' />
            </NavLink>

            <NavLink to='/distributed/notes'>
              <img src={ path.startsWith('/distributed/notes') ? NotesActive : NotesInactive } alt='Notes' />
            </NavLink>

            <NavLink to='/distributed/sprints'>
              <img src={ path.startsWith('/distributed/sprints') ? SprintsActive : SprintsInactive } alt='Sprints' />
            </NavLink>
          </nav>
        </div>

        <div className={`content ${className}`}>
          { this.props.children }
        </div>
      </div>
    );
  }
}