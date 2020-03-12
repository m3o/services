import { NgModule } from "@angular/core";
import { Routes, RouterModule } from "@angular/router";
import { ServicesComponent } from "./services/services.component";
import { ServiceComponent } from "./service/service.component";
import { NewProjectComponent } from "./new-project/new-project.component";
import { AuthGuard } from "./auth.guard";
import { WelcomeComponent } from "./welcome/welcome.component";
import { NotInvitedComponent } from "./not-invited/not-invited.component";
import { SettingsComponent } from "./settings/settings.component";
import { EventsComponent } from "./events/events.component";

const routes: Routes = [
  {
    path: "",
    component: WelcomeComponent,
    pathMatch: "full"
  },
  {
    path: "not-invited",
    component: NotInvitedComponent
  },
  {
    path: "project/new",
    component: NewProjectComponent,
    canActivate: [AuthGuard]
  },
  {
    path: "project/new/:id",
    component: NewProjectComponent,
    canActivate: [AuthGuard]
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
