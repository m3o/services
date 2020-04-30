import React from 'react';
import './Dashboard.scss';
import PageLayout from '../../components/PageLayout';
import ProjectIcon from './assets/project.png';
import AddIcon from './assets/add.png';
import NotificationsIcon from './assets/notifications.png';
import FeedbackIcon from './assets/feedback.png';
import DocsIcon from './assets/docs.png';

export default class Dashboard extends React.Component {
  render(): JSX.Element {
    return <PageLayout className='Dashboard'>
      <div className='left'>
        <section>
          <a href=''>
            <img src={NotificationsIcon} alt='Notifications' />
            <p>Notifications</p>
          </a>

          <a href=''>
            <img src={FeedbackIcon} alt='Feedback' />
            <p>Feedback</p>
          </a>

          <a href=''>
            <img src={DocsIcon} alt='Docs' />
            <p>Docs</p>
          </a>
        </section>

        <section>
          <p>ben-toogood</p>
          
          <a href=''>
            <img src={ProjectIcon} alt='ben-toogood/hello-world' />
            <p>ben-toogood/hello-world</p>
          </a>          

          <a href=''>
            <img src={AddIcon} alt='New Project' />
            <p>New Project</p>
          </a>
        </section>

        <section>
          <p>Kytra</p>
          
          <a href='' className='active'>
            <img src={ProjectIcon} alt='kytra/production' />
            <p>kytra/production</p>
          </a>
          
          <a href=''>
            <img src={ProjectIcon} alt='kytra/staging' />
            <p>kytra/staging</p>
          </a>
          
          <a href=''>
            <img src={ProjectIcon} alt='kytra/develpment' />
            <p>kytra/develpment</p>
          </a>
          

          <a href=''>
            <img src={AddIcon} alt='New Project' />
            <p>New Project</p>
          </a>
        </section>

        <section>
          <p>Micro</p>
          
          <a href='' >
            <img src={ProjectIcon} alt='micro/services' />
            <p>micro/services</p>
          </a>
          
          <a href=''>
            <img src={ProjectIcon} alt='micro/m3o' />
            <p>micro/m3o</p>
          </a>
          
          <a href=''>
            <img src={AddIcon} alt='New Project' />
            <p>New Project</p>
          </a>
        </section>
      </div>

      <div className='right'>

      </div>
    </PageLayout>
  }
}