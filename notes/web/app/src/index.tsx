import React from 'react';
import ReactDOM from 'react-dom';
import './index.scss';
import App from './App';
import * as serviceWorker from './serviceWorker';

// Declare global window interface so we can mount redux
declare global {
  interface Window {
    __REDUX_DEVTOOLS_EXTENSION__: any;
    store: any;
  }
}

// Mount React to #root element
ReactDOM.render(<App />, document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
