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
  username: string;
  paymentMethods: PaymentMethod[];
  subscriptions: Subscription[];
  roles: string[];

  constructor(args: any) {
    this.id = args.id;
    this.firstName = args.firstName;
    this.lastName = args.lastName;
    this.email = args.email;
    this.username = args.username;
    this.paymentMethods = (args.paymentMethods || []).map(p => new PaymentMethod(p));
    this.subscriptions = (args.subscriptions || []).map(p => new Subscription(p));
    this.roles = args.roles || [];
  }

  requiresOnboarding():boolean {
    // testing
    return this.email === 'ben@micro.mu';

    if(this.roles.includes('admin')) return false;
    if(this.paymentMethods.length === 0) return true;
    if(this.subscriptions.length === 0) return true;
    return false
  }
}

export class Subscription {
  constructor(args: any) {}
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

  constructor(args: any) {
    this.id = args.id;
    this.created = args.created;
    this.userId = args.userId;
    this.type = args.type;
    this.cardBrand = args.cardBrand;
    this.cardExpMonth = args.cardExpMonth;
    this.cardExpYear = args.cardExpYear;
    this.cardLast4 = args.cardLast4;
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
