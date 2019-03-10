#!/bin/sh

# Run some pre commit checks on the Go source code. Prevent the commit if any errors are found

# Check if the Go code is formatted
check_code_format (){
    gofmt -l . | grep -v vendor/ >/dev/null 2>&1
    if [ $? -ne 1 ]; then
       printf "\nErrors found in your code, please use 'go fmt' to format your code."
       exit 1
    else
       exit 0
    fi

}

# Check all files for errors
check_code_errors (){
    {
        errcheck -ignoretests $(go list ./... | grep -v /vendor/)
    } || {
        exitStatus=$?

        if [ $exitStatus ]; then
            printf "\nErrors found in your code, please fix them and try again."
            exit 1
        fi
    }
}

# Check all files for suspicious constructs
check_go_vet (){
    {
        go vet $(go list ./... | grep -v /vendor/)
    } || {
        exitStatus=$?

        if [ $exitStatus ]; then
            printf "\nIssues found in your code, please fix them and try again."
            exit 1
        fi
    }
}

# Calling the function
check_code_format
check_code_errors
check_go_vet
