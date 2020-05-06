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
  childRef?: React.RefObject<HTMLDivElement>;
  className?: string;
  projects: API.Project[];
  hideSidebar?: boolean;
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
            
            <NavLink exact to='/billing'>
              <p>Billing</p>
            </NavLink>

            <a href='https://account.micro.mu/' target='blank'>
              <p>Account</p>
            </a>
          </nav>
        </div>

        <div className='wrapper'>
          { this.props.hideSidebar ? null :this.renderSidebar() }
          
          <div className={`main ${this.props.className}`} ref={this.props.childRef}>
            { this.props.children }
          </div>
        </div>
      </div>
    );
  }

  renderSidebar(): JSX.Element {
    return(
      <div className='sidebar'>
        <section>
          <NavLink exact to='/'>
            <img src={NotificationsIcon} alt='Notifications' />
            <p>{this.props.projects.length === 0 ? 'Getting Started' : 'Notifications'}</p>
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

          { p.environments?.map(e => <NavLink key={e.id} to={`/projects/${p.name}/${e.name}`}>
            <img src={ProjectIcon} alt={`${p.name}/${e.name}`} />
            <p>{p.name}/{e.name}</p>
          </NavLink> ) }

          <NavLink to={`/new/environment/${p.name}`}>
            <img src={AddIcon} alt='New Enviroment' />
            <p>New Enviroment</p>
          </NavLink>
        </section>)}

        <section>
          <NavLink exact activeClassName='header active' className='header' to={`/new/project`}>
            <p>New Project</p>
          </NavLink>
        </section>
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