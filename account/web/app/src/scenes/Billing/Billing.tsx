import React from 'react';
import Call, { PaymentMethod } from '../../api';
import PaymentMethodComponent from './components/PaymentMethod';
import { ElementsConsumer, CardElement } from '@stripe/react-stripe-js';
import './Billing.scss';

interface Props {
  stripe?: any;
  elements?: any;

  paymentMethods: PaymentMethod[];
  deletePaymentMethod: (id: string) => void;
}

interface State {
  saving: boolean;
  error?: string;
}

class Billing extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { saving: true };
  }

  setError(error?: string) {
    this.setState({ error })
  }


  async onPaymentMethodSubmit(event: any) {
    event.preventDefault();
    this.setState({ error: undefined, saving: true });

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
      .catch((error) => this.setState({ error, saving: false }))
      // .then(res => {
      //   const user = new User({
      //     ...this.state.user,
      //     paymentMethods:[
      //       ...this.state.user!.paymentMethods,
      //       res.data.paymentMethod,
      //     ]
      //   });

      //   this.setState({ user, savingPaymentMethod: false });
      // });
  }

  render():JSX.Element {
    const { paymentMethods, deletePaymentMethod } = this.props;
    const { error, saving } = this.state;

    return(
      <div className='stripe'>
        <h3>Payment Methods</h3>
        { this.state.error ? <p>{error}</p> : null }

        { paymentMethods.map((pm: PaymentMethod) => {
          return <PaymentMethodComponent
                    key={pm.id}
                    paymentMethod={pm}
                    onError={this.setError.bind(this)}
                    onDelete={deletePaymentMethod} />
        })}

        <form onSubmit={this.onPaymentMethodSubmit.bind(this)}>
          <label>New Payment Method</label>
          <CardElement key={paymentMethods.length} />
          <input disabled={saving} type='submit' value={ saving ? 'Saving' : 'Create Payment Method' } />
        </form>
      </div>
    );
  }
}

export default function InjectedCheckoutForm(props: Props) {
  return (
    <ElementsConsumer>
      {({stripe, elements}) => (
        <Billing stripe={stripe} elements={elements} {...props} />
      )}
    </ElementsConsumer>
  );
}
