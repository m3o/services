

export interface User {
  name: string;
  email: string;
  avatarURL: string;
  teamName: string;
  teamURL: string;
  organizationAvatarURL: string;
  login: string;
}

export interface Value {
  name: string;
  type: string;
  values: Value[];
}

export interface Endpoint {
  name: string;
  request: Value;
  response: Value;
  metadata: Map<string, string>;
  title: string; // does not exist yet
  description: string; // does not exist yet
  //
  requestJSON: string;
  responseJSON: string;
}

export interface Node {
  id: string;
  address: string;
  metadata: Map<string, string>;
  // @TODO come up with a way to wrap all types in conenience interfaces
  show?: boolean;
}

// ... slightly different version of Service...
// this should be unified
export interface DebugService {
  name: string;
  version: string;
  metadata: Map<string, string>;
  endpoints: Endpoint[];
  node: Node;
}

export interface Service {
  name: string;
  version: string;
  metadata: Map<string, string>;
  endpoints: Endpoint[];
  nodes: Node[];
}

// taken from https://github.com/micro/micro/blob/master/debug/log/proto/log.proto
export interface LogRecord {
  timestamp: number;
  metadata: Map<string, string>;
  message: string;
}

export interface DebugSnapshot {
  service: DebugService;
  // Unix timestamp, e.g. 1575561487
  started: number;
  // Uptime in seconds
  uptime: number;
  // Heap allocated in bytes (TODO: change to resident set size)
  memory: number;
  // Number of Goroutines
  threads: number;
  // GC Pause total in ns
  gc: number;
  // Total number of request
  requests: number;
  // Total number of errors
  errors: number;
  timestamp: number;
}

export interface Span {
  // the trace id
  trace: string;
  // id of the span
  id: string;
  // parent span
  parent: string;
  // name of the resource
  name: string;
  // time of start in nanoseconds
  started: number;
  // duration of the execution in nanoseconds
  duration: number;
  // associated metadata
  metadata: Map<string, string>;
  type: number;
}

export interface TraceSnapshot {
  service: DebugService;
  spans: Span[];
}

export interface EventService {
  name: string;
  version?: string;
  source?: string;
  type?: string;
  metadata?: Map<string, string>;
}

// Platform event
export interface Event {
  type: number;
  timestamp: number;
  metadata: { string: string };
  service: EventService;
}
