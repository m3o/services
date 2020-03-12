import { Component, OnInit } from '@angular/core';
import * as types from "../types";
import { ProjectService } from '../project.service';

@Component({
  selector: 'app-new-project',
  templateUrl: './new-project.component.html',
  styleUrls: ['./new-project.component.css']
})
export class NewProjectComponent implements OnInit {
  organisations: types.Organisation[] = [];
  repositories: types.Repository[] = [];
  step = 0;
  alias = "my-first-app";
  projectExists = false;
  loadingProjects = false;
  loaded = true;
  selectedOrg: types.Organisation;
  selectedRepo: types.Repository;

  constructor(private ps: ProjectService) { }

  ngOnInit() {
    this.ps.listOrganisations().then(orgs => {
      this.organisations = orgs
    })
  }

  keyPress($event) {

  }

  orgSelected(v: string) {
    this.ps.listRepositories(v).then(repos => {
      this.repositories = repos
    })
  }

  repoSelected(v: string) {

  }
}
