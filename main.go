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
    var driverPath string
    var err error

    // Determine the operating system to choose the right driver
    osType := runtime.GOOS // "windows", "linux", etc.
    switch browserName {
    case "chrome":
        driverPath = "./chromedriver"
        if osType == "windows" {
            driverPath += ".exe"
        }
    case "firefox":
        driverPath = "./geckodriver"
        if osType == "windows" {
            driverPath += ".exe"
        }
    case "edge":
        driverPath = "./msedgedriver"
        if osType == "windows" {
            driverPath += ".exe"
        }
    default:
        return fmt.Errorf("unsupported browser: %s", browserName)
    }

    // Note: For Edge, using selenium.NewRemote directly as NewEdgeDriverService might not be available.
    caps := selenium.Capabilities{"browserName": browserName}
    wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515))
    if err != nil {
        log.Printf("Failed to create WebDriver session: %v", err)
        return err
    }
    defer wd.Quit()

    for i := 0; i < iterations; i++ {
        log.Printf("Navigating to http://www.google.com, iteration %d", i+1)
        err = wd.Get("http://www.google.com")
        if err != nil {
            log.Printf("Failed to navigate: %v", err)
            return err
        }
        // Here, you can add more interactions with the page.
    }

    log.Println("Completed browser interactions")
    return nil
}
