import axios from 'axios';

export const Domain = 'micro.mu';
const BaseURL = 'https://api.micro.mu/account/'

export default async function Call(path: string, params?: any): Promise<any> {
  return axios.post(BaseURL + path, params, { withCredentials: true });
}

export class User {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  paymentMethods: PaymentMethod[];
  subscriptions: Subscription[];
  roles: string[];

  constructor(args: any) {
    this.id = args.id;
    this.firstName = args.firstName || '';
    this.lastName = args.lastName || '';
    this.email = args.email || '';
    this.paymentMethods = (args.paymentMethods || []).map(p => new PaymentMethod(p));
    this.subscriptions = (args.subscriptions || []).map(p => new Subscription(p));
    this.roles = args.roles || [];
  }

  requiresOnboarding():boolean {
    if(this.roles.includes('admin')) return false;
    if(this.paymentMethods.length === 0) return true;
    if(this.subscriptions.length === 0) return true;
    return false
  }

  profileCompleted():boolean {
    if(this.firstName.length === 0) return false;
    if(this.lastName.length === 0) return false;
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
  userId: string;
  type: string;
  cardBrand: string;
  cardExpMonth: string;
  cardExpYear: string;
  cardLast4: string;
  default: boolean;

  constructor(args: any) {
    this.id = args.id;
    this.created = args.created;
    this.userId = args.userId;
    this.type = args.type;
    this.cardBrand = args.cardBrand;
    this.cardExpMonth = args.cardExpMonth;
    this.cardExpYear = args.cardExpYear;
    this.cardLast4 = args.cardLast4;
    this.default = args.default;
  }
}

export class Token {
  token: string;
  expires: Date;

  constructor(args: any) {
    this.token = args.token;
    this.expires = new Date(args.expires * 1000)
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