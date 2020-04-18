export interface User {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  roles: string[];
  me?: boolean;
}
