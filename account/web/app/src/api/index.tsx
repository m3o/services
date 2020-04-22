import axios from 'axios';

export const Domain = 'micro.mu';
const BaseURL = 'https://api.micro.mu/account/'

export default async function Call(path: string, params?: any): Promise<any> {
  return axios.post(BaseURL + path, params, { withCredentials: true });
}

export class User {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  roles: string[];
  invite_code: string;
  invite_verified: boolean;

  constructor(args: any) {
    this.id = args.id;
    this.first_name = args.first_name || '';
    this.last_name = args.last_name || '';
    this.email = args.email || '';
    this.roles = args.roles || [];
    this.invite_code = args.invite_code;
    this.invite_verified = args.invite_verified;
  }

  requiresOnboarding():boolean {
    if(this.roles.includes('admin')) return false;
    return !this.profileCompleted()
  }

  profileCompleted():boolean {
    if(this.first_name.length === 0) return false;
    if(this.last_name.length === 0) return false;
    if(!this.invite_verified) return false;
    return true
  }
}

export class Token {
  access_token: string;
  refresh_token: string;
  expiry: Date;
  created: Date;

  constructor(args: any) {
    this.access_token = args.access_token;
    this.refresh_token = args.refresh_token;
    this.created = new Date(args.created * 1000)
    this.expiry = new Date(args.expiry * 1000)
  }
}
