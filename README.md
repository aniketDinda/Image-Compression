# Image-Compression-using-pub-sub

# Golang Message Queuing System

## Description

This is a Golang-based message queuing system that utilizes MongoDB as the database for storing product information and RabbitMQ for establishing communication between the producer and consumer components. The system allows users to add products, which are then processed by the consumer to download and compress product images.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Installation

1. Clone the repository to your local machine:

   ```bash
   git clone <repository-url>

2. Change the directory to the project folder:

   ```bash
   cd <project-folder>
   
3. Install the project dependencies using go mod:

   ```bash
   go mod tidy


## Usage

Before using the application, make sure to set up MongoDB and RabbitMQ on your local machine or server. Update the configuration files with the appropriate connection details.

## Project Structure
The project has the following structure:

main.go: The entry point of the application.
models/: Contains the data models for users and products.
routes/: Defines the API routes.
helpers/: Contains helper functions for downloading and compressing images.
database/: Handles database connections.
controller/: Defines the application's controllers.

## API Endpoints
POST /user/new: Create a new user.
POST /product/add: Add a new product.
GET /product/view: View all products.

## Configuration
.env: Store environment variables, including MongoDB and RabbitMQ connection details.

## Testing
Run the tests for the application, including unit tests, integration tests, and benchmark tests:

## Contributing
Contributions are welcome! Feel free to open issues or pull requests for any improvements or bug fixes.

   
