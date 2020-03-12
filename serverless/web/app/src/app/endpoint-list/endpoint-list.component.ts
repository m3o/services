import { Component, OnInit, Input } from "@angular/core";
import * as types from "../types";
import { ServiceService } from "../service.service";

@Component({
  selector: "app-endpoint-list",
  templateUrl: "./endpoint-list.component.html",
  styleUrls: ["./endpoint-list.component.css"]
})
export class EndpointListComponent implements OnInit {
  @Input() services: types.Service[] = [];
  @Input() endpointQuery: string = "";
  @Input() selectedVersion: string = "";

  constructor(private ses: ServiceService) {}

  ngOnInit() {
    this.regenJSONs();
  }

  ngOnCange() {
    
    this.regenJSONs();
  }

  regenJSONs() {
    this.services.forEach(service => {
      service.endpoints.forEach(endpoint => {
        endpoint.requestJSON = this.valueToJson(endpoint.request, 1);
      });
    });
  }

  callEndpoint(service: types.Service, endpoint: types.Endpoint) {
    this.ses
      .call({
        endpoint: endpoint.name,
        service: service.name,
        address: service.nodes[0].address,
        method: "POST",
        request: endpoint.requestJSON
      })
      .then(rsp => {
        endpoint.responseJSON = rsp;
      });
  }

  valueToString(input: types.Value, indentLevel: number): string {
    if (!input) return "";

    const indent = Array(indentLevel).join("    ");
    const fieldSeparator = `,\n`;

    if (input.values) {
      return `${indent}${input.type} ${indentLevel == 1 ? "" : input.name} {
${input.values
  .map(field => this.valueToString(field, indentLevel + 1))
  .join(fieldSeparator)}
${indent}}`;
    } else if (indentLevel == 1) {
      return `${indent}${input.name} {}`;
    }

    return `${indent}${input.type} ${input.name}`;
  }

  // This is admittedly a horrible temporary implementation
  valueToJson(input: types.Value, indentLevel: number): string {
    const typeToDefault = (type: string): string => {
      switch (type) {
        case "string":
          return '""';
        case "int":
        case "int32":
        case "int64":
          return "0";
        case "bool":
          return "false";
        default:
          return "{}";
      }
    };

    if (!input) return "";

    const indent = Array(indentLevel).join("    ");
    const fieldSeparator = `,\n`;
    if (input.values) {
      return `${indent}${indentLevel == 1 ? "{" : '"' + input.name + '": {'}
${input.values
  .map(field => this.valueToJson(field, indentLevel + 1))
  .join(fieldSeparator)}
${indent}}`;
    } else if (indentLevel == 1) {
      return `{}`;
    }

    return `${indent}"${input.name}": ${typeToDefault(input.type)}`;
  }

  // code editor
  coptions = {
    theme: "vs-dark",
    language: "json"
  };

  pickVersion(services: types.Service[]): types.Service[] {
    return services.filter(s => {
      return s.version == this.selectedVersion;
    });
  }
}
