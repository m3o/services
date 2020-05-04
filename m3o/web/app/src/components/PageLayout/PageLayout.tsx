// Frameworks
import React from 'react';
import { connect } from 'react-redux';
import { NavLink } from 'react-router-dom';

// Utils
import * as API from '../../api';
import { State as GlobalState } from '../../store';

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
  projects: API.Project[];
}

class PageLayout extends React.Component<Props> {
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

            { this.props.projects.map(p => <section key={p.id}>
              <NavLink exact activeClassName='header active' className='header' to={`/projects/${p.name}`}>
                <p>{p.name}</p>
              </NavLink>

              <NavLink to={`/projects/${p.name}/production`.toLowerCase()}>
                <img src={ProjectIcon} alt={`${p.name}/production`} />
                <p>{p.name}/production</p>
              </NavLink>

              <NavLink to={`/projects/${p.name}/new`}>
                <img src={AddIcon} alt='New Enviroment' />
                <p>New Enviroment</p>
              </NavLink>
            </section>)}
          </div>

          <div className={`main ${this.props.className}`}>
            { this.props.children }
          </div>
        </div>
      </div>
    );
  }
}

function mapStateToProps(state: GlobalState): any {
  return({
    projects: state.project.projects,
  });
}

export default connect(mapStateToProps)(PageLayout);