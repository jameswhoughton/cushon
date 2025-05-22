# Cushon Interview Scenario

## Description

Cushon already offers ISAs and Pensions to Employees of Companies (Employers) who have an existing arrangement with
Cushon. Cushon would like to be able to offer ISA investments to retail (direct) customers who are not associated with an
employer. Cushon would like to keep the functionality for retail ISA customers separate from it’s Employer based offering
where practical.

When customers invest into a Cushon ISA they should be able to select a single fund from a list of available options. Currently they will be restricted to selecting a single fund however in the future we would anticipate allowing selection of multiple options.

Once the customer’s selection has been made, they should also be able to provide details of the amount they would like to
invest.

Given the customer has both made their selection and provided the amount the system should record these values and allow these details to be queried at a later date.

As a specific use case please consider a customer who wishes to deposit £25,000 into a Cushon ISA all into the Cushon Equities Fund.

## Assignment

Please provide your solution to the above scenario in whatever form you feel is appropriate, using your preferred tools.
Please spend the amount of time you feel appropriate to showcase your abilities and knowledge.

Please be prepared to discuss during an interview:

- What you have done and why.
- The specific decisions you made about your solution.
- Any assumptions you have made in the solution you have presented.
- Any enhancements you considered but decided not to cover.

## Key points

- Functionality for retail ISAs should be kept separate, I have however made some assumptions about core services that might be shared (e.g. performing trades).
- Initially customers can only select a single fund (however there should be scope to invest in multiple in the future)
- We need to keep record of investments and make this visible to the user
- I have been given a situation where a customer wishes to invest £25,000 into a fund (this is more than the annual ISA allowance).
- There is no mention of the type of ISA (stocks and shares/LISA) so my solution will focus on a stocks and shares ISA but aim to be flexible enough to accommodate different ISA types at a later date.

## Customer journeys

### Scenario: A new customer can register and create an ISA

As a new Cushon customer:

- I can visit a registration page and fill in personal information.
- If the information is invalid or I unable to create an ISA account, I cannot proceed.
- If the information is valid I can choose a fund and an amount to invest.
    - I cannot invest more than £20,000 in a tax year.
    - I must pick a fund to invest in.
- I can review my application before submitting.
- After submission I see a success page with a link to login

### Scenario: An existing customer can view information about their accounts

As an existing Cushon retail customer who has authenticated:

- I can visit the accounts page, where I can see a list of my accounts.
    - I can click on an account to see more information.
- On the individual account page:
    - I can see which funds I am invested in and how much is invested.
    - I can see a table of recent transactions on the account.

### Scenario: An existing customer can change the fund in which their money is invested

As an existing Cushon retail customer with an invested ISA who is authenticated:

- I can visit my ISA account page:
    - I can select a new fund in which to invest all of my money.
    - I can confirm the selection.
    - I see a success/failure message.

## Proposed Solution

*The research that led to my solution can be found [here](https://github.com/jameswhoughton/cushon/blob/main/RESEARCH.md).*

I propose the addition of two new services: 'retailCustomerService' to manage retail customers and 'retailAccountService' which manages retail accounts (this is intentionally generic in anticipation of Cushon offering other types of savings accounts).

I did consider the possibility of reusing an existing 'customer' service (that may already exist for creating employee customers), however having a separate one has several benefits:

- Can scale independently of existing employeeCustomerService
- Creates a clear boundary for auditing
- Allows either each entity to evolve independently (e.g. addition of new fields etc.)

On the other hand it does add some more complexity and possibly some duplication (this could be slightly mitigated through the use of a shared package as mentioned in improvements). Without knowing more about the existing system it is hard to say more at this time.

Both services will be stateless and have a single replica of each service to increase resiliency.

Databases will use MySQL, this aligns with other services as well as being capable of handling big data, when designing the schema I have tried to optimise the table sizes as much as possible. 

I have noted that the 'account_transactions' table will likely be the fastest growing table and therefore it is worth considering splitting the data (either by partitioning the table or sharding the database). As transactional information should be retained (indefinitely?) table partitioning may not be suitable as the number of tables might grow to be unmanageable, so a sharding approach might be more apporpriate.

I have focused my efforts on building out these services in Go.

### ERD

My proposed DB schema can be found [here](https://raw.githubusercontent.com/jameswhoughton/cushon/refs/heads/main/erd.svg), the diagram has some notes explaining my decisions.

### Assumptions

- The existing Authentication system can be used to support retail customers.
- The same funds are available to both retail and employee customers (and that a service already exists to manage them).
- For the purposes of this assignment, currency is assumed to be GBP.
- All monetary values are stored to the nearest penny as ints.
- There is an existing service to perform trades.
- All dates are stored as UTC.

## Future enhancements

- Store specific currency information.
- Consider notifications to the user (email/post).
- Explore the idea of a shared package for personal information types and validation (e.g. validating NI number)
- Customer personal information could be encrypted when inserted into the database, this would help to potentially reduce the impact of a data breach (direct DB access) at the cost of slight performance hit.
- Explore the best approach to splitting the data

## Endpoints

### RetailCustomerService

#### POST /v1/customer

- Create a new retail customer

### RetailAccountService

#### POST /v1/customer/{customer_id}/account

- Create a new account for the customer
- Customer is limited to one cushon ISA

#### GET /v1/account/{account_id}

- List of invested funds
- List of recent transactions

#### POST /v1/account/{account_id}/invest

- Assign money from the account to a fund


