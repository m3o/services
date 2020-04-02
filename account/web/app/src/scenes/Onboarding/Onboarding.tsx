import React, { createRef } from 'react';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import Call, { User, Plan } from '../../api';
import EditProfile from '../../components/EditProfile';
import EditPaymentMethods from '../../components/EditPaymentMethods';
import Subscribe from './Subscribe';
import './Onboarding.scss';

interface Props {
  user: User;
  history: any;
}

interface State {
  stage: number;
  plans?: Plan[];
  loadedPlans: boolean;
}

class Onboarding extends React.Component<Props, State> {
  readonly state: State = { stage: 0, loadedPlans: false };
  submitNewPaymentMethod: React.RefObject<() => Promise<any>>;

  constructor(props: Props) {
    super(props);
    this.submitNewPaymentMethod = createRef();
  }

  incrementStage() {
    this.setState({ stage: this.state.stage + 1 });
  }

  componentDidMount() {
    this.autoIncrement();

    Call("ListPlans")
      .then(res => {
        const plans = (res.data.plans || []).map(p => new Plan(p));
        this.setState({ plans: plans.sort((a,b) => a.amount - b.amount) });
      })
      .finally(() => this.setState({ loadedPlans: true }))
      .catch(console.warn);
  }

  componentDidUpdate(prevProps: Props, prevState: State) {
    if(!prevState || prevState.stage === this.state.stage) return;
    if(this.state.stage === 3) this.props.history.push('/');
    this.autoIncrement();
  }

  autoIncrement() {
    switch(this.state.stage) {
      case 0:
        // setup profile
        if(this.props.user.profileCompleted()) {
          this.incrementStage();
        }
        break
      case 1:
        // setup payment methods
        if(this.props.user.payment_methods.length > 0) {
          this.incrementStage();
        }
        break
      case 2:
        // setup subscription
        if(this.props.user.subscriptions.length > 0) {
          this.incrementStage();
        }
        break
    }

  }

  render(): JSX.Element {
    if(!this.state.loadedPlans && this.state.stage === 2) return null;

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
          <EditProfile buttonText='Continue →' onSave={this.incrementStage.bind(this)} />
        </div>
      );
    case 1:
      return(
        <div className='payment-methods'>
          <p>Please enter a payment method</p>
          <EditPaymentMethods singleCardMode={true} submitNewPaymentMethod={this.submitNewPaymentMethod} />
          <button onClick={() => this.submitNewPaymentMethod.current().then(this.incrementStage.bind(this))} className='continue'>Continue →</button>
        </div>
      )
    default:
      return(
        <div className='subscription'>
          <p>Please select a subscription</p>
          <Subscribe onComplete={this.incrementStage.bind(this)} plans={this.state.plans!} />
        </div>
      );
    }
  }
}

function mapStateToProps(state: any):any {
  return({
    user: state.user.user,
  });
}

export default withRouter(connect(mapStateToProps)(Onboarding));