import React from 'react';
import './Row.scss';

interface Props {
  status: string;
}

export default class SprintRow extends React.Component<Props> {
  render(): JSX.Element {
    const { status } = this.props;

    return(
      <div className={`SprintRow ${status}`}>
        <p className='title'>This is an example title</p>
        <p className={`status noselect ${status}`}>{status}</p>
      </div>
    )
  }
}