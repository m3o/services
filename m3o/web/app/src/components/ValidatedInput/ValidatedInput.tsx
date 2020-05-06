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
  validate: ((value: string) => string) | ((value: string) => Promise<string>);
}

interface State {
  status: string;
}

export default class ValidatedInput extends React.Component<Props, State> {
  readonly state: State = { status: 'pending' };

  async componentDidMount() {
    const status = await this.props.validate(this.props.value);
    this.setState({ status });
  }

  async onChange(e: any) {
    const value = e.target.value;
    this.props.onChange(this.props.name, value);

    const status = await this.props.validate(value);
    this.setState({ status });
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