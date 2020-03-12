import { Component, OnInit, Input } from "@angular/core";
import * as types from "../types";
import * as _ from "lodash";

@Component({
  selector: "app-trace-list",
  templateUrl: "./trace-list.component.html",
  styleUrls: ["./trace-list.component.css"]
})
export class TraceListComponent implements OnInit {
  @Input() serviceName: string;
  @Input() spans: any[] = [];

  traceDatas: any[] = [];
  traceDatasPart: any[] = [];
  public pageSize = 20;
  public currentPage = 0;
  public length = 0;

  public handlePage(e: any) {
    this.currentPage = e.pageIndex;
    this.pageSize = e.pageSize;
    this.iterator();
  }

  private iterator() {
    const end = (this.currentPage + 1) * this.pageSize;
    const start = this.currentPage * this.pageSize;
    const part = this.traceDatas.slice(start, end);
    this.traceDatasPart = part;
  }

  constructor() {}

  ngOnChanges(changes) {
    this.processTraces();
  }

  ngOnInit() {}

  prettyId(id: string) {
    return id.substring(0, 8);
  }

  show(td) {
    td.show = !td.show;
    return false;
  }

  prettyTime(ms: number): string {
    if (ms < 1000) {
      return Math.floor(ms) + "ms";
    }
    return (ms / 1000).toFixed(3) + "s";
  }

  traceDuration(spans: (String | Date)[][]): string {
    const durations = spans.slice(1).map(span => {
      return (span[3] as Date).getTime() - (span[2] as Date).getTime();
    });

    return this.prettyTime(durations.reduce((a, b) => a + b, 0));
  }

  getEndpointName(spans: (String | Date)[][]): string {
    return (spans.slice(1).filter(span => {
      return (span[1] as string).includes(this.serviceName);
    })[0][1] as string)
      .split(":")[1]
      .split(" ")[1];
  }

  processTraces() {
    const spans = this.spans;
    if (!spans) {
      return;
    }
    const groupedSpans = _.values(_.groupBy(_.uniqBy(spans, "id"), "trace"));
    let traceDatas: any[] = [];
    groupedSpans.forEach(spanGroup => {
      const spansToDisplay = _.orderBy(
        spanGroup.map((d, index) => {
          let start = d.started / 1000000;
          let end = (d.started + d.duration) / 1000000;
          let name = "Handle: " + d.name + " " + this.prettyTime(end - start);
          if (d.type == 1) {
            name = "Call: " + d.name + " " + this.prettyTime(end - start);
          }
          return ["", name, new Date(start), new Date(end)];
        }),
        sp => {
          const row = sp as Date[];
          return row[2];
        },
        ["asc"]
      );
      spansToDisplay.forEach((v, i) => {
        v[0] = "" + i;
      });

      const minMax = (): [Date, Date] => {
        const firstStart = (spansToDisplay[0][2] as Date).getTime();
        const lastEnd = (spansToDisplay[
          spansToDisplay.length - 1
        ][3] as Date).getTime();
        let leftPad = 1;
        let rightPad = 1;
        if (lastEnd - firstStart < 1000) {
          leftPad = (1000 - (lastEnd - firstStart)) / 2;
          rightPad = (1000 - (lastEnd - firstStart)) / 2;
        }
        const minDate = new Date(firstStart - leftPad);
        const maxDate = new Date(lastEnd + rightPad);
        return [minDate, maxDate];
      };

      const h = (spansToDisplay.length + 1) * 40 + 40;
      const [min, max] = minMax();
      let traceData = {
        // Display related things
        traceId: spanGroup[0].trace,
        divHeight: h + 20,
        // Chart related options
        chartType: "Timeline",

        dataTable: ([["Span", "Name", "From", "To"]] as any[][]).concat(
          spansToDisplay
        ),
        options: {
          height: h,
          timeline: {
            tooltipDateFormat: "HH:mm:ss.SSS"
          },
          hAxis: {
            format: "yyyy-MM-dd HH:mm:ss.SSS",
            minValue: min,
            maxValue: max
          }
        }
      };
      traceDatas.push(traceData);
    });
    this.traceDatas = _.orderBy(traceDatas, td => td.dataTable.length, [
      "desc"
    ]);
    this.length = this.traceDatas.length;
    this.iterator();
  }
}
