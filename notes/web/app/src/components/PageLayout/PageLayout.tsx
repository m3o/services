import React from 'react';
import './PageLayout.scss';

interface Props {
  className: string;
  match?: any;
}

export default class PageLayout extends React.Component<Props> {
  render(): JSX.Element {
    const { className } = this.props;

    return(
      <div className='PageLayout'>
        <div className={`content ${className}`}>
          { this.props.children }
        </div>
      </div>
    );
  }
}