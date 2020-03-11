import React from 'react';
import Call, { User, PaymentMethod } from './api';
import { ElementsConsumer, CardElement } from '@stripe/react-stripe-js';

// Components
import PaymentMethodComponent from './components/PaymentMethod';

// Assets
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
  savingPaymentMethod: boolean;
}

class App extends React.Component<Props, State> {
  readonly state: State = { saving: false, savingPaymentMethod: false };

  componentDidMount() {
    Call("ReadUser")
      .then(res => this.setState({ user: res.data.user }))
      .catch(err => this.setError(err.message));
  }

  componentDidUpdate(prevProps: Props, prevState: State) {
    if(this.state !== prevState) console.log("State: ", this.state);
    if(this.props !== prevProps) console.log("Props: ", this.props);
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

    Call("UpdateUser", { user })
      .then(() => this.setError(''))
      .catch(err => this.setError(err.message))
      .finally(() => this.setState({ saving: false }));
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
          { this.renderPaymentMethods() }
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

  async onPaymentMethodSubmit(event: any) {
    event.preventDefault();
    this.setState({ error: '', savingPaymentMethod: true });

    // Ensure stripe has loaded
    const { stripe, elements } = this.props;
    if (!stripe || !elements) return;

    // Get the card element from the dom
    const cardElement = elements.getElement(CardElement);

    // Create the card in the stripe api 
    const { error, paymentMethod } = await stripe.createPaymentMethod({
      type: 'card',
      card: cardElement,
    });

    // Handle the error
    if (error) {
      this.setError(error);
      return
    }

    // Submit to the API
    Call("CreatePaymentMethod", { id: paymentMethod.id })
      .catch((error) => this.setState({ error, savingPaymentMethod: false }))
      .then(res => {
        const user = new User({
          ...this.state.user,
          paymentMethods:[
            ...this.state.user!.paymentMethods,
            res.data.paymentMethod,
          ]
        });

        this.setState({ user, savingPaymentMethod: false });
      });
  }

  setError(error: string) {
    this.setState({ error });
  }

  deletePaymentMethod(id: string) {
    this.setState({ 
      user: new User({
        ...this.state.user,
        paymentMethods: this.state.user!.paymentMethods.filter(pm => pm.id !== id),
      }),
    });
  }

  renderPaymentMethods(): JSX.Element {
    const saving = this.state.savingPaymentMethod;
    const { paymentMethods } = this.state.user!;
    
    return(
      <div className='stripe'>
        <h3>Payment Methods</h3>

        { paymentMethods.map((pm: PaymentMethod) => {
          return <PaymentMethodComponent
                    key={pm.id}
                    paymentMethod={pm}
                    onError={this.setError.bind(this)}
                    onDelete={this.deletePaymentMethod.bind(this)} />
        })}

        <form onSubmit={this.onPaymentMethodSubmit.bind(this)}>
          <label>New Payment Method</label>
          <CardElement key={paymentMethods.length} />
          <input disabled={saving} type='submit' value={ saving ? 'Saving' : 'Create Payment Method' } />
        </form>
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
