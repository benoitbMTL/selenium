package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "runtime"
    "github.com/tebeka/selenium"
)

// BrowserRequest struct for parsing the JSON request
type BrowserRequest struct {
    Browser    string `json:"browser"`
    Iterations int    `json:"iterations"`
}

func main() {

	fs := http.FileServer(http.Dir("."))
    http.Handle("/", fs)

    http.HandleFunc("/openBrowser", openBrowserHandler)

    fmt.Println("Server started on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func openBrowserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req BrowserRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Calling the openBrowser function
    err = openBrowser(req.Browser, req.Iterations)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Process launched for browser %s with %d iterations", req.Browser, req.Iterations)
}

func openBrowser(browserName string, iterations int) error {
    log.Printf("Starting browser session for: %s", browserName)

    var driverPath string
    var err error

    // Determine the operating system to choose the right driver
    osType := runtime.GOOS
    log.Printf("Detected operating system: %s", osType)

    switch browserName {
    case "chrome":
        driverPath = "./webdrivers/chromedriver"
        if osType == "windows" {
            driverPath += ".exe"
        }
        log.Printf("Using Chrome driver at: %s", driverPath)
    case "firefox":
        driverPath = "./webdrivers/geckodriver"
        if osType == "windows" {
            driverPath += ".exe"
        }
        log.Printf("Using Firefox driver at: %s", driverPath)
    case "edge":
        driverPath = "./webdrivers/msedgedriver"
        if osType == "windows" {
            driverPath += ".exe"
        }
        log.Printf("Using Edge driver at: %s", driverPath)
    default:
        errMsg := fmt.Sprintf("Unsupported browser: %s", browserName)
        log.Println(errMsg)
        return fmt.Errorf(errMsg)
    }

    log.Println("Initializing WebDriver session...")
    caps := selenium.Capabilities{"browserName": browserName}
    wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515))
    if err != nil {
        log.Printf("Failed to create WebDriver session: %v", err)
        return err
    }
    defer func() {
        log.Println("Terminating WebDriver session.")
        wd.Quit()
    }()

    for i := 0; i < iterations; i++ {
        log.Printf("Attempt %d: Navigating to http://www.google.com", i+1)
        err = wd.Get("http://www.google.com")
        if err != nil {
            log.Printf("Failed to navigate on attempt %d: %v", i+1, err)
            return err
        }
        log.Printf("Successfully navigated to http://www.google.com on attempt %d", i+1)
        // Insert additional interactions with the page here.
        // For example, searching, clicking links, etc.
    }

    log.Println("Completed all browser interactions successfully.")
    return nil
}
