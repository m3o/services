// Libraries
import React from 'react';
import { BrowserRouter, Route } from 'react-router-dom';
import { Provider } from 'react-redux';

// Scenes
import GettingStartedScene from './scenes/GettingStarted';
import TeamScene from './scenes/Team';
import EditTeamMemberScene from './scenes/Team/scenes/EditTeamMember';
import InviteTeamMembersScene from './scenes/Team/scenes/InviteTeamMembers';
import ConfigurationScene from './scenes/Configuration';
import EditConfigurationScene from './scenes/Configuration/scenes/EditConfiguration';

// Reducer
import store from './store';

// Styling
import './App.scss';

// Redux Setup
window.store = store; 

// Declare global window interface so we can mount redux
declare global {
  interface Window {
    __REDUX_DEVTOOLS_EXTENSION__: any;
    store: any;
  }
}

function App() {
  return (
    <Provider store={window.store} >
      <BrowserRouter>
        <Route key='getting-started' exact path='/' component={GettingStartedScene} />
        <Route key='team' exact path='/team' component={TeamScene} />
        <Route key='edit-team-member' path='/team/members/:id/edit' component={EditTeamMemberScene} />
        <Route key='invite-team-members' path='/team/members/invite' component={InviteTeamMembersScene} />
        <Route key='configuration' exact path='/configuration' component={ConfigurationScene} />
        <Route key='edit-configuration' path='/configuration/:service/:key/edit' component={EditConfigurationScene} />
      </BrowserRouter>
    </Provider>
  );
}

export default App;
