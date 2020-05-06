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
  loading: boolean;
  error?: string;
  timer?: any;
}

export default class ValidatedInput extends React.Component<Props, State> {
  readonly state: State = { loading: false };

  async componentDidMount() {
    if(!this.props.validate) return;
    if(this.props.value.length === 0) return;
    this.validate();
  }

  validate() {
    if(this.props.value.length === 0) {
      this.setState({ error: undefined, loading: false });
      return;
    }

    this.props.validate(this.props.value)
      .then((error) => this.setState({ error, loading: false }))
      .catch((error) => this.setState({ error, loading: false }));
  }

  async onChange(e: any) {
    const value = e.target.value;
    this.props.onChange(this.props.name, value);

    if(this.state.timer) clearTimeout(this.state.timer);
    this.setState({ timer: setTimeout(this.validate.bind(this), 500), loading: true });
  }

  render(): JSX.Element {
    let status = '';
    if(!this.props.validate) {
      status = 'valid';
    } else if(this.state.loading) {
      status = 'pending';
    } else if(this.state.error || this.props.value.length === 0) {
      status = 'invalid';
    } else  {
      status = 'valid';
    }

    return <div className='ValidatedInput'>
      <input
        value={this.props.value}
        autoFocus={this.props.autoFocus}
        onChange={this.onChange.bind(this)}
        placeholder={this.props.placeholder} />

      { this.state.error ? <p className='error'>{this.state.error}</p> : null }
      <div className={`dot ${status}`} />
    </div>
  }
}