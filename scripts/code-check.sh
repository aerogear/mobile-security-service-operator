#!/bin/bash

# Run some pre commit checks on the Go source code. Prevent the commit if any errors are found

# Check if the Go code is formatted
check_code_format (){
    {
       go fmt  $(go list ./... )
    } || {
        exitStatus=$?

        if [ $exitStatus ]; then
            printf "\nErrors found in your code, please use 'make fmt' to format your code."
            exit 1
        fi
    }
}

# Check all files for suspicious constructs
check_go_vet (){
    {
        go vet $(go list ./... )
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
check_go_vet
