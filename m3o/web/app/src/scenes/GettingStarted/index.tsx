import React from 'react';
import PageLayout from '../../components/PageLayout';
import ContinuousDeployments from './assets/deployment.png';
import Configuration from './assets/configuration.png';
import Domain from './assets/domain.png';
import './style.scss';

interface Tutorial {
  key: string;
  title: string;
  description: string;
  duration: string;
  stepsRemaining: number;
  icon: string;
}

const tutorials: Tutorial[] = [
  {
    key: 'continuous-deployent',
    title: 'Setup Continuous Deployments',
    description: 'Configure Micro to auto deploy your services whenever a commit is added to your master branch. Simply drop in the GitHub action, configure it with your secret, then sit back and watch your services deploy.',
    duration: 'About 5 mins',
    stepsRemaining: 0,
    icon: ContinuousDeployments,
  },
  {
    key: 'configuration',
    title: 'Add Secrets and Configuration',
    description: 'Provide your services with the configuration they need to operate. Configuration can be specified at a global or service level and used RBAC to ensure services can only access the secrets they require.',
    duration: 'About 10 mins',
    stepsRemaining: 5,
    icon: Configuration,
  },
  {
    key: 'domain',
    title: 'Configure a Custom Domain',
    description: 'Configure your domain to expose Micro Web & API. Custom domains come free of charge in M3O.',
    duration: 'About 5 mins',
    stepsRemaining: 2,
    icon: Domain,
  },
]

export default class GettingStartedScene extends React.Component {
  render(): JSX.Element {
    return(
      <PageLayout className='GettingStarted'>
        <h1>Quick start guide</h1>
        { tutorials.map(this.renderTutorial) }
      </PageLayout>
    );
  }

  renderTutorial(section: Tutorial): JSX.Element {
    return(
      <section key={section.key} className={section.stepsRemaining > 0 ? '' : 'complete'}>
          <img src={section.icon} alt={section.title} />
          <div className='section-main'>
            <div className='section-upper'>
              <h4>{section.title}</h4>
              <p>{section.duration}<span>•</span>{section.stepsRemaining > 0 ? `${section.stepsRemaining} steps remaining` : 'Complete'}</p>
            </div>

            <div className='section-lower'>
              <p>{section.description}</p>
            </div>
          </div>
        </section>
    )
  }
}