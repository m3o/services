import axios from 'axios';

let BaseURL = 'https://api.micro.mu/home/';

// Toggle dev enviroment
if(window.location.protocol !== 'https:') {
  BaseURL = 'http://dev.micro.mu/home/'; 
}

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

export class App {
  id: string;
  name: string;
  category: string;
  icon: string;

  constructor(args: any) {
    this.id = args.id;
    this.name = args.name;
    this.category = args.category || '';
    this.icon = args.icon || '';
  }
}