# Golang Ecommerce API

This project is a simple ecommerce API built with Golang, focusing on user authentication, product management, and cart-based order creation.

## Business Rules

### Implemented Features

- **User Authentication**
  - Users can register with their first name, last name, email, and password.
  - Users can log in using their email and password.
  - Upon successful login, users receive a JWT token, which must be used for authenticated operations.

- **Product Management**
  - Administrators can create products with details such as name, description, price, and quantity.
  - Products are stored in the database and can be retrieved for display in the store.

- **Cart & Order Management**
  - Users can add products to a cart and proceed to checkout.
  - During checkout, the system checks if the requested quantities are available.
  - If the stock is sufficient, an order is created, and the products are deducted from inventory.

### Planned Features

- **Product Management Enhancements**
  - [ ] Administrators will be able to update existing products.
  - [ ] Administrators will be able to delete products.
  - [ ] Users will be able to retrieve detailed product information by product ID.

- **Order Management Enhancements**
  - [ ] Users will be able to view their order history.
  - [ ] Administrators will be able to update the status of an order (e.g., pending, shipped, delivered).
  - [ ] Users will be notified of order status changes via email.

- **User Account Management**
  - [ ] Users will be able to view and update their profile information.
  - [ ] Users will be able to change their password.

- **Payment Integration**
  - [ ] The system will integrate with payment gateways like Stripe or PayPal to process payments during checkout.
  - [ ] Orders will only be created after successful payment processing.

- **Reporting and Analytics**
  - [ ] Administrators will be able to generate sales reports.
  - [ ] The system will track and report inventory levels to prevent stockouts.

- **Notification System**
  - [ ] Implement email notifications for order confirmations and status updates.
  - [ ] Users will receive promotional emails based on their purchase history.

## Running the Project

### Prerequisites

- **Go**: Version 1.22.6 or later.
- **Docker**: For running the MySQL database.
- **Make**: Installed on your system (Linux and MacOS typically have this by default; for Windows, you can use **Make for Windows** or **Git Bash**).

### Setting Up the Project

1. **Clone the repository**:

    ```bash
    git clone https://github.com/joshbarros/golang-ecommerce-api.git
    cd golang-ecommerce-api
    ```

2. **Set up the environment variables**:
   - Create a `.env` file in the root directory with the following content:
     ```bash
     PUBLIC_HOST=http://localhost
     PORT=8080
     DB_USER=your_db_user
     DB_PASSWORD=your_db_password
     DB_HOST=127.0.0.1
     DB_PORT=3306
     DB_NAME=golang-ecommerce-api
     JWT_EXP=604800 # 7 days in seconds
     JWT_SECRET=please-dont-tell-anyone
     ```

3. **Start MySQL using Docker**:

    ```bash
    docker-compose up -d
    ```

4. **Apply database migrations**:

    ```bash
    make migrate-up
    ```

5. **Build and run the application**:

    - **On Linux**:

      ```bash
      make run
      ```

    - **On Windows** (using Git Bash or similar):

      ```bash
      make run
      ```

6. **Running Tests**:

    - To run the test suite, use the following command:

      ```bash
      make test
      ```

## Makefile Explanation

### Makefile Targets

- **build**: Compiles the Go application and outputs the binary to the `bin` directory.
- **test**: Runs the entire test suite.
- **run**: Builds the application and runs it.
- **migration**: Creates a new migration file with the specified name.
- **migrate-up**: Applies all up migrations to the database.
- **migrate-down**: Rolls back the last migration applied to the database.

### Running the Project on Linux / MacOs

- **Run the application**:

  ```bash
  make run
  ```

## Technical Explanation of Libraries

### Routing: `gorilla/mux`

Gorilla Mux is a powerful URL router and dispatcher for Golang. It allows you to define complex routing rules, including variables in the URL, regex-based matching, and more. This project uses Gorilla Mux to handle all incoming HTTP requests and route them to the appropriate handlers.

### Database: `go-sql-driver/mysql`

Go-MySQL-Driver is a MySQL driver for Go's `database/sql` package. It provides a robust and efficient way to connect to MySQL databases, execute queries, and manage transactions. This project uses this driver to interact with the MySQL database for storing and retrieving data.

### JWT Authentication: `golang-jwt/jwt`

Golang JWT is a widely used library for creating and verifying JSON Web Tokens (JWTs) in Go. JWTs are used for securely transmitting information between parties as a JSON object. This project uses JWTs for user authentication, ensuring that only authorized users can access certain endpoints.

### Environment Variables: `joho/godotenv`

Godotenv is a Go library that loads environment variables from a `.env` file into your Go application. This makes it easier to manage configuration settings like database credentials and secret keys. This project uses Godotenv to manage environment-specific configurations.

### Validation: `go-playground/validator`

Go-Playground Validator is a Go library for validating structs and fields based on tags. It provides a wide range of validation rules and allows you to define custom validation logic. This project uses this library to validate incoming requests, ensuring that data conforms to expected formats and rules before processing.
