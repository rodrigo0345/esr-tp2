#!/bin/bash

# Assign the first positional parameter to the variable 'name'
name=$1

# Ensure the 'name' variable is not empty
if [ -z "$name" ]; then
  echo "Usage: $0 <node_name>"
  exit 1
fi

# Build the project
make

# Run the main executable with the provided parameters
# this is likely crash a lot if you have race conditions
./bin/main presence -v 10.0.13.10:4242 -n "$name"

