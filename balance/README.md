# Balance Service

This is the Balance service. It tracks the balance for a customer in particular for the following scenarios
- Incrementing the balance when the customer tops up their account
- Decrementing the balance when the customer uses the API

This means that it is designed for potentially high reqs/sec per customer - e.g. you'd expect high traffic customers to cause high traffic to the balance service. With this in mind it trades off accuracy for speed - it is not a ledger of debits and credits. In the event of failure it could under count usage which would mean that some API usage is unaccounted for which is a failure in favour of the customer. The record of customer payments is stored in our payments provider and as such, the danger of missing a customer payment is reduced.
