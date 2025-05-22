# Research

## What is an ISA and who can have one?

- Tax free investment account
- Annual deposit limit of £20,000
- Only available to customers who:
    - are over the age of 18
    - Live in the UK (or have tax residency in the UK)
    - Have a valid National Insurance number

## What information do we need to on-board a retail customer?

- Full name
- Full address
- Contact information (email, phone number)
- National Insurance number
- Date of birth
- Country of tax residency

## Verifying suitability for a stocks and shares ISA account

[https://www.gov.uk/individual-savings-accounts#who-can-open-an-isa](https://www.gov.uk/individual-savings-accounts#who-can-open-an-isa)

NI number, DOB and country of tax residency will be used to determine if the customer can open an ISA

NI numbers follow the pattern: 2 letters followed by 6 numbers, then A, B, C or D. This is suitable for frontend form validation.

*In my code I make the assumption that there is a service inside Cushon to verify NI numbers.*

Country of tax residency is used over address as ISAs are also available to members of the armed forces and crown servants who may be stationed abroad. **This is my understanding, there maybe some nuance here**

## Key functionality of an ISA account

### Considered for the assignment

- Open an ISA account
- Pay money into the account
- Invest in a fund
- View transactions

### Future functionality

- Transfer money between funds
- Set up ongoing payments into a fund
- Withdraw money
- Close the account

## Scale

- Total number of working age adults (18-64) ~42 million [link](https://www.ibisworld.com/uk/bed/population-aged-18-to-64-years/44240/)
- Total number of adults with an ISA ~20 million [link](https://www.gov.uk/government/statistics/annual-savings-statistics-2024/commentary-for-annual-savings-statistics-september-2024)
- HMRC Registered list of ISA providers ~500 [link](https://www.gov.uk/government/publications/list-of-individual-savings-account-isa-managers-approved-by-hmrc/registered-individual-savings-account-isa-managers)

Based on the above we should consider a total number of customers on the level of 100,000s - 1,000,000, this would also align with the current number of employees (although they are not direct customers and instead registered in bulk by their employer.

The number of daily active users would likely be fairly low (relative to total users) as once a user has registered and invested, they will likely not be logging in frequently to check (assume perhaps once a month/quarter).

We would still have to consider spikes which could happen at the end/beginning of the month (payday), around the end of the financial year and during ad campaigns.

The system therefore should be able to accommodate these data spikes, a well written Go service on reasonable hardware should be able to handle this level of traffic but it would still be a good idea to have at least a single replica and distribute traffic with a load balancer, not only will this easily accommodate the traffic, it will also make the system more resilient.

We might also want to consider horizontal scaling of the database through sharding (specifically for the retail savings service as this keeps track of transactions and could grow very quickly to 10s of millions of rows as customers on board and make investments).
