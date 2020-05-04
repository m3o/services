import React from 'react';
import { CardElement, useStripe, useElements } from '@stripe/react-stripe-js';
import './PaymentMethod.scss';

interface Props {
  status: string;
  onSubmit: () => void;
  onComplete: (id: string) => void;
  onError: (err: string) => void;
}

export default (props: Props) => {
  const stripe = useStripe();
  const elements = useElements();

  const handleSubmit = async (event) => {
    event.preventDefault();
    props.onSubmit();

    const {error, paymentMethod} = await stripe.createPaymentMethod({
      type: 'card',
      card: elements.getElement(CardElement),
    });

    if(error) {
      props.onError(error.message);
      return
    }

    props.onComplete(paymentMethod.id);
  };

  return (
    <div className='PaymentMethod'>
      { props.status ? <p className='status'>{props.status}</p> :  null }

      <form onSubmit={handleSubmit}>
        <CardElement />

        <button className='btn' type="submit" disabled={!stripe}>
          Setup Subsciption
        </button>
      </form>
    </div>
  );
};