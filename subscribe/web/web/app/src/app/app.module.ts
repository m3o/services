import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";

import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { SubscribeFormComponent } from "./subscribe-form/subscribe-form.component";
import { ClientModule } from "@microhq/ng-client";
import { RouterModule } from "@angular/router";
import { FormsModule } from "@angular/forms";
import { SubscriberListComponent } from './subscriber-list/subscriber-list.component';

@NgModule({
  declarations: [AppComponent, SubscribeFormComponent, SubscriberListComponent],
  imports: [
    BrowserModule,
    AppRoutingModule,
    RouterModule,
    ClientModule,
    FormsModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {}
