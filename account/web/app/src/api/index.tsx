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
  payment_methods: PaymentMethod[];
  subscriptions: Subscription[];
  roles: string[];
  invite_code: string;
  invite_verified: boolean;

  constructor(args: any) {
    this.id = args.id;
    this.first_name = args.first_name || '';
    this.last_name = args.last_name || '';
    this.email = args.email || '';
    this.payment_methods = (args.payment_methods || []).map(p => new PaymentMethod(p));
    this.subscriptions = (args.subscriptions || []).map(p => new Subscription(p));
    this.roles = args.roles || [];
    this.invite_code = args.invite_code;
    this.invite_verified = args.invite_verified;
  }

  requiresOnboarding():boolean {
    if(this.roles.includes('admin')) return false;
    if(this.payment_methods.length === 0) return true;
    if(this.subscriptions.length === 0) return true;
    return false
  }

  profileCompleted():boolean {
    if(this.first_name.length === 0) return false;
    if(this.last_name.length === 0) return false;
    if(!this.invite_verified) return false;
    return true
  }
}

export class Subscription {
  id: string;
  plan: Plan;

  constructor(args: any) {
    this.id = args.id;
    this.plan = new Plan(args.plan);
  }
}

export class PaymentMethod {
  id: string;
  created: string;
  user_id: string;
  type: string;
  card_brand: string;
  card_exp_month: string;
  card_exp_year: string;
  card_last_4: string;
  default: boolean;

  constructor(args: any) {
    this.id = args.id;
    this.created = args.created;
    this.user_id = args.user_id;
    this.type = args.type;
    this.card_brand = args.card_brand;
    this.card_exp_month = args.card_exp_month;
    this.card_exp_year = args.card_exp_year;
    this.card_last_4 = args.card_last_4;
    this.default = args.default;
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

export class Plan {
  id: string;
  name: string;
  amount: number;
  interval: string;
  available: boolean;

  constructor(args: any) {
    this.id = args.id;
    this.name = args.name;
    this.amount = parseInt(args.amount) || 0;
    this.interval = args.interval;
    this.available = args.available;
  }
}