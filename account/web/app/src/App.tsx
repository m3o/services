import React from 'react';
import Call, { User } from './api';

// Scenes
import Profile from './scenes/Profile';
import Billing from './scenes/Billing';

// Assets
import Spinner from './assets/images/spinner.gif'; 
import './App.scss';

interface Props {
}

interface State {
  user?: User;
}

export default class App extends React.Component<Props, State> {
  readonly state: State = {};

  componentDidMount() {
    Call("ReadUser").then(res => this.setState({ user: res.data.user }))
  }

  componentDidUpdate(prevProps: Props, prevState: State) {
    if(this.state !== prevState) console.log("State: ", this.state);
    if(this.props !== prevProps) console.log("Props: ", this.props);
  }

  render(): JSX.Element {
    const { user } = this.state;
    if(!user) return this.renderLoading();

    return (
      <div className="App">
        <h1>Account</h1>

        <div className='inner'>
          <Profile user={user} />
          <Billing paymentMethods={user.paymentMethods} deletePaymentMethod={this.deletePaymentMethod.bind(this)} />
        </div>
      </div>
    );
  }

  renderLoading(): JSX.Element {
    return(
      <div className="App">
        <img className='spinner' src={Spinner} alt='Loading' />
      </div>
    );
  }

  deletePaymentMethod(id: string) {
    this.setState({ 
      user: new User({
        ...this.state.user,
        paymentMethods: this.state.user!.paymentMethods.filter(pm => pm.id !== id),
      }),
    });
  }
}