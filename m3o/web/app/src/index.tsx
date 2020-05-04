import React from 'react';
import ReactDOM from 'react-dom';
import * as serviceWorker from './serviceWorker';
import App from './App';
import './index.scss';
import { Elements } from '@stripe/react-stripe-js';
import { loadStripe } from '@stripe/stripe-js';

// Redux
import { Provider } from 'react-redux';
import store from './store';

// Redux Setup
window.store = store; 

// Stripe
const stripePromise = loadStripe('pk_test_wuI8wlKwKBUZ9iHnYlQPa8BH');

// Declare global window interface so we can mount redux
declare global {
  interface Window {
    __REDUX_DEVTOOLS_EXTENSION__: any;
    store: any;
  }
}

ReactDOM.render(
  <React.StrictMode>
    <Provider store={window.store}>
      <Elements stripe={stripePromise}>
        <App />
      </Elements>
    </Provider>
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
