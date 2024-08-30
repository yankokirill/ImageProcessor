# Simple Task Manager Server

_A minimalistic server setup for managing tasks._

## Overview

This project sets up a simple Go-based Task Manager Server, which is primarily designed for demonstration purposes.
It includes a basic API, although the server itself doesn't perform any real tasks beyond providing a framework
for potential expansions.

## Getting Started

### Building and Running the Docker Image

To build the Docker image for the Task Manager Server and run tests, follow these steps:

**Build the Docker image and run tests using Make:**

    ```bash
    make tests
    ```


### Shooter API

The Shooter API allows you to apply filters to images. Here’s how to use it:

1. **Navigate to the `shooter` directory and install dependencies:**

    ```bash
    cd shooter
    pip install -r requirements.txt
    ```

2. **Run the `shoot.py` script to apply a filter to an image:**

    ```bash
    python shoot.py <FilterName> [<Parameters>]
    ```

    - **`<FilterName>`**: Name of the filter (e.g., `Negative`, `Blur`, `Sharpen`).
    - **`<Parameters>`**: Optional filter parameters in JSON format (e.g., `{"sigma": 5.0}`).

The processed images will be saved in the 'shooter/results' directory.

**Examples:**

```bash
python shoot.py Negative
```

```bash
python shoot.py Blur '{"sigma": 5.0}'
```

```bash
python shoot.py Sharpen '{"sigma": 5.0}'
```

### Accessing the API Documentation

Once the server is up and running, you can view the API documentation by navigating to:

http://localhost:8000/swagger/index.html#/

This link will direct you to the Swagger UI, where you can explore the available API endpoints.