// Frameworks
import React from 'react';
import { NavLink } from 'react-router-dom';

// Styling
import Logo from './assets/logo.png';
import ProjectIcon from './assets/project.png';
import AddIcon from './assets/add.png';
import NotificationsIcon from './assets/notifications.png';
import FeedbackIcon from './assets/feedback.png';
import DocsIcon from './assets/docs.png';
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

        <div className='wrapper'>
          <div className='sidebar'>
            <section>
              <NavLink exact to='/'>
                <img src={NotificationsIcon} alt='Notifications' />
                <p>Notifications</p>
              </NavLink>

              <NavLink exact to='/feedback'>
                <img src={FeedbackIcon} alt='Feedback' />
                <p>Feedback</p>
              </NavLink>

              <NavLink exact to='/docs'>
                <img src={DocsIcon} alt='Docs' />
                <p>Docs</p>
              </NavLink>
            </section>

            <section>
              <NavLink exact activeClassName='header active' className='header' to='/projects/ben-toogood'>
                <p>ben-toogood</p>
              </NavLink>
              
              <NavLink to='/projects/ben-toogood/hello-world'>
                <img src={ProjectIcon} alt='ben-toogood/hello-world' />
                <p>ben-toogood/hello-world</p>
              </NavLink>

              <NavLink to='/projects/ben-toogood/new'>
                <img src={AddIcon} alt='New Project' />
                <p>New Enviroment</p>
              </NavLink>
            </section>

            <section>
              <NavLink exact activeClassName='header active' className='header' to='/projects/kytra'>
                <p>Kytra</p>
              </NavLink>
              
              <NavLink to='/projects/kytra/production'>
                <img src={ProjectIcon} alt='kytra/production' />
                <p>kytra/production</p>
              </NavLink>
              
              <NavLink to='/projects/kytra/staging'>
                <img src={ProjectIcon} alt='kytra/staging' />
                <p>kytra/staging</p>
              </NavLink>
              
              <NavLink to='/projects/kytra/develpment'>
                <img src={ProjectIcon} alt='kytra/develpment' />
                <p>kytra/develpment</p>
              </NavLink>
              

              <NavLink to='/projects/kytra/new'>
                <img src={AddIcon} alt='New Project' />
                <p>New Enviroment</p>
              </NavLink>
            </section>

            <section>
              <NavLink exact activeClassName='header active' className='header' to='/project/Micro'>
                <p>Micro</p>
              </NavLink>
              
              <NavLink to='/projects/micro/services'>
                <img src={ProjectIcon} alt='micro/services' />
                <p>micro/services</p>
              </NavLink>
              
              <NavLink to='/projects/micro/m3o'>
                <img src={ProjectIcon} alt='micro/m3o' />
                <p>micro/m3o</p>
              </NavLink>
              
              <NavLink to='/projects/micro/new'>
                <img src={AddIcon} alt='New Project' />
                <p>New Enviroment</p>
              </NavLink>
            </section>
          </div>

          <div className={`main ${this.props.className}`}>
            { this.props.children }
          </div>
        </div>
      </div>
    );
  }
}