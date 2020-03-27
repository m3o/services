import React from 'react';
import PageLayout from '../../components/PageLayout';
import EditPaymentMethods from '../../components/EditPaymentMethods';
import './Billing.scss';

export default class Billing extends React.Component {
  render():JSX.Element {
    return(
      <PageLayout className='Billing' {...this.props}>
        <EditPaymentMethods />
      </PageLayout>
    );
  }
}