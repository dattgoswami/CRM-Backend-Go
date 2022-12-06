# CRM Backend Project

## Getting Started

This repo contains a basic GO app to get started with constructing an API using GO. To get started, clone this repo and run `go run main.go` in your terminal at the project root.

# API Requirements

The company stakeholders want to create a CRM backend. You have been tasked with building the API that will support this application, and your coworker is building the frontend.

These are the notes from a meeting with the frontend developer that describe what endpoints the API needs to supply, as well as data shapes the frontend and backend have agreed meet the requirements of the application.

## API Endpoints

#### Customers

- A SHOW route: [GET] '/customers/{id}'
- An INDEX route: [GET] '/customers'
- A CREATE route: [POST] '/customers'
- An UPDATE route: [PATCH] '/customers/{id}'
- A DELETE route: [DELETE] '/customers/{id}'

## Data Shapes

#### Customers

- Id
- Name
- Role
- Email
- Phone
- Contacted

## Required Technologies

Your application must make use of the following libraries:

- Postgres instance
- gorrila mux (use `go get` to install this and lib pq)
- lib pq

## Once the project is up and running we can test it using postman(/curl)

1. Send a GET request to url [http://0.0.0.0:3000/customers/]
2. Send a GET request to url [http://0.0.0.0:3000/customers/1222]
3. Send a POST request to url [http://0.0.0.0:3000/customers] with the body containing the following raw json

```
   {
   "id": 1456,
   "name": "Jack",
   "role": "Product Manager",
   "email": "jack@example.com",
   "phone": "67898989",
   "contacted": false
   }
```

4. Send a PATCH request to url [http://localhost:3000/customers/1456] with the body containing the following raw json

```
   {
   "name": "Jackie",
   "role": "Product Manager",
   "email": "jackie@example.com",
   "phone": "67898989",
   "contacted": true
   }
```

5. Send a DELETE request to url [http://localhost:3000/customers/1456]

## Database Creation

### 1. Setup Database

It is time to create the database and the user with access privileges.

```
psql -U postgres
CREATE DATABASE crm_customers;
CREATE USER example_admin WITH PASSWORD 'somepassword';
GRANT ALL PRIVILEGES ON DATABASE crm_customers TO example_admin;
\c crm_customers
CREATE EXTENSION citext;
CREATE TABLE customers(
id INT PRIMARY KEY,
name VARCHAR(30),
role TEXT,
email CITEXT UNIQUE,
phone TEXT,
contacted BOOLEAN);
GRANT ALL PRIVILEGES ON TABLE customers TO example_admin;
```

### 2. Insert Data to the table

```
INSERT INTO customers VALUES (1367, 'Rick', 'Software Engineer', 'rick@example.com', '55198989', true);
INSERT INTO customers VALUES (1222, 'John', 'Data Analyst', 'john@example.com', '75198888', false);
INSERT INTO customers VALUES (3243, 'Ron', 'Data Scientist', 'ron@example.com', '45201787', true);
```

## References:

1. https://drstearns.github.io/tutorials/gojson/
2. https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/
3. https://hugo-johnsson.medium.com/rest-api-with-golang-mux-mysql-c5915347fa5b
4. https://dba.stackexchange.com/questions/68266/what-is-the-best-way-to-store-an-email-address-in-postgresql
5. https://techinscribed.com/different-approaches-to-pass-database-connection-into-controllers-in-golang/
