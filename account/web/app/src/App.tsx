import React from 'react';
import Call, { User } from './api';
import { connect } from 'react-redux';
import queryString from 'query-string';
import { BrowserRouter , Route } from 'react-router-dom';

// Scenes
import Profile from './scenes/Profile';
import Billing from './scenes/Billing';
import Login from './scenes/Login';

// Assets
import Spinner from './assets/images/spinner.gif'; 
import './App.scss';
import { setUser } from './store/User';
import { setRedirect } from './store/Redirect';

interface Props {
  user?: User;
  setUser: (user: User) => void;
  setRedirect: (path: string) => void;
}

interface State {
  loading: boolean;
}

interface Params {
  redirect_to?: string;
}

const Routes = [
  <Route key='profile' exact path='/account/' component={Profile}/>,
  <Route key='billing' exact path='/account/billing' component={Billing}/>,
];

const UnauthenticatedRoutes = [
  <Route key='login' exact path='/account/' component={Login}/>,
]

class App extends React.Component<Props, State> {
  state = { loading: true };

  componentDidMount() {
    const params: Params = queryString.parse(window.location.search);
    if(params.redirect_to) this.props.setRedirect(params.redirect_to);
    
    Call("ReadUser")
      .then(res => this.props.setUser(res.data.user))
      .catch(console.warn)
      .finally(() => this.setState({ loading: false }));
  }

  render(): JSX.Element {
    if(this.state.loading) return this.renderLoading();

    return (
      <BrowserRouter>
        <div className='App'>
          { this.props.user ? Routes : UnauthenticatedRoutes }
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
    setRedirect: (path: string) => dispatch(setRedirect(path)),
  });
}

export default connect(mapStateToProps, mapDispatchToProps)(App);