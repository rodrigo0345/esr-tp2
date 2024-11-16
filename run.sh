#!/bin/bash

$name = @1

# This script works for the presence nodes
make
./bin/main presence -v 10.0.13.10:4242 -n $name
