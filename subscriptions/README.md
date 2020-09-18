# Subscription Service

This is the Subscription service. It manages the lifecycle of subscriptions and does most of its work through the payment provider. Services should (probably) interact with this service rather than the payment provider service.

In our modelling, a subscription is mapped to an M3O account and is basically 1:1. So even if you're an additional user invited to someone else's namespace you will have your own subscription object.

An M3O subscription may have multiple Stripe subscriptions due to the way we model the charging. For example a typical M3O subscription will have at least 1 Stripe subscription, the core M3O platform subscription. It may also have N "additional user" subscriptions and "additional service" ones too. Everything we charge for will end up as a separate Stripe subscription, so in the future you may pay for extra resource and that will be an additional Stripe subscription but it will end up mapping to this single M3O subscription object. 