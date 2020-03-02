import React from 'react';
import { NavLink } from 'react-router-dom';
import Logo from '../../assets/images/logo.png';
import './PageLayout.scss';

interface Props {
  className: string;
}

export default class PageLayout extends React.Component<Props> {
  render(): JSX.Element {
    const { className } = this.props;
    
    return(
      <div className='PageLayout'>
        <div className='sidebar'>
          <div className='upper'>
            <img src={Logo} alt='logo'/>
            {/* <h1>Distributed</h1> */}
          </div>

          <nav>
            <NavLink exact to='/distributed'>Home</NavLink>
            <NavLink to='/distributed/notes'>Notes</NavLink>
          </nav>
        </div>

        <div className={`content ${className}`}>
          { this.props.children }
        </div>
      </div>
    );
  }
}