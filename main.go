package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "runtime"
	"os"
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

    err = openBrowser(req.Browser, req.Iterations)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Process launched for browser %s with %d iterations", req.Browser, req.Iterations)
}

func openBrowser(browserName string, iterations int) error {
    log.Printf("Starting browser session for: %s", browserName)

    var driverPath, driverPort string
    var cmd *exec.Cmd
    var err error

    osType := runtime.GOOS
    driverPort = "9515" // Assuming all drivers use the same port for simplicity, adjust if necessary

    switch browserName {
    case "chrome":
        driverPath = "./webdrivers/chromedriver"
        if osType == "windows" {
            driverPath += ".exe"
        }
        cmd = exec.Command(driverPath, "--port="+driverPort)
    case "firefox":
        driverPath = "./webdrivers/geckodriver"
        if osType == "windows" {
            driverPath += ".exe"
        }
        cmd = exec.Command(driverPath, "--port="+driverPort)
    case "edge":
        driverPath = "./webdrivers/msedgedriver"
        if osType == "windows" {
            driverPath += ".exe"
        }
        cmd = exec.Command(driverPath, "--port="+driverPort)
    default:
        errMsg := fmt.Sprintf("Unsupported browser: %s", browserName)
        log.Println(errMsg)
        return fmt.Errorf(errMsg)
    }

	cmd, err = startDriver(driverPath, driverPort)
	if err != nil {
		log.Printf("Failed to start %s driver: %v", browserName, err)
		return err
	}
    log.Printf("%s Driver Started", browserName)
    
	defer func() {
        if err := cmd.Process.Kill(); err != nil {
            log.Printf("Failed to kill %s Driver process: %v", browserName, err)
        } else {
            log.Printf("%s Driver Terminated", browserName)
        }
    }()

    caps := selenium.Capabilities{"browserName": browserName}
    wd, err := selenium.NewRemote(caps, "http://localhost:"+driverPort+"/wd/hub")
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
    }

    log.Println("Completed all browser interactions successfully.")
    return nil
}

func startDriver(driverPath, driverPort string) (*exec.Cmd, error) {
    cmd := exec.Command(driverPath, "--port="+driverPort)
    cmd.Stdout = os.Stdout // Redirect driver stdout to os.Stdout
    cmd.Stderr = os.Stderr // Redirect driver stderr to os.Stderr
    err := cmd.Start()
    if err != nil {
        log.Printf("Failed to start driver: %v", err)
        return nil, err
    }
    log.Printf("Driver started successfully on port %s", driverPort)
    return cmd, nil
}