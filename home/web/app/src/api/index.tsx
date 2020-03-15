import axios from 'axios';

const BaseURL = 'http://dev.micro.mu:8080/home/'
// const BaseURL = 'https://api.micro.mu/home/'

export default async function Call(path: string, params?: any): Promise<any> {
  return axios.post(BaseURL + path, params, { withCredentials: true });
}

export class User {
  firstName: string;
  lastName: string;

  constructor(args: any) {
    this.firstName = args.firstName;
    this.lastName = args.lastName;
  }
}