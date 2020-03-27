import React from 'react';
import Call, { Plan } from '../../../api';
import './Subscribe.scss';

interface Props {
  plans: Plan[];
  onComplete: () => void;
}

interface State {
  saving: boolean;
  selectedPlanID: string;
  error?: string;
}

export default class Subscribe extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { selectedPlanID: props.plans[0]?.id, saving: false };
  }

  setPlan(id: string) {
    if(this.state.saving) return;
    this.setState({ selectedPlanID: id });
  }

  onSubmit(e: any) {
    e.preventDefault();
    this.setState({ saving: true });

    Call("CreateSubscription", { plan_id: this.state.selectedPlanID })
      .then(this.props.onComplete)
      .catch(err => this.setState({ error: err.message }))
  }

  render(): JSX.Element {
    const renderPlan = (p: Plan): JSX.Element => {
      return(
        <div key={p.id} className='row' onClick={() => this.setPlan(p.id)} >
          <input
            type="radio"
            value={p.id}
            name="subscriptions" 
            disabled={this.state.saving}
            checked={p.id === this.state.selectedPlanID}
            onChange={(e: any) => this.setPlan(e.target.value) } />

          <label htmlFor={p.id}>
            <p className='name'>{p.id}</p>
            <p className='price'>${p.amount / 100.0} per {p.interval}</p>
          </label>
        </div>
      )
    }

    return(
      <div className='Subscribe'>
        { this.state.error ? <p className='error'>{this.state.error}</p> : null }

        <form onSubmit={this.onSubmit.bind(this)}>
          { this.props.plans.map(renderPlan) }
          <input disabled={this.state.saving} type="submit" value={this.state.saving ? "Subscribing" : "Subscribe"}/>
        </form>
      </div>
    );
  }
}