// Libraries
import React from 'react';
import { connect } from 'react-redux';
import { BrowserRouter, Route } from 'react-router-dom';

// Utils
import { State as GlobalState } from './store';
import { setUser } from './store/Account';
import * as API from './api';

// Scenes
import Dashboard from './scenes/Dashboard';

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
        <Route key='dashboard' exact path='/' component={Dashboard} />
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
