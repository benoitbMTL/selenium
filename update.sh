#!/bin/bash

echo "Fetching the latest updates from Git repository..."
git pull
if [ $? -ne 0 ]; then
    echo "Failed to fetch updates from Git. Please check your Git configuration."
    exit 1
fi

# Step 2: Run the Go program
echo "Running Go program..."
go run .
if [ $? -ne 0 ]; then
    echo "Failed to run Go program. Please check for compilation errors."
    exit 1
fi


# git status
# git add .
# git commit -m "message"
# git push origin main