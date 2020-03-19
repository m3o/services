import { Component, OnInit, Input, Pipe } from "@angular/core";
import * as types from "../types";
import * as testEvents from "./mock-events";

const eventTypes = {
  1: "ServiceCreated",
  2: "ServiceDeleted",
  3: "ServiceUpdated",
  4: "SourceCreated",
  5: "BuildStarted",
  6: "BuildFinished",
  7: "BuildFailed",
  8: "SourceUpdated",
  9: "SourceDeleted"
};

const eventTypesNice = {
  1: "service created",
  2: "service deleted",
  3: "service updated",
  4: "source created",
  5: "build started",
  6: "build finished",
  7: "build failed",
  8: "source updated",
  9: "source deleted"
};

@Component({
  selector: "app-events-list",
  templateUrl: "./events-list.component.html",
  styleUrls: ["./events-list.component.css"]
})
export class EventsListComponent implements OnInit {
  @Input() events: types.Event[];
  searched: types.Event[];
  eventsPart: types.Event[] = [];
  testEvents: types.Event[] = testEvents.default;
  query: string = "";

  public pageSize = 30;
  public currentPage = 0;
  public length = 0;

  constructor() {}

  ngOnInit() {
    //this.events = this.testEvents;
    this.refresh();
  }

  ngOnChanges(changes) {
    this.refresh();
  }

  refresh() {
    this.searched = this.events
      ? this.events.filter(e => {
          if (!e || !e.service || !e.service.name) {
            return false;
          }
          return true;
        })
      : [];
    this.length = this.searched.length;
    this.iterator();
  }

  eventTypeToString(e: types.Event): string {
    return eventTypesNice[e.type];
  }

  commitUrl(e: types.Event): string {
    if (!e.metadata) {
      return "";
    }
    const repo = e.metadata["repo"];
    const commitHash = e.metadata["commit"];
    // https://github.com/micro/services/commit/f291afc98f624c44e34e758efab89e77546b709d
    return "https://" + repo + "/commit/" + commitHash;
  }

  buildUrl(e: types.Event): string {
    if (!e.metadata) {
      return "";
    }
    const repo = e.metadata["repo"];
    const buildId = e.metadata["build"];
    // eg. https://github.com/micro/services/runs/466859781
    return "https://" + repo + "/actions/runs/" + buildId;
  }

  hasMeta(e: types.Event): boolean {
    return e.metadata && (e.metadata["commit"] || e.metadata["build"]);
  }

  visibleMeta(e: types.Event): Map<string, string> {
    if (!e.metadata) {
      return new Map();
    }
    return new Map([
      ["commit", e.metadata["commit"]],
      ["build", e.metadata["build"]]
    ]);
  }

  shortHash(s: string): string {
    if (!s) {
      return "";
    }
    return s.slice(0, 8);
  }

  public handlePage(e: any) {
    this.currentPage = e.pageIndex;
    this.pageSize = e.pageSize;
    this.iterator();
  }

  private iterator() {
    const end = (this.currentPage + 1) * this.pageSize;
    const start = this.currentPage * this.pageSize;
    const part = this.searched ? this.searched.slice(start, end) : [];
    this.eventsPart = part;
  }

  search() {
    this.searched = this.events.filter(e => {
      if (!this.query || this.query.length == 0) {
        return true;
      }
      if (!e.service || !e.service.name) {
        return false;
      }
      return e.service.name.includes(this.query);
    });
    this.currentPage = 0;
    this.length = this.searched.length;
    this.iterator();
  }
}
