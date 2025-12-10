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
  - Go for all logic. Server-side rendered pages with just plain HTML + CSS.
  - Gorm as the ORM layer. Although the database usage will be super simple I'd like to have an abstraction layer to allow for extensibility. And I need something to do migrations which I'll need for the schema and seeding the data.
  - [Tachyons](https://tachyons.io/) for basic front-end styling. External style links only (not interested in hosting assets).
  - Docker builds.
  - Github Actions for CI/CD pipeline.
  - AWS for IaaS. I would like to learn a bit more Azure at this stage but just don't think there's time for it.
  - AWS AppRunner to host the application.
  - RDS Aurora serverless with PostgreSQL for the database.
- Not entirely sure how the PII aspect will fit in here. As a precaution I'll have each quote have a UUID primary key. We cannot avoid transmitting the quote data to the frontend, and there are no other tables. We're also not doing any analytics where we'd need to anonymise or just use the ID.
  - How will random selection of a single row work? This seems simple on the surface but could be surprisingly complex.

## Non goals

- Upon re-reading the specs there is no need for querying. Elastic/Opensearch was completely out of the question due to time constraints, but now I am realising I don't even need `tsvector/tsquery` or even `LIKE` queries.
