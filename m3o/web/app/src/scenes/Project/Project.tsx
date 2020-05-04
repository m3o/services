import React from 'react';
import PageLayout from '../../components/PageLayout';
import './Project.scss';

interface Props {
  match: any;
}

export default class Project extends React.Component<Props> {
  render(): JSX.Element {
    // const { project } = this.props.match.params;

    return <PageLayout className='Project'>
      <div className='center'>
        <div className='header'>
          <h1>{this.props.match.params.project}</h1>
        </div>

        <section>
          <h2>Project Details</h2>
          <p>These details are only visible to you and collaborators. All M3O projects are private.</p>

          <form>
            <div className='row'>
              <label>Name *</label>
              <input required type='text' value={this.props.match.params.project} placeholder='My Awesome Project' name='name' />
            </div>
            
            <div className='row'>
              <label>Description</label>
              <input type='text' value='Description' placeholder='Description' name='description' />
            </div>
          </form>
        </section>

        <section>
          <h2>GitHub</h2>
          <p>M3O connects to GitHub and builds your services in your repo, keeping your source and builds firmly in your control. The <a href='https://github.com/micro/actions' target='blank'>micro/actions</a> GitHub action automatically builds your services when any changes are detected and triggers a release. Find our more at our <a href='/todo'>docs</a>.</p>

          <form>
            <div className='row'>
              <label>Repository</label>
              <input disabled type='text' value='kytra/backend' name='repository' />
            </div>
            <div className='row'>
              <label>Client ID</label>
              <input disabled type='text' value='************' name='repository' />
            </div>
            <div className='row'>
              <label>Client Secret</label>
              <input disabled type='text' value='************************' name='repository' />
            </div>
          </form>

          <p>Configure the GitHub action using your Client ID and Secret. If you loose your ID/Secret, click the regenerate button below to generate a new set of credentials.</p>
          <button className='btn warning'>Regenerate Credentials</button>
        </section>

        <section>
          <h2>Collaborators</h2>
          <p>Collaborators have full access to all enviroments, but only the owner (you) can invite additional collaborators or delete enviroments.</p>

          <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Actions</th>
            </tr>
          </thead>

          <tbody>
            <tr key={'asim'}>
              <td>Asim Aslam</td>
              <td>asim@micro.mu</td>
              <td>
                <button className='danger'>Remove</button>
              </td>
            </tr>
            <tr key={'jake'}>
              <td>Jake Sanders</td>
              <td>jake@micro.mu</td>
              <td>
                <button className='danger'>Remove</button>
              </td>
            </tr>
            <tr key={'ben'}>
              <td>Ben Toogood</td>
              <td>ben@micro.mu</td>
              <td>
                <button className='danger'>Remove</button>
              </td>
            </tr>
          </tbody>
        </table>

        <p>Invite users to your project. Collaborators will recieve an email invite which is valid for 24 hours. </p>
        <form>
          <div className='row'>
            <label>Name</label>
            <input required type='text' placeholder='John Doe' name='name' />
          </div>

          <div className='row'>
            <label>Email</label>
            <input required type='email' placeholder='john@doe.com' name='email' />
          </div>

          <button className='btn'>Send Invite</button>
        </form>
      </section>
     </div>
    </PageLayout>
  }
}