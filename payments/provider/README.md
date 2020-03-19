# Payments Provider

A payments provider is a third party service which facilitates collecting / managing payments.

Payment providers are registered as services with the format "go.micro.service.payments.*".

To use the package, generate a new provider using the `provider.NewProvider("stripe")` method. If a provider with this identifier is not found, or has not been registered to service discovery then an error will be returned.