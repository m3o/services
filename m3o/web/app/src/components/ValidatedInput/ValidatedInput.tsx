// Frameworks
import React from 'react';

// Assets
import './ValidatedInput.scss';


interface Props {
  name?: string;
  value: string;
  disabled?: boolean;
  autoFocus?: boolean;
  placeholder?: string;
  validateDelay?: number;

  validate?: (value: string) => Promise<void>;
  onChange?: (name: string, value: string) => void;
  onValid?: () => void;
  onInvalid?: () => void;
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
    // do a basic check here to check for a missing value
    if(this.props.value.length === 0) {
      this.setState({ error: undefined, loading: false });
      if(this.props.onInvalid) this.props.onInvalid();
      return;
    }

    this.props.validate(this.props.value)
      .then(() => {
        this.setState({ error: undefined, loading: false });
        if(this.props.onValid) this.props.onValid();
      })
      .catch((error) => {
        this.setState({ error, loading: false })
        if(this.props.onInvalid) this.props.onInvalid();
      });
  }

  async onChange(e: any) {
    const value = e.target.value;
    this.props.onChange(this.props.name, value);

    if(!this.props.validate) return;
    
    const delay = this.props.validateDelay || 500;
    if(this.state.timer) clearTimeout(this.state.timer);
    this.setState({ timer: setTimeout(this.validate.bind(this), delay), loading: true });
  }

  render(): JSX.Element {
    let status = '';
    if(this.props.disabled) {
      status = 'disabled';
    } else if(!this.props.validate) {
      status = 'valid';
    } else if(this.state.loading) {
      status = 'loading';
    } else if(this.state.error) {
      status = 'invalid';
    } else if(this.props.value.length === 0) {
      status = 'pending';
    } else {
      status = 'valid';
    }

    return <div className='ValidatedInput'>
      <input
        value={this.props.value}
        disabled={this.props.disabled}
        autoFocus={this.props.autoFocus}
        onChange={this.onChange.bind(this)}
        placeholder={this.props.placeholder} />

      { this.state.error ? <p className={`error ${status}`}>{this.state.error}</p> : null }
      <div className={`dot ${status}`} />
    </div>
  }
}