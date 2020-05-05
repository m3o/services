// Frameworks
import React from 'react';

// Components
import PageLayout from '../../components/PageLayout';

export default class NotFound extends React.Component {
  render(): JSX.Element {
    return (
      <PageLayout>
        <div className='center'>
          <h1>Not Found</h1>
        </div>
      </PageLayout>
    );
  }
}