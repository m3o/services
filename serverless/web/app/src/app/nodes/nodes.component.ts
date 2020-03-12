import { Component, OnInit, Input } from "@angular/core";
import * as types from "../types";
import * as _ from "lodash";

@Component({
  selector: "app-nodes",
  templateUrl: "./nodes.component.html",
  styleUrls: ["./nodes.component.css"]
})
export class NodesComponent implements OnInit {
  @Input() services: types.Service[] = [];
  nodes: types.Node[];
  constructor() {}

  ngOnInit() {
    this.nodes = _.uniqBy(
      _.flatten(this.services.map(s => s.nodes)),
      n => n.id
    );
    //this.nodes.push(this.nodes[0])
  }

  metadata(node: types.Node) {
    let serialised = "No metadata.";
    if (!node.metadata) {
      return serialised;
    }
    const v = JSON.parse(JSON.stringify(node.metadata));
    serialised = "";
    let maxKeyLength = 0;
    for (var key in v) {
      if (maxKeyLength < key.length) {
        maxKeyLength = key.length;
      }
    }
    console.log(maxKeyLength)
    for (var key in v) {
      console.log(maxKeyLength - key.length)
      serialised +=
        key.padEnd(maxKeyLength + 3, " ") +
        node.metadata[key] +
        "\n";
    }
    return serialised;
  }
}
