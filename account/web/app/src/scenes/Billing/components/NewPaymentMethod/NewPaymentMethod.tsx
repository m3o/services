import React from 'react';
import Call, { PaymentMethod } from '../../../../api';
import { CardElement, useStripe, useElements } from '@stripe/react-stripe-js';
import './NewPaymentMethod.scss';

interface Props {
  onSubmit: () => void,
  onSuccess: (pm: PaymentMethod) => void,
  onError: (msg: string) => void,
  saving: boolean,
}

export default ({ onSuccess, onError, onSubmit, saving }: Props) => {
  const stripe = useStripe();
  const elements = useElements();

  const onFormSubmit = async (event: any) => {
    event.preventDefault();
    onSubmit();

    // Ensure stripe has loaded
    if (!stripe || !elements) return;

    // Get the card element from the dom
    const cardElement = elements.getElement(CardElement);

    // Create the card in the stripe api 
    const { error, paymentMethod } = await stripe.createPaymentMethod({
      type: 'card',
      card: cardElement!,
    });

    // Handle the error
    if (error) {
      onError(error.message!);
      return;
    }

    // Submit to the API
    Call("CreatePaymentMethod", { id: paymentMethod!.id })
    .then(res => onSuccess(res.data.paymentMethod))
    .catch(onError);
  }

  return(
    <form className='NewPaymentMethod' onSubmit={onFormSubmit}>
      <label>New Payment Method</label>
      <CardElement />
      <input type='submit' disabled={saving} value={saving ? 'Creating' : 'Create Payment Method'} />
    </form>
  );
}
