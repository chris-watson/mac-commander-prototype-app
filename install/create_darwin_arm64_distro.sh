#!/bin/bash

set -e  

# for now there is nothing platform-specific
# in our package script, so just call the
# amd64 version
./create_darwin_amd64_distro.sh $1
