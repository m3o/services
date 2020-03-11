import React from 'react';
import { connect } from 'react-redux';
import { PaymentMethod } from '../../api';
import PageLayout from '../../components/PageLayout';
import NewPaymentMethod from './components/NewPaymentMethod';
import PaymentMethodComponent from './components/PaymentMethod';
import './Billing.scss';
import { removePaymentMethod, addPaymentMethod } from '../../store/User';

interface Props {
  stripe?: any;
  elements?: any;

  paymentMethods: PaymentMethod[];
  addPaymentMethod: (pm: PaymentMethod) => void;
  removePaymentMethod: (pm: PaymentMethod) => void;
}

interface State {
  saving: boolean;
  error?: string;
}

class Billing extends React.Component<Props, State> {
  readonly state: State = { saving: false };
  
  setError(error?: string) {
    this.setState({ error, saving: false })
  }

  render():JSX.Element {
    const { paymentMethods } = this.props;
    const { error, saving } = this.state;

    return(
      <PageLayout className='Billing' {...this.props}>
        { this.state.error ? <p>{error}</p> : null }

        <h3>Existing Payment Methods</h3>
        { paymentMethods.map((pm: PaymentMethod) => {
          return <PaymentMethodComponent
                    key={pm.id}
                    paymentMethod={pm}
                    onError={this.setError.bind(this)} 
                    onDelete={this.props.removePaymentMethod} />
        })}

        <NewPaymentMethod
          saving={saving}
          key={paymentMethods.length}
          onError={this.setError.bind(this)}
          onSuccess={this.props.addPaymentMethod}
          onSubmit={() => this.setState({ saving: true })}  />
      </PageLayout>
    );
  }
}

function mapStateToProps(state: any): any {
  return({
    paymentMethods: state.user.user.paymentMethods,
  });
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    addPaymentMethod: (pm: PaymentMethod) => dispatch(addPaymentMethod(pm)),
    removePaymentMethod: (pm: PaymentMethod) => dispatch(removePaymentMethod(pm)),
  });
}

export default connect(mapStateToProps, mapDispatchToProps)(Billing);