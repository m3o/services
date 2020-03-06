import { NgModule } from "@angular/core";
import { Routes, RouterModule } from "@angular/router";
import { SubscribeFormComponent } from "./subscribe-form/subscribe-form.component";
import { SubscriberListComponent } from "./subscriber-list/subscriber-list.component";

const routes: Routes = [
  {
    path: "",
    component: SubscribeFormComponent
  },
  {
    path: "list",
    component: SubscriberListComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
