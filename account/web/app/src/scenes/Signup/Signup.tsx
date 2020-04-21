import React from 'react';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import { User } from '../../api';
import EditProfile from '../../components/EditProfile';
import './Signup.scss';

interface Props {
  user: User;
  history: any;
}

class Signup extends React.Component<Props> {
  render(): JSX.Element {
    return(
      <div className='Signup'>
        <div className='inner'>
          <h1>Welcome to Micro</h1>

          <div className='profile'>
            <p>Let's get started by completing your Micro profile</p>
            <EditProfile onSave={() => window.location.href = 'https://m3o.micro.mu/'} />
        </div>
        </div>
      </div>
    );
  }
}

function mapStateToProps(state: any):any {
  return({
    user: state.user.user,
  });
}

export default withRouter(connect(mapStateToProps)(Signup));