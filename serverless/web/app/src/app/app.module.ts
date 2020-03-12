import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";

import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { HeaderComponent } from "./header/header.component";
import { HomeComponent } from "./home/home.component";

import {
  MatTabsModule,
  MatSidenavModule,
  MatToolbar,
  MatList,
  MatMenu,
  MatProgressSpinnerModule
} from "@angular/material";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";

import { MatToolbarModule } from "@angular/material";
import {
  MatIconModule,
  MatButtonModule,
  MatMenuModule,
  MatCardModule,
  MatChipsModule,
  MatFormFieldModule,
  MatInputModule,
  MatExpansionModule,
  MatProgressBarModule,
  MatCheckboxModule
} from "@angular/material";
import { MatPaginatorModule } from "@angular/material/paginator";
import { MatListModule } from "@angular/material";
import { FlexLayoutModule } from "@angular/flex-layout";
import { LoginComponent } from "./login/login.component";
import { ServicesComponent } from "./services/services.component";

import { CookieService } from "ngx-cookie-service";
import { UserService } from "./user.service";
import { HttpClientModule } from "@angular/common/http";
import { SimpleNotificationsModule } from "angular2-notifications";
import { ServiceComponent } from "./service/service.component";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { SearchPipe } from "./search.pipe";
import { NewServiceComponent } from "./new-service/new-service.component";

import { ChartsModule } from "ng2-charts";
import { WelcomeComponent } from "./welcome/welcome.component";
import { LogUserInComponent } from "./log-user-in/log-user-in.component";

import { ClipboardModule } from "ngx-clipboard";
import { HighlightModule, HIGHLIGHT_OPTIONS } from "ngx-highlightjs";
import { NotInvitedComponent } from "./not-invited/not-invited.component";

import { Ng2GoogleChartsModule } from "ng2-google-charts";
import { MonacoEditorModule } from "ngx-monaco-editor";
import { SettingsComponent } from "./settings/settings.component";
import { EventsComponent } from "./events/events.component";
import { DateAgoPipe } from "./dateago.pipe";
import { EventsListComponent } from "./events-list/events-list.component";
import { StatChartsComponent } from "./stat-charts/stat-charts.component";
import { TraceListComponent } from "./trace-list/trace-list.component";
import { EndpointListComponent } from "./endpoint-list/endpoint-list.component";
import { LogsComponent } from "./logs/logs.component";
import { NodesComponent } from "./nodes/nodes.component";

/**
 * Import specific languages to avoid importing everything
 * The following will lazy load highlight.js core script (~9.6KB) + the selected languages bundle (each lang. ~1kb)
 */
export function getHighlightLanguages() {
  return {
    typescript: () => import("highlight.js/lib/languages/typescript"),
    css: () => import("highlight.js/lib/languages/css"),
    xml: () => import("highlight.js/lib/languages/xml"),
    bash: () => import("highlight.js/lib/languages/bash")
  };
}

@NgModule({
  declarations: [
    AppComponent,
    HeaderComponent,
    HomeComponent,
    LoginComponent,
    ServicesComponent,
    ServiceComponent,
    SearchPipe,
    NewServiceComponent,
    WelcomeComponent,
    LogUserInComponent,
    NotInvitedComponent,
    SettingsComponent,
    EventsComponent,
    DateAgoPipe,
    EventsListComponent,
    StatChartsComponent,
    TraceListComponent,
    EndpointListComponent,
    LogsComponent,
    NodesComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    MatSidenavModule,
    MatTabsModule,
    MatToolbarModule,
    MatIconModule,
    MatButtonModule,
    MatListModule,
    FlexLayoutModule,
    MatMenuModule,
    HttpClientModule,
    SimpleNotificationsModule.forRoot({
      //position: ["top", "right"],
    }),
    MatCardModule,
    MatChipsModule,
    MatFormFieldModule,
    MatInputModule,
    FormsModule,
    ReactiveFormsModule,
    MatProgressSpinnerModule,
    MatExpansionModule,
    MatProgressBarModule,
    ChartsModule,
    ClipboardModule,
    HighlightModule,
    Ng2GoogleChartsModule,
    MatPaginatorModule,
    MatCheckboxModule,
    MonacoEditorModule.forRoot()
  ],
  providers: [
    CookieService,
    UserService,
    {
      provide: HIGHLIGHT_OPTIONS,
      useValue: {
        languages: getHighlightLanguages()
      }
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule {}
