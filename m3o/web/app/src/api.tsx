export interface User {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  roles: string[];
  me?: boolean;
}

export interface EnvVar {
  id?: string;
  key: string;
  value: string;
  service: string;
  secret?: boolean;
}
