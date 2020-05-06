// Frameworks
import React from 'react';
import { connect } from 'react-redux';

// Utils
import { State as GlobalState } from '../../store';

// Components
import PageLayout from '../../components/PageLayout';

// Styling
import './Billing.scss';

interface Props {}

class Billing extends React.Component<Props> {
  render(): JSX.Element {
    return (
      <PageLayout className='Billing' hideSidebar>
        <div className='center'>
          <h1>Billing</h1>
        </div>
      </PageLayout>
    );
  }
}

function mapStateToProps(state: GlobalState): any {
  return ({});
}

export default connect(mapStateToProps)(Billing);