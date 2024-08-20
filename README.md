# **URL Shortener and Manager CLI**

A simple and extensible URL shortener CLI application written in Go. This project allows you to shorten URLs, track how many times a shortened URL has been accessed, and automatically expire URLs after a configurable time. The data is persisted to a file, so URLs and statistics remain consistent across program restarts.

## **Features**
- **URL Shortening**: Convert long URLs into shorter, more manageable links.
- **URL Redirection**: Redirect users to the original URL when they access the shortened link.
- **Access Statistics**: Track the number of times each short URL is accessed.
- **URL Expiry**: Automatically expire URLs after a configurable time period.
- **File Persistence**: URLs and their statistics are saved to a file and reloaded on program restart.
- **CLI Interface**: Easy-to-use command-line interface for managing URLs.

## **Go Features & Packages Used**

### **1. Goroutines and Concurrency**
- **Goroutines:** goroutines are utilized for handling multiple HTTP requests concurrently
- **Synchronization:** The `sync.RWMutex` is used to safely manage concurrent access to shared data (the URL map) by multiple goroutines. The `RWMutex` provides a mechanism for synchronizing access, where multiple readers can read from the shared data simultaneously, but only one writer can write to it at a time.

### **2. `net/http` Package**
- To define endpoints for URL creation and redirection.
  
### **3.`encoding/json` Package**
- To easily marshal (convert Go structs into JSON) and unmarshal (parse JSON into Go structs) data

### **4. `time` Package**
- To manage the expiration of URLs. By setting expiry times using `time.Duration` and `time.Now()`, this project can determine whether a shortened URL has expired and prevent its use after the specified duration.

### **5. `math/rand` Package**
- To generate random strings for creating unique short URLs. By seeding the random number generator with `time.Now().UnixNano()`, the project ensures that the generated URLs are different for each request.

### **6. `github.com/gorilla/mux` Package**
- To handle URL parameters and create RESTful APIs.

### **7. File I/O**
- To read and write URL data to a file. This ensures that URLs and their metadata(such as access counts and expiration times are persisted.

### **8. Modular Code Structure**
- The project is structured in a modular way, with functions defined for specific tasks such as generating short URLs, handling HTTP requests, and managing data storage

### **9. Unit Testing with `testing` Package**
- The `testing` package is used to write and run unit tests for the project's core functionality. The tests ensure that the URL generation, storage, redirection, and expiry logic work as expected. This helps maintain code quality and reliability.

## **Installation**

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/RaajeevChandran/url-shortener-go.git
   cd url-shortener-go

2. **Build the Application:**
   ```bash
   go run main.go

## **To run the tests**
```bash
go test
```

## **Configuration**

- **URL Expiry Duration:** The default expiry time for URLs is set to 24 hours. You can change the `expiryDuration` constant in the code to adjust this time.
- **Auto-save Interval:** The URL data is auto-saved every 5 minutes by default. This interval can be adjusted by modifying the `saveInterval` constant.

## **Usage**

1. **Create Short URL:**
   - Choose the option to create a new short URL.
   - Enter the URL you wish to shorten.
   - The application will generate a shortened URL.

2. **Redirect to Original URL:**
   - Enter the short URL to be redirected to the original URL.
   - The application will display the original URL if the short URL is valid and not expired.

3. **View Statistics:**
   - Enter the short URL to view the number of times it has been accessed.

4. **Exit the Application:**
   - Choose the exit option to safely terminate the program, ensuring data is saved.

### **Example:**

```bash
  Choose an option:

  Create Short URL
  Redirect to Original URL
  View Statistics
  Exit
  Enter your choice: 1
  Enter the URL: https://www.example.com
  Short URL: abc123

  Enter your choice: 3
  Enter the Short URL to view statistics: abc123
  Short URL abc123 has been accessed 0 times.
```


