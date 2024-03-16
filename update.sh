#!/bin/bash

# Navigate to your Go project directory
# Replace /path/to/your/project with the actual path to your Go project
cd /path/to/your/project || exit

# Step 1: Fetch the latest Git updates
echo "Fetching the latest updates from Git repository..."
git pull
if [ $? -ne 0 ]; then
    echo "Failed to fetch updates from Git. Please check your Git configuration."
    exit 1
fi

# Step 2: Run the Go program
# Replace main.go with the path to your Go program's entry point if it's different
echo "Running Go program..."
go run main.go
if [ $? -ne 0 ]; then
    echo "Failed to run Go program. Please check for compilation errors."
    exit 1
fi
