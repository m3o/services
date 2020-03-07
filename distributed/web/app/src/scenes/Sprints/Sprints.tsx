import React from 'react';
import './Sprints.scss';
import PageLayout from '../../components/PageLayout';
import Arrow from '../../assets/images/arrow.png';
import ChatIcon from '../../assets/images/chat-icon-white.png';
import AddIcon from '../../assets/images/add-icon.png';
import SprintRow from './components/Row';

export default class SprintsScene extends React.Component {
  render() {
    return <PageLayout className='SprintsScene' {...this.props}>
      { this.renderUpper() }
      { this.renderLower() }
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
          <div className='chat-icon active noselect'>
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