// Frameworks
import React from 'react';

// Assets
import './ValidatedInput.scss';


interface Props {
  name: string;
  value: string;
  autoFocus?: boolean;
  placeholder?: string;
  onChange: (name: string, value: string) => void;
  validate?: (value: string) => Promise<string>;
}

interface State {
  status: string;
  error?: string;
}

export default class ValidatedInput extends React.Component<Props, State> {
  readonly state: State = { status: 'pending' };

  async componentDidMount() {
    if(!this.props.validate) {
      this.setState({ status: 'valid' });
      return;
    }

    this.setState({ status: 'pending' });
    if(this.props.value.length === 0) return;

    this.props.validate(this.props.value)
      .then((error) => this.setState({ error, status: error ? 'invalid' : 'valid' }))
      .catch((error) => this.setState({ error, status: 'invalid' }));
  }

  async onChange(e: any) {
    const value = e.target.value;
    this.props.onChange(this.props.name, value);
    if(!this.props.validate) return;

    this.setState({ status: 'pending' });
    if(value.length === 0) return;

    this.props.validate(value)
      .then((error) => this.setState({ error, status: error ? 'invalid' : 'valid' }))
      .catch((error) => this.setState({ error, status: 'invalid' }));
  }

  render(): JSX.Element {
    return <div className='ValidatedInput'>
      <input
        value={this.props.value}
        autoFocus={this.props.autoFocus}
        onChange={this.onChange.bind(this)}
        placeholder={this.props.placeholder} />

      <div className={`dot ${this.state.status}`} />
    </div>
  }
}