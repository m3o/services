import axios from 'axios';

const BaseURL = 'http://localhost:8080/account/'
// const BaseURL = 'https://api.micro.mu/account/'

export default async function Call(path: string, params?: any): Promise<any> {
  return axios.post(BaseURL + path, params, { withCredentials: true });
}

export class User {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  username: string;

  constructor(args: any) {
    this.id = args.id;
    this.firstName = args.firstName;
    this.lastName = args.lastName;
    this.email = args.email;
    this.username = args.username;
  }
}