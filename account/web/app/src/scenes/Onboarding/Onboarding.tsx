import React from 'react';
import { connect } from 'react-redux';
import { User } from '../../api';
import EditProfile from '../../components/EditProfile';
import './Onboarding.scss';
import EditPaymentMethods from '../../components/EditPaymentMethods';

interface Props {
  user: User;
}

interface State {
  stage: number;
}

class Onboarding extends React.Component<Props, State> {
  readonly state: State = { stage: 0 };

  incrementStage() {
    this.setState({ stage: this.state.stage + 1 });
  }

  componentDidMount() {
    this.autoIncrement();
  }

  componentDidUpdate(prevProps: Props, prevState: State) {
    if(!prevState || prevState.stage === this.state.stage) return;
    this.autoIncrement();
  }

  autoIncrement() {
    switch(this.state.stage) {
      case 0:
        // setup profile
        if(this.props.user.profileCompleted()) {
          this.incrementStage();
          return
        }
      case 1:
        // setup payment methods
        if(this.props.user.paymentMethods.length > 0) {
          this.incrementStage();
        }
    }

  }

  render(): JSX.Element {
    return(
      <div className='Onboarding'>
        <div className='inner'>
          <h1>Welcome to Micro</h1>
          { this.renderStage() }
        </div>
      </div>
    );
  }

  renderStage(): JSX.Element {
    switch(this.state.stage) {
    case 0: 
      return(
        <div className='profile'>
          <p>Let's get started by completing your Micro profile</p>
          <EditProfile onSave={this.incrementStage.bind(this)} />
        </div>
      );
    case 1:
      return(
        <div className='payment-methods'>
          <p>Please enter a payment method</p>
          <EditPaymentMethods />
        </div>
      )
    default:
      return <div />
    }
  }
}

function mapStateToProps(state: any):any {
  return({
    user: state.user.user,
  });
}

export default connect(mapStateToProps)(Onboarding);