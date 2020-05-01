import React from 'react';
import PageLayout from '../../components/PageLayout';
import * as API from '../../api';
import './NewProject.scss';

interface Props {}

interface State {
  project: API.Project;
  token: string;
  tokenStatus: string;
  repos: string[];
}

export default class NewProject extends React.Component<Props, State> {
  readonly state: State = {
    token: '',
    repos: [],
    tokenStatus: 'Waiting for token...',
    project: { name: '', description: '' },
  };

  onInputChange(e: any): void {
    this.setState({ project: { ...this.state.project, [e.target.name]: e.target.value } });
  }

  onTokenChange(e: any): void {
    if(this.state.repos.length > 0) return;
    this.setState({ token: e.target.value, tokenStatus: "Validating token, please wait" });

    API.Call("ProjectService/VerifyGithubToken", { token: e.target.value })
      .then((res) => this.setState({ tokenStatus: "Token Valid. Please select a repository from the list below.", repos: res.data.repos }))
      .catch((err) => this.setState({ tokenStatus: err.response.data.detail }));
  }

  render(): JSX.Element {
    const { token, tokenStatus, repos } = this.state;
    const { name, description, repository } = this.state.project;

    return(
      <PageLayout className='NewProject'>
        <div className='center'>
          <div className='header'>
            <h1>New Project</h1>
          </div>

          <section>
            <h2>Project Details</h2>
            <p>Let's start by entering some basic project information</p>

            <form>
              <div className='row'>
                <label>Name *</label>
                <input required type='text' value={name} placeholder='My Awesome Project' name='name' onChange={this.onInputChange.bind(this)} />
              </div>
              
              <div className='row'>
                <label>Description</label>
                <input type='text' value={description} placeholder='' name='description'  onChange={this.onInputChange.bind(this)} />
              </div>
            </form>
          </section>

          <section>
            <h2>Connect to GitHub Repository</h2>
            <p>Enter a personal access token below. The token will need the <strong>repo</strong> and <strong>read:packages</strong> scopes. You can generate a new token at <a href='https://github.com/settings/tokens/new' target='blank'>this link</a>. Read more at the <a href=''>docs</a>.</p>

            <p className='token-status'>{tokenStatus}</p>

            <form>
              <div className='row'>
                <label>Token *</label>
                <input required disabled={repos.length > 0} type='text' value={token} onChange={this.onTokenChange.bind(this)} />
              </div>

              <div className='row'>
                <label>Repository *</label>
                <select value={repository}>
                  <option value=''>{repos.length > 0 ? 'Select a repository' : ''}</option>
                  { repos.map(r => <option key={r} value={r}>{r}</option>) }
                </select>
              </div>
            </form>
          </section>
        </div>
      </PageLayout>
    );
  }
}