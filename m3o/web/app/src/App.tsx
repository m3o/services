// Libraries
import React from 'react';
import { connect } from 'react-redux';
import { BrowserRouter, Route } from 'react-router-dom';

// Utils
import { State as GlobalState } from './store';
import { setUser } from './store/Account';
import * as API from './api';

// Scenes
import Billing from './scenes/Billing';
import Notifications from './scenes/Notifications';
import Enviroment from './scenes/Enviroment';
import Project from './scenes/Project';
import NewProject from './scenes/NewProject';
import NewEnvironment from './scenes/NewEnvironment';
import GettingStarted from './scenes/GettingStarted';
import NotFound from './scenes/404';

// Styling
import Logo from './components/PageLayout/assets/logo.png';
import './App.scss';
import { setProjects } from './store/Project';

interface Props {
  user?: API.User;
  projects: API.Project[];
  setUser: (user: API.User) => void;
  setProjects: (projects: API.Project[]) => void;
}

interface State {
  loaded: boolean;
}

class App extends React.Component<Props, State> {
  readonly state: State = { loaded: false };

  render(): JSX.Element {
    if(this.state.loaded) return this.renderLoggedIn();
    return this.renderLoading();
  }

  componentDidMount() {
    API.Call("Accounts/Read").then((res) => {
      this.props.setUser(res.data.user);

      API.Call("Projects/ListProjects").then((res) => {
        this.props.setProjects(res.data.projects || []);
        this.setState({ loaded: true });
      });  
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
        { this.props.projects.length > 0 ? <Route key='notificiations' exact path='/' component={Notifications} />  : null }
        { this.props.projects.length === 0 ? <Route key='getting-started' exact path='/' component={GettingStarted} /> : null }

        <Route key='billing' exact path='/billing' component={Billing} />
        <Route key='new-project' exact path='/new/project' component={NewProject} />
        <Route key='new-environmnt' exact path='/new/environment/:project' component={NewEnvironment} />
        <Route key='project' exact path='/projects/:project' component={Project} />
        <Route key='environment' exact path='/projects/:project/:environment' component={Enviroment} />
        <Route key='not-found' exact path='/not-found' component={NotFound} />
      </BrowserRouter>
    );  
  }
}

function mapStateToProps(state: GlobalState): any {
  return({
    user: state.account.user,
    projects: state.project.projects,
  });
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    setUser: (user: API.User) => dispatch(setUser(user)),
    setProjects: (projects: API.Project[]) => dispatch(setProjects(projects)),
  });
}

export default connect(mapStateToProps, mapDispatchToProps)(App);
