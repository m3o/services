import React from 'react';
import Cookies from 'universal-cookie';
import Call, { User, Domain } from './api';
import { connect } from 'react-redux';
import { withRouter } from "react-router-dom";
import queryString from 'query-string';
import { Route } from 'react-router-dom';

// Scenes
import Profile from './scenes/Profile';
import Billing from './scenes/Billing';
import Settings from './scenes/Settings';
import Login from './scenes/Login';
import Onboarding from './scenes/Onboarding';

// Assets
import Spinner from './assets/images/spinner.gif'; 
import './App.scss';
import { setUser } from './store/User';
import { setRedirect } from './store/Redirect';

interface Props {
  user?: User;
  history: any;
  redirect?: string;
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
  <Route key='settings' exact path='/account/settings' component={Settings}/>,
  <Route key='onboarding' exact path='/account/onboarding' component={Onboarding}/>,
];

const UnauthenticatedRoutes = [
  <Route key='login' exact path='/account/' component={Login}/>,
  <Route key='onboarding' exact path='/account/onboarding' component={Onboarding}/>,

]

class App extends React.Component<Props, State> {
  state = { loading: true };

  componentDidMount() {
    const params: Params = queryString.parse(window.location.search);
    
    if(params.redirect_to) {
      this.props.setRedirect(params.redirect_to);
    } else {
      const cookies = new Cookies();
      this.props.setRedirect(cookies.get('micro-account-redirect'));
      cookies.remove('micro-account-redirect', { path: '/', domain: Domain });
    }
    
    Call("ReadUser")
      .then(this.setUser.bind(this))
      .catch(console.warn)
      .finally(() => this.setState({ loading: false }));
  }

  setUser(res: any) {
    // construct the user
    const user = new User(res.data.user);

    // redirect the user upon login
    if(this.props.redirect && !user.requiresOnboarding()) {
      window.location.href = this.props.redirect;
    }
    
    // set the user in the redux store
    this.props.setUser(user);

    // check to see if the user requires onboarding
    if(user.requiresOnboarding()) this.props.history.push('/account/onboarding');
  }

  render(): JSX.Element {
    if(this.state.loading) return this.renderLoading();

    return (
      <div className='App'>
        { this.props.user ? Routes : UnauthenticatedRoutes }
      </div>
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
    redirect: state.redirect.path,
  });
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    setUser: (user: User) => dispatch(setUser(user)),
    setRedirect: (path: string) => dispatch(setRedirect(path)),
  });
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(App));