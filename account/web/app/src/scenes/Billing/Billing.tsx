import React from 'react';
import { connect } from 'react-redux';
import { PaymentMethod } from '../../api';
import PageLayout from '../../components/PageLayout';
import NewPaymentMethod from './components/NewPaymentMethod';
import PaymentMethodComponent from './components/PaymentMethod';
import './Billing.scss';

interface Props {
  stripe?: any;
  elements?: any;

  paymentMethods: PaymentMethod[];
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
                    onError={this.setError.bind(this)} />
        })}

        <NewPaymentMethod
          saving={saving}
          onSuccess={console.log}
          key={paymentMethods.length}
          onError={this.setError.bind(this)}
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
    // paymentMethods: state.user.user.paymentMethods,
  });
}

export default connect(mapStateToProps, mapDispatchToProps)(Billing);