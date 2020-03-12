import { Component, OnInit, ViewEncapsulation } from "@angular/core";
import { ServiceService } from "../service.service";
import * as types from "../types";
import { Location } from "@angular/common";
import { ActivatedRoute } from "@angular/router";
import { Subject } from "rxjs";
import * as _ from "lodash";
import { NotificationsService } from "angular2-notifications";

const tabNamesToIndex = {
  "": 0,
  logs: 1,
  stats: 2,
  nodes: 3,
  traces: 4,
  events: 5
};

const tabIndexesToName = {
  0: "",
  1: "logs",
  2: "stats",
  3: "nodes",
  4: "traces",
  5: "events"
};

@Component({
  selector: "app-service",
  templateUrl: "./service.component.html",
  styleUrls: [
    "./service.component.css",
    "../../../node_modules/nvd3/build/nv.d3.css"
  ],
  encapsulation: ViewEncapsulation.None
})
export class ServiceComponent implements OnInit {
  services: types.Service[];
  logs: types.LogRecord[];
  stats: types.DebugSnapshot[];
  traceSpans: types.Span[];
  events: types.Event[];

  selectedVersion = "";
  serviceName: string;
  endpointQuery: string;
  intervalId: any;
  // refresh stats
  refresh = true;
  refreshLogs = true;

  selected = 0;
  tabValueChange = new Subject<number>();

  constructor(
    private ses: ServiceService,
    private activeRoute: ActivatedRoute,
    private location: Location,
    private notif: NotificationsService
  ) {}

  ngOnInit() {
    this.activeRoute.params.subscribe(p => {
      if (this.intervalId) {
        clearInterval(this.intervalId);
      }
      this.serviceName = <string>p["id"];
      this.ses.list().then(servs => {
        this.services = servs.filter(s => s.name == this.serviceName);
        this.selectedVersion =
          this.services.filter(s => s.version == "latest").length > 0
            ? "latest"
            : this.services[0].version;
      });
      this.ses
        .events(this.serviceName)
        .then(events => {
          this.events = events;
        })
        .catch(e => {
          this.notif.error(
            "Error listing events",
            JSON.parse(e.error.error).detail
          );
        });
      this.loadVersionData();
      const tab = <string>p["tab"];
      if (tab) {
        this.selected = tabNamesToIndex[tab];
      }
    });
  }

  loadVersionData() {
    this.ses
      .trace(this.serviceName)
      .then(spans => {
        this.traceSpans = spans;
      })
      .catch(e => {
        console.log(e);
        this.notif.error(
          "Error listing trace",
          JSON.parse(e.error.error).detail
        );
      });
    // stats subscriptions
    let statsFailure = false;
    this.intervalId = setInterval(() => {
      if (this.selected !== 2 || !this.refresh) {
        return;
      }
      this.ses
        .stats(this.serviceName)
        .then(stats => {
          this.stats = stats;
        })
        .catch(e => {
          if (statsFailure) {
            return;
          }
          statsFailure = true;
          this.notif.error(
            "Error reading stats",
            JSON.parse(e.error.error).detail
          );
        });
    }, 2000);
    this.tabValueChange.subscribe(index => {
      if (index !== 2 || !this.refresh) {
        return;
      }
      this.ses
        .stats(this.serviceName)
        .then(stats => {
          this.stats = stats;
        })
        .catch(e => {
          if (statsFailure) {
            return;
          }
          statsFailure = true;
          this.notif.error(
            "Error reading stats",
            JSON.parse(e.error.error).detail
          );
        });
    });
    // logs subscriptions
    let logFailure = false;
    this.intervalId = setInterval(() => {
      if (this.selected !== 1 || !this.refreshLogs) {
        return;
      }
      this.ses
        .logs(this.serviceName)
        .then(logs => {
          this.logs = logs;
        })
        .catch(e => {
          if (logFailure) {
            return;
          }
          logFailure = true;
          this.notif.error(
            "Error reading logs",
            JSON.parse(e.error.error).detail
          );
        });
    }, 2000);
    this.tabValueChange.subscribe(index => {
      if (index !== 1 || !this.refreshLogs) {
        return;
      }
      this.ses
        .logs(this.serviceName)
        .then(logs => {
          this.logs = logs;
        })
        .catch(e => {
          if (logFailure) {
            return;
          }
          logFailure = true;
          this.notif.error(
            "Error reading logs",
            JSON.parse(e.error.error).detail
          );
        });
    });
  }

  versionSelected(service: types.Service) {
    if (this.selectedVersion == service.version) {
      this.selectedVersion = "";
      return;
    }
    this.selectedVersion = service.version;
    this.loadVersionData();
  }

  tabChange($event: number) {
    this.selected = $event;
    this.location.replaceState(
      "/service/" + this.serviceName + "/" + tabIndexesToName[this.selected]
    );
    this.tabValueChange.next(this.selected);
  }

  ngOnDestroy() {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
  }

  code: string = "{}";
}
