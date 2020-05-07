// Frameworks
import React from 'react';
import Gist from 'react-gist';
import { connect } from 'react-redux';

// Components
import PageLayout from '../../components/PageLayout';
import ValidatedInput from '../../components/ValidatedInput';
import PaymentMethod from './components/PaymentMethod';

// Utils
import * as API from '../../api';
import { createProject } from '../../store/Project';

// Styling
import OpenSourceIcon from './assets/opensource.png';
import DeveloperIcon from './assets/developer.png';
import TeamIcon from './assets/team.png';
import './NewProject.scss';

interface Props {
  history: any;
  createProject: (project: API.Project) => void;
}

interface Repository {
  name: string;
  public?: boolean;
}

interface State {
  project: API.Project;
  nameValid: boolean;
  token: string;
  tokenStatus: string;
  repos: Repository[];
  repository?: Repository;
  clientID?: string;
  clientSecret?: string;
  paymentPlan?: string;
  paymentMethodStatus?: string;
  paymentMethodDisabled: boolean;
}

// regex to check for specical chars
var regex = /[^\w]|_/g


class NewProject extends React.Component<Props, State> {
  readonly ref: React.RefObject<HTMLDivElement> = React.createRef();

  readonly state: State = {
    token: '',
    repos: [],
    tokenStatus: 'Waiting for token...',
    project: { name: '', description: '' },
    paymentMethodDisabled: false,
    nameValid: false,
  };

  onRepositoryChange(e: any): void {
    const repoName: string = e.target.value;
    const repo = this.state.repos.find(r => r.name === repoName);

    if(!repo) {
      this.setState({ 
        repository: undefined,
        project: { ...this.state.project, repository: '' },
      });
      return;
    };

    this.setState({
      project: {...this.state.project, repository: repo.name },
      repository: repo,
    });

    setTimeout(this.scrollToBottom.bind(this), 100);
  }
  
  render(): JSX.Element {
    const { repository, project, paymentPlan, nameValid } = this.state;

    return(
      <PageLayout className='NewProject' childRef={this.ref}>
        <div className='center'>
          <div className='header'>
            <h1>New Project</h1>
          </div>

          { this.renderProjectDetails() }
          { nameValid ? this.renderGithubToken() : null }
          { repository ? this.renderPlans() : null }
          { paymentPlan ? this.renderPaymentMethod() : null }
          { project.id ? this.renderSecrets() : null }
        </div>
      </PageLayout>
    );
  }

  renderProjectDetails(): JSX.Element {
    const { id, name, description } = this.state.project;

    const validateName = async (name: string): Promise<void> => {      
      return new Promise(async (resolve: Function, reject: Function) => {
        if(name.length < 3) {
          reject("Name must be at least 3 characters long");
          return
        }

        if(regex.test(name)) {
          reject("Name cannot contain any special characters");
          return;
        }

        API.Call('Projects/ValidateProjectName', { name })
          .then(() => resolve())
          .catch(err => reject(err.response ? err.response.data.detail : err.message));
      });
    }

    const setNameValid = () => this.setState({ nameValid: true });
    const setNameInvalid = () => this.setState({ nameValid: false });
    
    const onChange = (key: string, value: string) => {
      if(key === 'name') value = value.toLowerCase();
      this.setState({ project: { ...this.state.project, [key]: value } });
    };

    return(
      <section className='complete'>
        <h2>Project Details</h2>
        <p>Let's start by entering some basic project information</p>

        <form>
          <div className='row'>
            <label>Name *</label>

            <ValidatedInput
              autoFocus
              name='name'
              value={name}
              disabled={!!id}
              onChange={onChange}
              onValid={setNameValid}
              validate={validateName} 
              onInvalid={setNameInvalid}
              placeholder='helloworld' />
          </div>
          
          <div className='row'>
            <label>Description</label>

            <ValidatedInput
              disabled={!!id}
              name='description'
              value={description}
              onChange={onChange}
              placeholder='My Awesome Project' />
          </div>
        </form>
      </section>
    );
  }

  renderGithubToken(): JSX.Element {
    const { token, repos } = this.state;
    const { repository } = this.state.project;

    const validateToken = (token: string): Promise<void> => {
      return new Promise(async (resolve: Function, reject: Function) => {
        API.Call("Projects/ValidateGithubToken", { token })
          .then((res) => {
            this.setState({ repos: res.data.repos });
            resolve();
          })
          .catch(err => reject(err.response ? err.response.data.detail : err.message));
      });
    }

    return (
      <section>
        <h2>Connect to GitHub Repository</h2>
        <p>Enter a personal access token below. The token will need the <strong>repo</strong> and <strong>read:packages</strong> scopes. You can generate a new token at <a href='https://github.com/settings/tokens/new' target='blank'>this link</a>. Read more at the <a href='/todo'>docs</a>.</p>

        <form>
          <div className='row'>
            <label>Token *</label>

            <ValidatedInput
              name='name'
              value={token}
              validate={validateToken} 
              disabled={repos.length > 0}
              onChange={((_, token: string) => this.setState({ token }))} />
          </div>

          <div className='row'>
            <label>Repository *</label>
            <select value={repository} onChange={this.onRepositoryChange.bind(this)}>
              <option value=''>{repos.length > 0 ? 'Select a repository' : ''}</option>
              { repos.map(r => <option key={r.name} value={r.name}>{r.name}</option>) }
            </select>
          </div>
        </form>
      </section>
    );
  }

  renderPlans(): JSX.Element {
    const setPlan = (paymentPlan: string) => {
      this.setState({ paymentPlan });
      setTimeout(this.scrollToBottom.bind(this), 100);
    };

    return(
      <section>
        <h2>Payment Tiers</h2>
        <p>Select one of the payment tiers below. The community tier is only available to public repositories with an Apache License. See <a href='/todo'>the docs</a> for more information on pricing.</p>

        <div className='payment-plans'>
          <div className='plan'>
            <div className='img-wrapper'>
              <img src={OpenSourceIcon} alt='Community'/>
            </div>

            <h5>Community</h5>
            <h6>Built for open-source</h6>
            
            <p className='attr'><span>Single</span> Enviroment</p>
            <p className='attr'><span>Unlimited</span> Collaborators</p>
            
            <p className='price'><span>$0</span>/month</p>

            <button onClick={() => setPlan('community')} className='btn info'><p>Choose Community</p></button>
          </div>

          <div className='plan'>
            <div className='img-wrapper'>
              <img src={DeveloperIcon} alt='Community'/>
            </div>

            <h5>Developer</h5>
            <h6>Perfect for Indie Hackers</h6>
            
            <p className='attr'><span>Single</span> Enviroment</p>
            <p className='attr'><span>No</span> Collaborators</p>
            
            <p className='price'><span>$35</span>/month</p>

            <button onClick={() => setPlan('developer')} className='btn info'><p>Choose Developer</p></button>
          </div>

          <div className='plan'>
            <div className='img-wrapper'>
              <img src={TeamIcon} alt='Community'/>
            </div>

            <h5>Team</h5>
            <h6>Ideal for Startups</h6>
            
            <p className='attr'><span>5</span> Enviroments</p>
            <p className='attr'><span>Unlimited</span> Collaborators</p>
            
            <p className='price'><span>$45</span>/user per month</p>

            <button onClick={() => setPlan('team')} className='btn info'><p>Choose Team</p></button>
          </div>
        </div>
      </section>
    );
  }

  renderPaymentMethod(): JSX.Element {
    const onComplete = (paymentMethodID: string) => {
      const params = {
        github_token: this.state.token,
        // payment_method_id: paymentMethodID,
        project: {
          repository: this.state.project.repository,
          name: this.state.project.name,
          description: this.state.project.description,
        },
      };

      API.Call("Projects/CreateProject", params)
        .then(res => {
          this.setState({ 
            project: res.data.project,
            clientID: res.data.client_id,
            clientSecret: res.data.client_secret,
            paymentMethodStatus: 'Subscription Setup.'
          })

          setTimeout(this.scrollToBottom.bind(this), 100);    
        })
        .catch(err => onError(err.response.data.detail));
    }

    const onError = (err: string) => {
      this.setState({
        paymentMethodDisabled: false,
        paymentMethodStatus: `Error: ${err}`,
      });
    }

    const onSubmit = () => {
      this.setState({ 
        paymentMethodDisabled: true,
        paymentMethodStatus: 'Creating Subscription...',
      });
    }

    return(
      <section>
        <h2>Setup Billing</h2>
        <p>Add a payment method for your project. Payments are processed by <a href='/todo'>Stripe</a> and taken on the first of each month. For more information, see <a href='/todo'>the docs</a>.</p>
        <PaymentMethod status={this.state.paymentMethodStatus} onSubmit={onSubmit} onComplete={onComplete} onError={onError} />
      </section>
    )
  }

  renderSecrets(): JSX.Element {
    const { project, clientID, clientSecret } = this.state;
    const addSecretsLink = `https://github.com/${project.repository}/settings/secrets`;

    return(
      <section>
        <h2>Setup Github Action</h2>
        <p>M3O provides a GitHub action <a href='https://github.com/micro/actions' target='blank'>(micro/actions)</a> which builds packages within your repository, giving you full ownership over your source and builds. The GitHub action requires the following secrets to authenticate with M3O. You can add the secrets at <a href={addSecretsLink} target='blank'>this link</a>.</p>

        <form onSubmit={null}>
          <div className='row'>
            <label>M3O_CLIENT_ID</label>
            <input type='text' disabled value={clientID} />
          </div>
          <div className='row'>
            <label>M3O_CLIENT_SECRET</label>
            <input type='text' disabled value={clientSecret} />
          </div>
        </form>

        <p>Commit the following file to your repo as <strong>.github/workflows/m3o.yaml</strong></p>
        <Gist id="cd6ed0ae96e83c49569f877be7a22b32" />

        <button className='btn' onClick={this.done.bind(this)}>Done</button>
      </section>
    );
  }

  done(): void {
    this.props.createProject(this.state.project);
    this.props.history.push(`/projects/${this.state.project.name}`);
  }

  scrollToBottom(): void {
    this.ref.current.scroll({
      top: this.ref.current.scrollHeight, 
      left: 0, 
      behavior: 'smooth'
    });
  }
}

function mapDispatchToProps(dispatch: Function): any {
  return({
    createProject: (project: API.Project) => dispatch(createProject(project)),
  });
}

export default connect(null, mapDispatchToProps)(NewProject);