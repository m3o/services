import axios from 'axios';

export const Domain = 'dev.micro.mu';

// const BaseURL = 'http://dev.micro.mu:8080/account/'
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

  constructor(args: any) {
    this.id = args.id;
    this.firstName = args.firstName;
    this.lastName = args.lastName;
    this.email = args.email;
    this.username = args.username;
    this.paymentMethods = (args.paymentMethods || []).map(p => new PaymentMethod(p));
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
