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

  onRepositoryChange(e: any): void {
    if(e.target.value === "") e.target.value = undefined;
    this.setState({ project: {...this.state.project, repository: e.target.value}});
  }
  
  render(): JSX.Element {
    const { repository } = this.state.project;

    return(
      <PageLayout className='NewProject'>
        <div className='center'>
          <div className='header'>
            <h1>New Project</h1>
          </div>

          { this.renderProjectDetails() }
          { this.renderGithubToken() }
          { repository ? this.renderSecrets() : null }
        </div>
      </PageLayout>
    );
  }

  renderProjectDetails(): JSX.Element {
    const { name, description } = this.state.project;

    return(
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
    );
  }

  renderGithubToken(): JSX.Element {
    const { token, tokenStatus, repos } = this.state;
    const { repository } = this.state.project;

    return (
      <section>
        <h2>Connect to GitHub Repository</h2>
        <p>Enter a personal access token below. The token will need the <strong>repo</strong> and <strong>read:packages</strong> scopes. You can generate a new token at <a href='https://github.com/settings/tokens/new' target='blank'>this link</a>. Read more at the <a href=''>docs</a>.</p>

        <p className='status'>{tokenStatus}</p>

        <form>
          <div className='row'>
            <label>Token *</label>
            <input required disabled={repos.length > 0} type='text' value={token} onChange={this.onTokenChange.bind(this)} />
          </div>

          <div className='row'>
            <label>Repository *</label>
            <select value={repository} onChange={this.onRepositoryChange.bind(this)}>
              <option value=''>{repos.length > 0 ? 'Select a repository' : ''}</option>
              { repos.map(r => <option key={r} value={r}>{r}</option>) }
            </select>
          </div>
        </form>
      </section>
    );
  }

  renderSecrets(): JSX.Element {
    const addSecretsLink = `https://github.com/${this.state.project.repository}/settings/secrets`;

    return(
      <section>
        <h2>Setup Github Action</h2>
        <p>M3O provides a GitHub action <a href='https://github.com/micro/actions' target='blank'>(micro/actions)</a> which builds packages within your repository, giving you full ownership over your source and builds. The GitHub action requires the following secrets to authenticate with M3O. You can add the secrets at <a href={addSecretsLink} target='blank'>this link</a>.</p>
      </section>
    );
  }
}