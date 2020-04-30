import React from 'react';
import './Dashboard.scss';
import PageLayout from '../../components/PageLayout';

export default class Dashboard extends React.Component {
  render(): JSX.Element {
    return <PageLayout className='Dashboard'>
      <h1>Foo</h1>
    </PageLayout>
  }
}