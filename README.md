# URL Shortener Project

This project is a simple URL shortener built using Go and the Gin framework. It allows users to create short URLs that redirect to original URLs.

## Project Structure

```
url-shortener
├── main.go                # Entry point of the application
├── go.mod                 # Module dependencies
├── go.sum                 # Module dependency checksums
├── controllers            # Contains the URL controller
│   └── url_controller.go  # Handles URL creation and retrieval
├── models                 # Contains the data models
│   └── url.go             # Defines the URL struct
├── routes                 # Contains route definitions
│   └── routes.go          # Sets up application routes
├── services               # Contains business logic
│   └── shortener_service.go # Handles URL shortening logic
├── utils                  # Contains utility functions
│   └── helper.go          # Assists with URL operations
├── config                 # Contains configuration settings
│   └── config.go          # Holds app configuration
└── README.md              # Project documentation
```

## Setup Instructions

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/url-shortener.git
   cd url-shortener
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the application:**
   ```bash
   go run main.go
   ```

4. **Access the API:**
   The API will be available at `http://localhost:8080`. You can use tools like Postman or curl to interact with the endpoints.

## Usage Examples

- **Create a Short URL:**
  Send a POST request to `/shorten` with a JSON body containing the original URL.
  
- **Retrieve Original URL:**
  Send a GET request to `/url/{shortUrl}` to get the original URL associated with the short URL.

## Contributing

Feel free to submit issues or pull requests if you have suggestions or improvements for the project.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.