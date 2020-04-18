import React from 'react';
import { BrowserRouter, Route } from 'react-router-dom';
import './App.scss';
import GettingStartedScene from './scenes/GettingStarted';

function App() {
  return (
    <BrowserRouter>
      <Route key='getting-started' exact path='/' component={GettingStartedScene} />
    </BrowserRouter>
  );
}

export default App;
