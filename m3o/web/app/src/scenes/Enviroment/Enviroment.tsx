import React from 'react';
import PageLayout from '../../components/PageLayout';
import EditIcon from './assets/edit.png';
import './Enviroment.scss';

interface Props {
  match: any;
}

export default class Enviroment extends React.Component<Props> {
  render(): JSX.Element {
    const { project, enviroment } = this.props.match.params;

    return <PageLayout className='Enviroment'>
      <div className='center'>
        <div className='header'>
          <h1><span>{project}</span> / <span>{enviroment}</span></h1>
          <img src={EditIcon} alt='Edit Name' />

          <button className='btn'><p>Launch Dashboard</p></button>
        </div>

        <section>
          <h2>Enviroment Details</h2>
          <p>These details are only visible to you and collaborators. All M3O projects are private.</p>

          <form>
            <div className='row'>
              <label>Name *</label>
              <input required type='text' value='Production' placeholder='My Awesome Project' name='name' />
            </div>
            
            <div className='row'>
              <label>Description</label>
              <input type='text' value='The Kytra production environment' placeholder='Description' name='description' />
            </div>
          </form>
        </section>

        <section>
          <h2>DNS</h2>
          <p>Your default domain is <a href='https://production.kytra.m3o.app' target='blank'>https://production.kytra.m3o.app</a>. Your web domain is served at <a href='/todo'>/ (root)</a> and your API is available at <a href='/todo'>/api</a>. To configure a custom domain, enter the domains below and then setup CNAME records for each domain pointing at <strong>m3o.app</strong>. For more information about custom domains, see <a href='/todo'>the docs</a>.</p>
          <form>
            <div className='row'>
              <label>Web Domain</label>
              <input type='text' value='' placeholder='myapp.com' name='web_domain' />
            </div>
            
            <div className='row'>
              <label>API Domain</label>
              <input type='text' value='' placeholder='api.myapp.com' name='api_domain' />
            </div>
          </form>
        </section>

        <section>
          <h2>CLI</h2>
          <p>Configure your CLI to use the {project}/{enviroment} enviroment. Firstly, all calls made to your enviroment are authenticated so if you aren't already, login using the following command and a token you can get at <a href='/todo'>this link</a>.</p>
          <p className='code'>
            micro login [token]
          </p>

          <p>Once you're logged in, add your enviroment and configure micro to use it with the following commands.</p>
          <p className='code'>
            micro env add {project}/{enviroment} {project}.{enviroment}.proxy.m3o.app
            <br />
            micro env set {project}/{enviroment}
          </p>
        </section>

        <section>
          <h2>Settings</h2>
          <p><strong>Warning:</strong> Deleting your enviroment cannot be undone and all data will be lost.</p>
          <button className='btn danger'>Delete {project}/{enviroment}</button>
        </section>
     </div>
    </PageLayout>
  }
}