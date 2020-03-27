import { User, PaymentMethod } from "../api";

const SET_USER = 'SET_USET';
const ADD_PAYMENT_METHOD = 'ADD_PAYMENT_METHOD';
const REMOVE_PAYMENT_METHOD = 'REMOVE_PAYMENT_METHOD';

interface Action {
  type: string;
  user?: User;
  paymentMethod?: PaymentMethod;
}

interface State {
  user?: User;
}

export function setUser(user: User): Action {
  return { type: SET_USER, user };
}

export function addPaymentMethod(pm: PaymentMethod): Action {
  return { type: ADD_PAYMENT_METHOD, paymentMethod: pm };
}

export function removePaymentMethod(pm: PaymentMethod): Action {
  return { type: REMOVE_PAYMENT_METHOD, paymentMethod: pm };
}

const defaultState: State = {
  // user: new User({ id: "aeecc8a8-eb28-4b17-bb19-7d103f125e66", email: "ben@micro.mu" }),
};
export default function(state = defaultState, action: Action): State {
  switch (action.type) {
    case SET_USER: 
      return { ...state, user: action.user! };
    case ADD_PAYMENT_METHOD:
      var user = new User({
        ...state.user, paymentMethods: [
          ...state.user!.paymentMethods,
          action.paymentMethod,
        ],
      });

      return { ...state, user };
    case REMOVE_PAYMENT_METHOD:
      user = new User({
        ...state.user, paymentMethods: [
          ...state.user!.paymentMethods.filter(p => p.id !== action.paymentMethod!.id),
        ],
      });

      return { ...state, user };
    default:
      return state;
  }
}