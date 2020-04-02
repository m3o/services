import React from 'react';
import Call, { PaymentMethod } from '../../../../api';
import { CardElement, useStripe, useElements } from '@stripe/react-stripe-js';
import PlusIcon from './plus.png';
import './NewPaymentMethod.scss';

interface Props {
  onSubmit: () => void,
  onSuccess: (pm: PaymentMethod) => void,
  onError: (msg: string) => void,
  submitRef: any;
  saving: boolean,
  hideButton: boolean,
}

export default ({ onSuccess, onError, onSubmit, saving, submitRef, hideButton }: Props) => {
  const stripe = useStripe();
  const elements = useElements();

  const onFormSubmit = async (event?: any): Promise<any> => {
    return new Promise(async (resolve, reject) => {
    if(event) event.preventDefault();
    onSubmit();
    
    // Ensure stripe has loaded
    if (!stripe || !elements) {
      reject();
      return;
    }

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
      reject();
      return;
    }

    // Submit to the API
      Call("CreatePaymentMethod", { id: paymentMethod!.id })
        .then(res => {
          resolve();
          onSuccess(res.data.payment_method);
        })
        .catch(err => {
          reject();
          onError(err.message);
        });
    })
  }

  // this is a hack to allow us to trigger the callback in a parent
  // component, this is used in onboarding.
  submitRef.current = onFormSubmit;

  return(
    <form className='NewPaymentMethod' onSubmit={onFormSubmit}>
      <label>New Payment Method</label>

      <div className='payment-method'>
        <CardElement />
        
        { hideButton ? null : <button disabled={saving} onClick={onFormSubmit}>
          <img src={PlusIcon} alt='Add payment method' />
        </button> }
      </div>
    </form>
  );
}
