import axios from 'axios';

// const BaseURL = 'http://localhost:8080/notes/'
const BaseURL = 'https://api.micro.mu/notes/'

export default async function Call(path: string, params?: any): Promise<any> {
  return axios.post(BaseURL + path, params)
}

export class Note {
  id: string;
  title: string;
  text: string;
  created: Date;

  constructor(args: any) {
    this.id = args.id;
    this.title = args.title;
    this.text = args.text;
    this.created = new Date(parseInt(args.created) * 1000);
  }
}