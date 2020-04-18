import React from 'react';
import { BrowserRouter, Route } from 'react-router-dom';
import './App.scss';
import GettingStartedScene from './scenes/GettingStarted';
import TeamScene from './scenes/Team';

function App() {
  return (
    <BrowserRouter>
      <Route key='getting-started' exact path='/' component={GettingStartedScene} />
      <Route key='team' exact path='/team' component={TeamScene} />
    </BrowserRouter>
  );
}

export default App;
