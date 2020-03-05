import { NgModule } from "@angular/core";
import { Routes, RouterModule } from "@angular/router";
import { SubscribeFormComponent } from "./subscribe-form/subscribe-form.component";

const routes: Routes = [
  {
    path: "",
    component: SubscribeFormComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
