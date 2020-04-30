// Libraries
import React from 'react';
import { connect } from 'react-redux';
import { BrowserRouter, Route } from 'react-router-dom';

// Utils
import { State as GlobalState } from './store';
import { setUser } from './store/Account';
import * as API from './api';

// Scenes
import Notifications from './scenes/Notifications';
import Enviroment from './scenes/Enviroment';
import Project from './scenes/Project';

// Styling
import Logo from './components/PageLayout/assets/logo.png';
import './App.scss';

interface Props {
  user?: API.User;
  setUser: (user: API.User) => void;
}

class App extends React.Component<Props> {
  render(): JSX.Element {
    if(this.props.user) return this.renderLoggedIn();
    return this.renderLoading();
  }

  componentDidMount() {
    API.Call("AccountService/Read").then((res) => {
      this.props.setUser(res.data.user);
    });  
  }

  renderLoading(): JSX.Element {
    return <div className='loading'>
      <img src={Logo} alt='M3O' />
    </div>
  }

  renderLoggedIn(): JSX.Element {
    return (
      <BrowserRouter>
        <Route key='notificiations' exact path='/' component={Notifications} />
        <Route key='project' exact path='/projects/:project' component={Project} />
        <Route key='enviroment' exact path='/projects/:project/:enviroment' component={Enviroment} />
      </BrowserRouter>
    );  
  }
}

function mapStateToProps(state: GlobalState): any {
  return({
    user: state.account.user,
  });
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    setUser: (user: API.User) => dispatch(setUser(user)),
  });
}

export default connect(mapStateToProps, mapDispatchToProps)(App);
