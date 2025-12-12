# quotes-database

## Top-level goal

- Web application seeded with famous quotes
- Uses an external database (as a service)
- When site is visited, a random quote is selected and displayed
- Treats data as PII
- Highly available
- Less than 5 hours time to implement, in 5 days.

## Sub goals

- Simplicity is key. Only have five hours so there is no time for anything fancy. I will use:
  - Go for all logic. Server-side rendered pages with just plain HTML. No CSS.
  - Gorm as the ORM layer. Although the database usage will be super simple I'd like to have an abstraction layer to allow for extensibility. And I need something to do migrations which I'll need for the schema and seeding the data.
  - Docker builds.
  - Github Actions for CI/CD pipeline.
  - AWS for IaaS. I would like to learn a bit more Azure at this stage but just don't think there's time for it.
  - AWS AppRunner to host the application.
  - RDS Aurora serverless with PostgreSQL for the database.
- Not entirely sure how the PII aspect will fit in here. As a precaution I'll have each quote have a UUID primary key. We cannot avoid transmitting the quote data to the frontend, and there are no other tables. We're also not doing any analytics where we'd need to anonymise or just use the ID.
  - How will random selection of a single row work? This seems simple on the surface but could be surprisingly complex.

## Non goals

- Upon re-reading the specs there is no need for querying. Elastic/Opensearch was completely out of the question due to time constraints, but now I am realising I don't even need `tsvector/tsquery` or even `LIKE` queries.
- No CSS.
- Noticed that [Atlas](https://gorm.io/docs/migration.html#Atlas-Integration) would be a better database migration option than what GORM provides natively but I don't think there's time to worry about that.

## Work log

- Created a simple program structure in Go, table schema, just enough to get an application started.
- Created the Dockerfile.
- Created a basic CI/CD pipeline and Dependabot configuration (Copilot aided).
- Generated the Terraform boilerplate and a basic Github Actions pipeline for infrastructure (Copilot aided).
- Generated a Terraform VPC networking boilerplate configuration (Copilot aided).
- Created the database cluster, using serverless Aurora (PG compatible).
- Created the AppRunner configuration and connected to the deployment workflow.
- Completed the AppRunner configuration with IAM roles and a secret for database access. App connects to the database now.
- Split code logic out to a separate struct for handling user requests, start planning database structure migration / data ingestion.