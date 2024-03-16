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

package main

import (
    "fmt"
    "log"
    "net"
    "os"
    "os/exec"
    "time"
)

func startDriver(driverPath, driverPort string) (*exec.Cmd, error) {
    // Check if the port is already in use
    log.Printf("Checking if port %s is available...", driverPort)
    if portInUse(driverPort) {
        return nil, fmt.Errorf("port %s is already in use", driverPort)
    }

    log.Printf("Starting driver with command: %s --port=%s", driverPath, driverPort)
    cmd := exec.Command(driverPath, "--port="+driverPort)
    cmd.Stdout = os.Stdout // Redirect stdout of the driver to os.Stdout of the program
    cmd.Stderr = os.Stderr // Redirect stderr of the driver to os.Stderr of the program
    err := cmd.Start()
    if err != nil {
        log.Printf("Failed to start the driver: %v", err)
        return nil, err
    }

    // Wait a moment to ensure driver has started
    time.Sleep(2 * time.Second)

    return cmd, nil
}

// portInUse tries to make a connection to the given port and returns true if it's already in use.
func portInUse(port string) bool {
    conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", port), time.Second)
    if err != nil {
        return false // Port is not in use
    }
    conn.Close() // Close the connection if open
    return true
}
