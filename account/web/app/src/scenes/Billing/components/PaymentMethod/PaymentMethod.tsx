import React from 'react';
import Call, { PaymentMethod } from '../../../../api';
import DeleteIcon from '../../../../assets/images/bin.png';
import './PaymentMethod.scss';

interface Props {
  paymentMethod: PaymentMethod;
  onDelete: (id: string) => void;
  onError: (msg: string) => void;
}

interface State {
  deleting: boolean;
}

export default class PaymentMethodComponent extends React.Component<Props, State> {
  readonly state: State = { deleting: false };

  async onDelete() {
    // eslint-disable-next-line no-restricted-globals
    if(!confirm('Are you sure you want to delete this payment method?')) return;

    this.setState({ deleting: true });

    const { id } = this.props.paymentMethod;
    Call("DeletePaymentMethod", { id })
      .then(() => this.props.onDelete(id))
      .catch(err => {
        this.props.onError('Error deleting payment method: ' + err);
        this.setState({ deleting: false });
      });
  }

  render():JSX.Element {
    const pm = this.props.paymentMethod;

    return(
      <div className='PaymentMethod'>
        <div className='pm-left'>
          <p><span>{pm.cardBrand}</span> ending in {pm.cardLast4}</p>
          <p>Exp: {pm.cardExpMonth}/{pm.cardExpYear}</p>
        </div>

        <div className='pm-right' onClick={this.onDelete.bind(this)}>
          <img src={DeleteIcon} alt='Delete Payment Method' />
        </div>
      </div>
    );
  }
}