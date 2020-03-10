import React from 'react';
import './Messenger.scss';
import Cross from '../../assets/images/close-icon-white.png';

interface Props {
  title: string;
  hidden: boolean;
  toggleHidden: () => void;
}

export default class Messenger extends React.Component<Props> {
  render():JSX.Element {
    const { title, hidden, toggleHidden } = this.props;

    return(
      <div className={`Messenger ${hidden ? 'hidden' : 'visible'}`}>
        <div className='messenger-wrapper'>
          <div className='header'>
            <div className='left'>
              <h3>{title}</h3>
              <p>Asim and 4 others online</p>
            </div>

            <img src={Cross} alt='Close Chat' onClick={toggleHidden} />
          </div>
        </div>
      </div>
    )
  }
}