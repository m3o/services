import React from 'react';
import Cookies from 'universal-cookie';
import Serverless from  './assets/images/serverless.png';
import Distributed from  './assets/images/distributed.png';
import Notes from  './assets/images/notes.png';
import Person from './assets/images/person.png';
import Call, { User } from './api';
import './App.scss';

export default class App extends React.Component {
  state: { user?: User } = {};

  componentDidMount() {
    Call('ReadUser').then((res: any) => this.setState({ user: res.data.user }));
  }

  render() {
    const now = new Date();
    
    const timeOpts =  { hour: "2-digit", minute: "2-digit", hour12: false }
    const time = now.toLocaleTimeString("en-uk", timeOpts);
    
    const dateOpts = { weekday: 'long', month: 'long', day: 'numeric' };
    const date = now.toLocaleDateString("en-uk", dateOpts);

    const { user } = this.state;

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
                <p onClick={() => window.location.href='/account'}>My Account</p>
                <p onClick={this.onLogoutPressed} className='logout'>Logout</p>
              </div>
            </div>
          </div>
        </div>

        <div className={`main ${user ? '' : 'hidden'}`}>
          <div className='section'>
            <div className='section-upper'>
              <h3>Apps</h3>
              <p className='action'>Browse</p>
            </div>

            <div className='AppCard'>
              <img src={Serverless} alt='Serverless' />
              <p className='name'>Serverless</p>
              <p className='category'>Development</p>
            </div>

            <div className='AppCard'>
              <img src={Distributed} alt='Distributed' />
              <p className='name'>Distributed</p>
              <p className='category'>Development</p>
            </div>
          
            <div className='AppCard'>
              <img src={Notes} alt='Notes' />
              <p className='name'>Notes</p>
              <p className='category'>Productivity</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  onLogoutPressed() {
        // eslint-disable-next-line no-restricted-globals
    if(!confirm("Are you sure you want to logout?")) return;
    const cookies = new Cookies();
    cookies.remove('micro-token');
    window.location.reload();
  }
}