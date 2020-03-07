import React from 'react';
import { createStore } from 'redux';
import { Provider } from 'react-redux';
import { BrowserRouter , Route } from 'react-router-dom';

// Scenes 
import { rootReducer } from './store';
import HomeScene from './scenes/Home';
import NotesScene from './scenes/Notes';
import SprintsScene from './scenes/Sprints';

// Redux
window.store = createStore(
  rootReducer,
  window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__()
);

export default class App extends React.Component {
  render():JSX.Element {
    return(
      <Provider store={window.store} basename='/distributed'>
        <BrowserRouter>
          <div className='App'>
            <Route exact path='/distributed/' component={HomeScene}/>
            <Route exact path='/distributed/notes' component={NotesScene}/>
            <Route exact path='/distributed/notes/:id' component={NotesScene}/>
            <Route exact path='/distributed/notes/:id/:options' component={NotesScene}/>
            <Route exact path='/distributed/sprints' component={SprintsScene}/>
          </div>
        </BrowserRouter>
      </Provider>
    );
  }
}