import React, { createRef } from 'react';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import Call, { User, Plan } from '../../api';
import EditProfile from '../../components/EditProfile';
import EditPaymentMethods from '../../components/EditPaymentMethods';
import Subscribe from './Subscribe';
import './Signup.scss';

interface Props {
  user: User;
  history: any;
}

interface State {
  stage: number;
  plans?: Plan[];
  loadedPlans: boolean;
  creatingPaymentMethod: boolean;
}

class Signup extends React.Component<Props, State> {
  readonly state: State = { stage: 0, loadedPlans: false, creatingPaymentMethod: false };
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
        this.incrementStage();
        // setup payment methods
        break
      case 2:
        this.incrementStage();
        // setup subscriptions
        break
    }

  }

  render(): JSX.Element {
    if(!this.state.loadedPlans && this.state.stage === 2) return null;

    return(
      <div className='Signup'>
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
          
          <button
            className='continue'
            disabled={this.state.creatingPaymentMethod}
            onClick={this.onSubmitPaymentMethod.bind(this)}>
              Continue →
          </button>
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

  onSubmitPaymentMethod() {
    this.setState({ creatingPaymentMethod: true });

    this.submitNewPaymentMethod.current()
      .then(this.incrementStage.bind(this))
      .finally(() => this.setState({ creatingPaymentMethod: false }));
  }
}

function mapStateToProps(state: any):any {
  return({
    user: state.user.user,
  });
}

export default withRouter(connect(mapStateToProps)(Signup));