import React from 'react';
import Call, { PaymentMethod } from '../../../../api';
import DeleteIcon from '../../../../assets/images/bin.png';
import './PaymentMethod.scss';

interface Props {
  paymentMethod: PaymentMethod;
  onDelete: (pm: PaymentMethod) => void;
  onError: (msg: string) => void;
  setDefault: (pm: PaymentMethod) => void;
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
      .then(() => this.props.onDelete(this.props.paymentMethod))
      .catch(err => {
        this.props.onError('Error deleting payment method: ' + err);
        this.setState({ deleting: false });
      });
  }

  async onMakeDefault() {
    const pm = this.props.paymentMethod;

    // eslint-disable-next-line no-restricted-globals
    if(!confirm(`Are you sure you want to make ${pm.card_brand} ending in ${pm.card_last_4} your default payment method?`)) return;

    Call("DefaultPaymentMethod", { id: pm.id })
      .then(() => this.props.setDefault(pm))
      .catch(console.warn);
  }

  render():JSX.Element {
    const pm = this.props.paymentMethod;

    return(
      <div className='PaymentMethod'>
        <div className='pm-left'>
          <p><span>{pm.card_brand}</span> ending in {pm.card_last_4}</p>
          <p>Exp: {pm.card_exp_month}/{pm.card_exp_year}</p>
        </div>

        <div className='pm-right'>
          { this.props.paymentMethod.default ? <p>Default</p> : <p className='make-default' onClick={this.onMakeDefault.bind(this)}>Make Default</p> }
          <img src={DeleteIcon} alt='Delete Payment Method' onClick={this.onDelete.bind(this)} />
        </div>
      </div>
    );
  }
}
