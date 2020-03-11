import React from 'react';
import Call, { User } from './api';
import { connect } from 'react-redux';
import { BrowserRouter , Route } from 'react-router-dom';

// Scenes
import Profile from './scenes/Profile';
import Billing from './scenes/Billing';

// Assets
import Spinner from './assets/images/spinner.gif'; 
import './App.scss';
import { setUser } from './store/User';

interface Props {
  user?: User;
  setUser: (user: User) => void;
}

class App extends React.Component<Props> {
  componentDidMount() {
    Call("ReadUser")
      .then(res => this.props.setUser(res.data.user))
      .catch(console.warn);
  }

  render(): JSX.Element {
    const { user } = this.props;
    if(!user) return this.renderLoading();

    return (
      <BrowserRouter>
        <div className='App'>
          <Route exact path='/account/' component={Profile}/>
          <Route exact path='/account/billing' component={Billing}/>
        </div>
      </BrowserRouter>
    );
  }

  renderLoading(): JSX.Element {
    return(
      <div className="Loading">
        <img className='spinner' src={Spinner} alt='Loading' />
      </div>
    );
  }
}

function mapStateToProps(state: any): any {
  return({
    user: state.user.user,
  });
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    setUser: (user: User) => dispatch(setUser(user)),
  });
}

export default connect(mapStateToProps, mapDispatchToProps)(App);