import React from 'react';

// Components
import PageLayout from '../../components/PageLayout';
import Messenger from '../../components/Messenger';
import SprintRow from './components/Row';

// Assets
import Arrow from '../../assets/images/arrow.png';
import ChatIcon from '../../assets/images/chat-icon-white.png';
import AddIcon from '../../assets/images/add-icon.png';
import './Sprints.scss';

interface Props {}

interface State {
  chatHidden: boolean;
}

export default class SprintsScene extends React.Component<Props,State> {
  readonly state:State = { chatHidden: true };

  toggleChat():void {
    this.setState({ chatHidden: !this.state.chatHidden });
  }

  render() {
    return <PageLayout className='SprintsScene' {...this.props}>
      <div className='inner'>
        { this.renderUpper() }
        { this.renderLower() }
      </div>

      <Messenger
        title='Sprint #1 Chat'
        hidden={this.state.chatHidden}
        toggleHidden={this.toggleChat.bind(this)} />
    </PageLayout>
  }

  renderUpper():JSX.Element {
    return(
      <div className='upper'>
        <div className='left'>
          <div className='left-upper'>
            <h1>Sprint #1</h1>
            <img src={Arrow} className='arrow left' alt='Previous Sprint'/>
            <img src={Arrow} className='arrow right' alt='Next Sprint'/>
          </div>

          <div className='left-lower'>
            <p>12th Jan- 19th Jan 2020<span className='split'>•</span>1/3 Objectives completed</p>
          </div>
        </div>

        <div className='right'>
          <div className='chat-icon active noselect' onClick={this.toggleChat.bind(this)}>
            <img src={ChatIcon} alt='Chat' />
            <p>Chat Active</p>
          </div>
        </div>
      </div>
    );
  }

  renderLower():JSX.Element {
    return(
      <div className='lower'>
        <div className='section'>
          <div className='section-upper'>
            <h2>Objectives</h2>
            <img src={AddIcon} alt='Add Objective' />
          </div>

          <SprintRow status='completed' />
          <SprintRow status='pending' />
          <SprintRow status='pending' />
        </div>

        <div className='section'>
          <div className='section-upper'>
            <h2>Tasks</h2>
            <img src={AddIcon} alt='Add Task' />
          </div>

          <SprintRow status='completed' />
          <SprintRow status='completed' />
          <SprintRow status='pending' />
        </div>
      </div>
    );
  }
}