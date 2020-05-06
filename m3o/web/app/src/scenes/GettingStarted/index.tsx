// Frameworks
import React from 'react';

// Components
import PageLayout from '../../components/PageLayout';

// Assets
import M3OIcon from './assets/tri.png';
import './GettingStarted.scss';

interface Props {
  history: any;
}

export default class GettingStarted extends React.Component<Props> {
  render(): JSX.Element {
    return (
      <PageLayout className='GettingStarted'>
        <div className='center'>
          <h1>Getting Started</h1>
          <p>Welcome to M3O! Let's get started by creating your first project</p>
          <button onClick={() => this.props.history.push('/new/project')} className='btn'>Create a project</button>
          <img src={M3OIcon} alt='M3O' />
        </div>
      </PageLayout>
    );
  }
}