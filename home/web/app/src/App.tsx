import React from 'react';
import Cookies from 'universal-cookie';
import DefaultIcon from  './assets/images/default-icon.png';
import Person from './assets/images/person.png';
import Call, { User, App } from './api';
import './App.scss';

interface Props {}

interface State {
  user?: User;
  apps: App[];
}

export default class AppComponent extends React.Component<Props, State> {
  readonly state: State = { apps: [] };

  componentDidMount() {
    Call('ReadUser').then((res: any) => this.setState({ user: new User(res.data.user) }));
    Call('ListApps').then((res: any) => this.setState({ apps: res.data.apps.map((a:any) => new App(a)) }));
  }

  render() {
    const now = new Date();
    
    const timeOpts =  { hour: "2-digit", minute: "2-digit", hour12: false }
    const time = now.toLocaleTimeString("en-uk", timeOpts);
    
    const dateOpts = { weekday: 'long', month: 'long', day: 'numeric' };
    const date = now.toLocaleDateString("en-uk", dateOpts);

    const { user, apps } = this.state;

    return (
      <div className="App">
        <div className='upper'>
          <div className='left'>
            <h1>{time}</h1>
            <p>{date}</p>
          </div>

          <div className={`right ${user ? '' : 'hidden'}`}>
            <p>Welcome back {user?.firstName}</p>

            <div className='dropdown'>
              <img src={Person} alt='My Account' />

              <div className="dropdown-content">
                <p onClick={() => window.location.href='/account?back=home'}>My Account</p>
                <p onClick={this.onLogoutPressed} className='logout'>Logout</p>
              </div>
            </div>
          </div>
        </div>

        <div className={`main ${apps.length > 0 ? '' : 'hidden'}`}>
          <div className='section'>
            <div className='section-upper'>
              <h3>Apps</h3>
              <p className='action'>Browse</p>
            </div>

            { apps.map(this.renderApp) }
          </div>
        </div>
      </div>
    );
  }

  renderApp(app: App): JSX.Element {
    return(
      <a key={app.id} className='AppCard' href={`/${app.id}`}>
        <img src={app.icon.length > 0 ? app.icon : DefaultIcon} alt='' />
        <p className='name'>{app.name}</p>
        <p className='category'>{app.category}</p>
      </a>
    )
  }

  onLogoutPressed() {
        // eslint-disable-next-line no-restricted-globals
    if(!confirm("Are you sure you want to logout?")) return;
    const cookies = new Cookies();
    cookies.remove('micro-token');
    window.location.reload();
  }
}