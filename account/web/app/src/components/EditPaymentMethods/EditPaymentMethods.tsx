import React from 'react';
import { connect } from 'react-redux';
import { PaymentMethod } from '../../api';
import NewPaymentMethod from './components/NewPaymentMethod';
import PaymentMethodComponent from './components/PaymentMethod';
import './EditPaymentMethods.scss';
import { setDefaultPaymentMethod, addPaymentMethod, removePaymentMethod } from '../../store/User';

interface Props {
  stripe?: any;
  elements?: any;
  singleCardMode?: boolean;
  submitNewPaymentMethod: React.RefObject<() => Promise<any>>;

  paymentMethods: PaymentMethod[];
  setDefault: (pm: PaymentMethod) => void;
  addPaymentMethod: (pm: PaymentMethod) => void;
  removePaymentMethod: (pm: PaymentMethod) => void;
}

interface State {
  saving: boolean;
  error?: string;
}

class EditPaymentMethods extends React.Component<Props, State> {
  readonly state: State = { saving: false };

  setError(error?: string) {
    this.setState({ error, saving: false })
  }

  addPaymentMethod(pm: PaymentMethod) {
    this.props.addPaymentMethod(pm);
    this.setState({ error: undefined, saving: false })
  }

  render():JSX.Element {
    const { paymentMethods, singleCardMode } = this.props;
    const { error, saving } = this.state;

    return(
      <div className='EditPaymentMethods'>
        { this.state.error ? <p className='error'>{error}</p> : null }

        { singleCardMode ? null : this.renderPaymentMethods() }

        <NewPaymentMethod
          saving={saving}
          key={paymentMethods.length}
          hideButton={singleCardMode}
          onError={this.setError.bind(this)}
          onSuccess={this.addPaymentMethod.bind(this)}
          submitRef={this.props.submitNewPaymentMethod}
          onSubmit={() => this.setState({ saving: true })} />
      </div>
    );
  }

  renderPaymentMethods(): JSX.Element {
    if(this.props.paymentMethods.length === 0) return null;
    
    return(
      <div className='existing'>
        <h3>Existing Payment Methods</h3>
        { this.props.paymentMethods.map((pm: PaymentMethod) => {
          return <PaymentMethodComponent
                    key={pm.id}
                    paymentMethod={pm}
                    onError={this.setError.bind(this)} 
                    setDefault={this.props.setDefault}
                    onDelete={this.props.removePaymentMethod} />
        })}
      </div>
    )
  }
}

function mapStateToProps(state: any): any {
  return({
    paymentMethods: state.user.user.payment_methods,
  });
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    setDefault: (pm: PaymentMethod) => dispatch(setDefaultPaymentMethod(pm)),
    addPaymentMethod: (pm: PaymentMethod) => dispatch(addPaymentMethod(pm)),
    removePaymentMethod: (pm: PaymentMethod) => dispatch(removePaymentMethod(pm)),
  });
}

export default connect(mapStateToProps, mapDispatchToProps)(EditPaymentMethods);