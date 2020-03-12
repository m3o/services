import { Component, OnInit, Input } from "@angular/core";
import * as types from "../types";

@Component({
  selector: "app-stat-charts",
  templateUrl: "./stat-charts.component.html",
  styleUrls: ["./stat-charts.component.css"]
})
export class StatChartsComponent implements OnInit {
  @Input() stats: types.DebugSnapshot[] = [];
  constructor() {}

  ngOnInit() {}

  ngOnChanges(changes) {
    this.processStats();
  }

  processStats() {
    if (!this.stats) {
      return;
    }
    function onlyUnique(value, index, self) {
      return self.indexOf(value) === index;
    }
    const STAT_WINDOW = 8 * 60 * 1000; /* ms */
    this.stats = this.stats.filter(stat => {
      return Date.now() - stat.timestamp * 1000 < STAT_WINDOW;
    });
    const nodes = this.stats
      .map(stat => stat.service.node.id)
      .filter(onlyUnique);
    this.requestRates.data = nodes.map(node => {
      return {
        label: node,
        name: node,
        type: "line",
        pointRadius: 0,
        fill: false,
        lineTension: 0,
        borderWidth: 2,
        data: this.stats
          .filter(stat => stat.service.node.id == node)
          .map((stat, i) => {
            let value = stat.requests;
            if (i == 0 && this.stats.length > 0) {
              const first = this.stats[0].requests ? this.stats[0].requests : 0;
              value = this.stats[1].requests - first;
            } else {
              const prev = this.stats[i - 1].requests
                ? this.stats[i - 1].requests
                : 0;
              value = this.stats[i].requests - prev;
            }
            return {
              x: new Date(stat.timestamp * 1000),
              y: value ? value : 0
            };
          })
      };
    });

    this.memoryRates.data = nodes.map(node => {
      return {
        label: node,
        type: "line",
        pointRadius: 0,
        fill: false,
        lineTension: 0,
        borderWidth: 2,
        data: this.stats
          .filter(stat => stat.service.node.id == node)
          .map((stat, i) => {
            let value = stat.memory;
            return {
              x: new Date(stat.timestamp * 1000),
              y: value ? value / (1000 * 1000) : 0
            };
          })
      };
    });
    this.errorRates.data = nodes.map(node => {
      return {
        label: node,
        type: "line",
        pointRadius: 0,
        fill: false,
        lineTension: 0,
        borderWidth: 2,
        data: this.stats
          .filter(stat => stat.service.node.id == node)
          .map((stat, i) => {
            let value = stat.errors;
            if (i == 0 && this.stats.length > 0) {
              const first = this.stats[0].errors ? this.stats[0].errors : 0;
              value = this.stats[1].errors - first;
            } else {
              const prev = this.stats[i - 1].errors
                ? this.stats[i - 1].errors
                : 0;
              value = this.stats[i].errors - prev;
            }
            return {
              x: new Date(stat.timestamp * 1000),
              y: value ? value : 0
            };
          })
      };
    });
    let concMax = 0;
    this.concurrencyRates.data = nodes.map(node => {
      return {
        label: node,
        type: "line",
        pointRadius: 0,
        fill: false,
        lineTension: 0,
        borderWidth: 2,
        data: this.stats
          .filter(stat => stat.service.node.id == node)
          .map((stat, i) => {
            let value = stat.threads;
            if (value > concMax) {
              concMax = value;
            }
            return {
              x: new Date(stat.timestamp * 1000),
              y: value ? value : 0
            };
          })
      };
    });
    //this.concurrencyRates.options.scales.yAxes[0].ticks.max = concMax * 1.5;
    this.gcRates.data = nodes.map(node => {
      return {
        label: node,
        name: node,
        type: "line",
        pointRadius: 0,
        fill: false,
        lineTension: 0,
        borderWidth: 2,
        data: this.stats
          .filter(stat => stat.service.node.id == node)
          .map((stat, i) => {
            let value = stat.gc;
            if (i == 0 && this.stats.length > 0) {
              const first = this.stats[0].gc ? this.stats[0].gc : 0;
              value = this.stats[1].gc - first;
            } else {
              const prev = this.stats[i - 1].gc ? this.stats[i - 1].gc : 0;
              value = this.stats[i].gc - prev;
            }
            return {
              x: new Date(stat.timestamp * 1000),
              y: value ? value : 0
            };
          })
      };
    });
    this.uptime.data = nodes.map(node => {
      return {
        label: node,
        name: node,
        type: "line",
        pointRadius: 0,
        fill: false,
        lineTension: 0,
        borderWidth: 2,
        data: this.stats
          .filter(stat => stat.service.node.id == node)
          .map((stat, i) => {
            return {
              x: new Date(stat.timestamp * 1000),
              y: stat.uptime ? stat.uptime : 0
            };
          })
      };
    });
  }

  // config options taken from https://www.chartjs.org/samples/latest/scales/time/financial.html
  options(title: string, ylabel: string, distribution?: string) {
    if (!distribution) {
      distribution = "series";
    }
    return {
      options: {
        title: {
          display: true,
          text: title
        },
        //maintainAspectRatio: false,
        animation: {
          duration: 0
        },
        scales: {
          xAxes: [
            {
              type: "time",
              distribution: distribution,
              offset: true,
              ticks: {
                major: {
                  enabled: true,
                  fontStyle: "bold"
                },
                source: "data",
                autoSkip: true,
                autoSkipPadding: 75,
                maxRotation: 0,
                sampleSize: 100
              }
            }
          ],
          yAxes: [
            {
              gridLines: {
                drawBorder: false
              },
              scaleLabel: {
                display: true,
                labelString: ylabel
              }
            }
          ]
        },
        tooltips: {
          intersect: false,
          mode: "index",
          callbacks: {
            label: function(tooltipItem, myData) {
              var label = myData.datasets[tooltipItem.datasetIndex].label || "";
              if (label) {
                label += ": ";
              }
              label += parseFloat(tooltipItem.value).toFixed(2);
              return label;
            }
          }
        }
      },
      data: [],
      chartColors: [
        {
          // first color
          backgroundColor: "rgba(10,24,225,0.6)",
          borderColor: "rgba(10,24,225,0.6)",
          pointBackgroundColor: "rgba(10,24,225,0.6)",
          pointBorderColor: "#fff",
          pointHoverBackgroundColor: "#fff",
          pointHoverBorderColor: "rgba(10,24,225,0.6)"
        },
        {
          // second color
          backgroundColor: "rgba(10,24,225,0.6)",
          borderColor: "rgba(10,24,225,0.6)",
          pointBackgroundColor: "rgba(10,24,225,0.6)",
          pointBorderColor: "#fff",
          pointHoverBackgroundColor: "#fff",
          pointHoverBorderColor: "rgba(10,24,225,0.6)"
        }
      ],
      lineChartType: "line"
    };
  }
  memoryRates = this.options("Memory Usage", "memory usage (MB)");
  requestRates = this.options("Requests", "requests/second");
  errorRates = this.options("Errors", "errors/second");
  concurrencyRates = this.options("Goroutines", "goroutines");
  gcRates = this.options(
    "Garbage Collection",
    "garbage collection (nanoseconds/seconds)"
  );
  uptime = this.options("Uptime", "uptime (seconds)");
}
