# Cushon Interview Scenario

*The task description can be found [here](https://github.com/jameswhoughton/cushon/blob/main/TASK.md)*

## Key Points

- Functionality for retail ISAs should be kept separate, I have however made some assumptions about core services that might be shared (e.g. funds and performing trades).
- Initially, customers can only select a single fund (however there should be scope to invest in multiple in the future)
- We need to keep a record of investments and make this visible to the user
    - Should also consider long term storage of transactional data (even after an account has been closed).
- I have been given a situation in which a customer wishes to invest £25,000 in a fund (this is more than the annual ISA allowance).
- I assume that different types of ISA will be offered to retail customers (e.g. Junior ISA and LISA). The different types of ISA have different requirements to create an account and restrictions when depositing into them. Our solution should account for these variations.

## My Assumptions

- The existing authentication system can be used to support retail customers.
- The same funds are available to both retail and employee customers (and that a service already exists to manage them).
- For this assignment, currency is assumed to be GBP, and all monetary values are stored to the nearest penny as ints.
- There is an existing service to perform trades.
- There is an existing service to process deposits (card payment/bank transfer etc.).
- There is an API gateway in place to serve the web UI and provide authentication.
    - Requests from the gateway contain a session with the current customerID

## Proposed Solution

*The research that led to my solution can be found [here](https://github.com/jameswhoughton/cushon/blob/main/RESEARCH.md).*

### Overview

I propose the addition of two new services: 

1. Retail Customer Service - Manage retail customers
2. Retail Account Service - Manage retail accounts.

Both services will use a REST API to expose their functionality to the wider system.

### Reasoning For Separate Services

Separating the logic for retail customers and accounts from the existing employee services has several benefits:

- Allows the customer/employees services to scale independently of one another.
- Provides a clear boundary for auditing/regulatory purposes.
- Allows entities to evolve independently (e.g. addition of new fields etc.)

There are however some trade-offs:

- Adds more complexity to the system (for example if an employee also opened a retail account we might want to be able to link the accounts).
- At least in the beginning there will likely be some duplication between the services (although this could be mitigated though the use of a shared package).

### Scale And Storage

- I have estimated a number of registered users on the scale of 100,000 - 1,000,000 with a much smaller fraction of daily active users (my reasoning can be found in RESEARCH.md). 
- While data scaling should be a consideration, MySQL should be able to handle this size of data (especially as the schema is fairly simple). If performance begins to become an issue, one option would be to create read replicas of the database. This would be preferable to sharding which would introduce a lot of complexity to the system.
- The application should be designed with horizontal scaling in mind in order to handle peaks in traffic (for example, around the end of the tax year).

A RDMS such as MySQL is a good choice of database for these services for the following reasons:

- Can handle large datasets.
- The data is relational so more suited to a relational database.
- As far as I'm aware, other services in Cushon rely on MySQL so that might be the better choice (while not critical for microservices to use the same technology, team familiarity is a benefit).

### Business considerations

There are different types of ISA which Cushon may wish to offer. The new service should be able to accommodate the following:

- Different requirements to open an account.
- Different restrictions on Depositing into an account.
- Different restrictions on withdrawing from an account.

### Service Architecture

I have begun to build out the retail account microservice using the service/repository pattern and TDD. The service is partially complete, handlers and the frontend still need to be added along with some additional tests.

- Services contain the business logic
- Repositories are only responsible for communicating with the MySQL database.
- Handlers format and pass information to and from the services.
- I am using docker to run a testing MySQL instance which is migrated/rolled back between tests.

In order to accommodate different types of ISA, I have defined a Service interface that can be implemented for each type of account, on top of this I have created a factory to return the correct service where required. For common service actions I have created some generic unexported functions, these can be used to compose the account-specific services to reduce code duplication, for example, `createAccount()` in `internal/account/service.go` is responsible for validating the account passed in and calling the underlying repository method to store the account the database (a group of actions common to all account types).



### Scenario: customer who wishes to deposit £25,000 into a Cushon ISA all into the Cushon Equities Fund

In this specific case the customer would not be permitted to deposit £25,000 into a Cushon ISA in a single transaction. Assuming this is a standard ISA the annual limit is £20,000. The only situation in which this would be possible is if they were transferring the money from an existing ISA with another provider, in this case they would be able to transfer the whole amount into a new Cushon ISA and invest it all into the Cuson Equities fund.

### Schema

My proposed DB schema can be found [here](https://raw.githubusercontent.com/jameswhoughton/cushon/refs/heads/main/schema.png).

## Future Enhancements

- Build out the frontend.
- I have focused on a basic ISA, and I would look to extend my solution to cover other ISA accounts.
- Store specific currency information.
- Consider notifications to the user (email/post).
- Explore the idea of a shared package for personal information types and validation (e.g. validating NI number)
- Customer personal information could be encrypted when inserted into the database, this would help to potentially reduce the impact of a data breach (direct DB access) at the cost of a slight performance hit.
- Consider pagination for transactions. My current solution limits the date range for returning transactions to 1 year.
- Consider permissions/admin routes for account management and reporting.
- Consider external ISA to Cushon ISA transfers.
- Introduce in-memory implementations of repositories (backed by contract tests) to improve test performance.
- Support cash balances (uninvested money in the ISA)


## Running Tests

To run the tests against the service, follow these steps.

1. Ensure Go and Docker are installed.
2. Run `docker compose up -d` at the project root.
3. Once the database is up and running, run `go test ./...`
