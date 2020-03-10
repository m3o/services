import React from 'react';
import PageLayout from '../../components/PageLayout';
import logo from '../../assets/images/logo.svg';
import './Home.scss';

export default class HomeScene extends React.Component {
  render():JSX.Element {
    return(
      <PageLayout className='HomeScene' {...this.props}>
        <img src={logo} className="logo" alt="Distributed Logo" />
        <h1>Welcome to Distributed</h1>
      </PageLayout>
    );
  }
}