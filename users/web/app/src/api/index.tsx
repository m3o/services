import axios from 'axios';

// const BaseURL = 'http://localhost:8080/users/'
const BaseURL = 'https://api.micro.mu/users/'

export default async function Call(path: string, token: string, params?: any): Promise<any> {
  const headers = { 'Authorization': 'Bearer ' + token };
  return axios.post(BaseURL + path, params, { headers });
}

export class User {
  id: string;
  firstName: string;
  lastName: string;
  email: string;

  constructor(args: any) {
    this.id = args.id;
    this.firstName = args.firstName;
    this.lastName = args.lastName;
    this.email = args.email;
  }
}