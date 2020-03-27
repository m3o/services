import React from 'react';
import PageLayout from '../../components/PageLayout';
import EditProfile from '../../components/EditProfile';
import './Profile.scss';

export default class Profile extends React.Component {
  render(): JSX.Element {
    return(
      <PageLayout className='Profile' {...this.props}>
        <EditProfile />
      </PageLayout>
    );
  }
}