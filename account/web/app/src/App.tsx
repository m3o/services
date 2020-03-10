import React from 'react';
import Call, { User } from './api';
import { ElementsConsumer, CardElement } from '@stripe/react-stripe-js';
import Spinner from './assets/images/spinner.gif'; 
import './App.scss';

interface Props {
  stripe: any;
  elements: any;
}

interface State {
  error?: string;
  user?: User;
  saving: boolean;
}

class App extends React.Component<Props, State> {
  readonly state: State = { saving: false };

  componentDidMount() {
    Call("ReadUser")
      .then(res => this.setState({ user: res.data.user }))
      .catch(err => this.setState({ error: err.message }))
  }

  onChange(e:any) {
    this.setState({ user: new User({
      ...this.state.user,
      [e.target.name]: e.target.value,
    })});
  };

  async onSubmit(e:any) {
    e.preventDefault();
    this.setState({ saving: true });

    const { user } = this.state;
    const { stripe, elements } = this.props;

    // Stripe.js has not yet loaded.
    if (!stripe || !elements) {
      return;
    }

    const result = await stripe.createPaymentMethod({
      type: 'card',
      card: elements.getElement(CardElement),
      billing_details: {
        email: user!.email,
      },
    });

    console.log(result);  

    Call("UpdateUser", { user })
      .then(() => this.setState({ error: '' }))
      .catch(err => this.setState({ error: err.message }))
      .finally(() => this.setState({ saving: false }))
  }

  render(): JSX.Element {
    const { error, user } = this.state;
    if(!user) return this.renderLoading();

    return (
      <div className="App">
        <h1>Account</h1>
        <p className='error'>{error}</p>

        <div className='inner'>
          { this.renderForm() }
          { this.renderBilling() }
        </div>
      </div>
    );
  }

  renderLoading(): JSX.Element {
    return(
      <div className="App">
        <img className='spinner' src={Spinner} alt='Loading' />
      </div>
    );
  }

  renderBilling(): JSX.Element {
    return(
      <div className='stripe'>
        <h3>Billing</h3>
      </div>
    );
  }

  renderForm(): JSX.Element {
    const { user, saving } = this.state;

    return(
      <form onSubmit={this.onSubmit.bind(this)}>
        <h3>Profile</h3>

        <label>First Name</label>
        <input
          type='text'
          name='firstName'
          value={user!.firstName} 
          disabled={this.state.saving}
          onChange={this.onChange.bind(this)} />
        
        <label>Last Name</label>
        <input
          type='text'
          name='lastName'
          value={user!.lastName} 
          disabled={this.state.saving}
          onChange={this.onChange.bind(this)} />
        
        <label>Email</label>
        <input
          name='email'
          type='email'
          value={user!.email}
          disabled={this.state.saving}
          onChange={this.onChange.bind(this)} />
        
        <label>Username</label>
        <input
          name='username'
          type='text'
          value={user!.username}
          disabled={this.state.saving}
          onChange={this.onChange.bind(this)} />

        <label>Card Details</label>
        <CardElement />

        <input disabled={this.state.saving} type='submit' value={ saving ? 'Saving' : 'Save Changes' } />
      </form>
    );
  }
}

export default function InjectedCheckoutForm() {
  return (
    <ElementsConsumer>
      {({stripe, elements}) => (
        <App  stripe={stripe} elements={elements} />
      )}
    </ElementsConsumer>
  );
}
