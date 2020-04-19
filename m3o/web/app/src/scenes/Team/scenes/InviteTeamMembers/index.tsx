import React from 'react';
import PageLayout from '../../../../components/PageLayout';
import BinIcon from './assets/bin.png';
import './style.scss';

interface Props {
  history: any;
}

interface State {
  emails: string[];
  currentIndex: number;
}

export default class InviteTeamMembersScene extends React.Component<Props, State> {
  readonly state: State = { emails: [], currentIndex: 0 };

  render(): JSX.Element {
    return(
      <PageLayout className='InviteTeamMembers'>
        <header>
          <h1>Invite Team Members</h1>

          <button className='btn danger' onClick={this.onCancel.bind(this)}>
            <p>Cancel</p>
          </button>

          <button className='btn' onClick={this.onSave.bind(this)}>
            <p>Send Invites</p>
          </button>
        </header>

        <p className='instructions'>Enter the emails of your team members below and they will recieve an email containing a signup link which is valid for 24 hours.</p>

        <form onSubmit={(e: any) => {e.preventDefault(); this.onSave()}}>
          { this.state.emails.map(this.renderRow.bind(this)) }
          { this.renderRow('', this.state.emails.length)}
        </form>
      </PageLayout>
    );
  }

  renderRow(email: string, index: number, exists?: boolean): JSX.Element {
    const onChange = (e: any) => {
      let emails = [...this.state.emails];
      emails[index] = e.target.value;
      this.setState({ emails, currentIndex: index });
    };

    const onDelete = () => {
      const emails = this.state.emails.filter((e, i) => i !== index);
      this.setState({ emails });
    }

    return (
      <div className='row' key={index}>
        <input
          value={email}
          onChange={onChange}
          autoFocus={this.state.currentIndex === index}
          placeholder={exists ? 'Email Address' : 'Add an email'} />

        { exists ? <img src={BinIcon} alt='Remove Icon' onClick={onDelete} /> : null }
      </div>
    );
  }

  onSave(): void {
    // eslint-disable-next-line no-restricted-globals
    const c = this.state.emails.length
    alert(`Invites have been sent to ${c} email${c > 1 ? 's' : ''}`);
    this.props.history.push('/team');
  }

  onCancel(): void {
    // eslint-disable-next-line no-restricted-globals
    if (!confirm(`Are you sure you want to cancel? No invites will be sent.`)) return;
    this.props.history.push('/team');
  }
}