#!/bin/bash

set -e

# This command will run pip install with the given requirements.txt file inside
# a docker container compatible with our function runtimes, and pull the installed
# dependencies locally to your package directory. As these dependencies have been
# installed on top of alpine Linux with our compatible system libraries, you will
# be able to upload your source code and deploy your function properly.
PYTHON_VERSION=3.10
docker run --rm -v $(pwd):/home/app/function --workdir /home/app/function rg.fr-par.scw.cloud/scwfunctionsruntimes-public/python-dep:$PYTHON_VERSION pip install -r ./requirements.txt --target ./package
