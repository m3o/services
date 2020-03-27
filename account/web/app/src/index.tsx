import React from 'react';
import ReactDOM from 'react-dom';
import './index.scss';
import App from './App';
import * as serviceWorker from './serviceWorker';
import { Provider } from 'react-redux';
import { BrowserRouter } from 'react-router-dom';
import { createStore, combineReducers } from 'redux';
import UserReducer from './store/User';
import RedirectReducer from './store/Redirect';
import { Elements } from '@stripe/react-stripe-js';
import { loadStripe } from '@stripe/stripe-js';

// Redux Setup
window.store = createStore(combineReducers({
  user: UserReducer,
  redirect: RedirectReducer,
})); 

// Declare global window interface so we can mount redux
declare global {
  interface Window {
    __REDUX_DEVTOOLS_EXTENSION__: any;
    store: any;
  }
}

// Stripe
const stripePromise = loadStripe('pk_test_wuI8wlKwKBUZ9iHnYlQPa8BH');

// Wrap the app
const WrappedApp = ():JSX.Element => {
  return(
    <Provider store={window.store} >
      <Elements stripe={stripePromise}>
        <BrowserRouter>
          <App />
        </BrowserRouter>
      </Elements>
    </Provider>
  );
}

ReactDOM.render(<WrappedApp />, document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
