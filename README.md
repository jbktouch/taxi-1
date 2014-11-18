# Taxi

## What
Taxi ships your local code to a remote docker server for building containers and
running tests against them.

## Where
Taxi is used by me in TravisCI, but would probably work anywhere which supports
separate install / test / cleanup phases.

## Why
On one project with extensive dependencies, the build time was shortened from 45
minutes to 36 seconds.

## How

Integrating Taxi into your TravisCI workflow is simple if you know all the right
tricks. You probably don't. There is a tool at
`contrib/certificate_authority.sh` which will generate your certificate
authority, docker server keys, and client keys. It will also generate a userdata
script for a CoreOS node to run as your remote.

Here is what your .travis.yml might look like:

```yaml
env:
  global:
    - TAXI_VERSION=0.1.0-alpha
    - TAXI_URL=https://github.com/grahamc/taxi/releases/download/v0.1.0-alpha/taxi-0.1.0-alpha.0-prerelease
cache:
  directories:
    - .taxi-cache/

install:
 - if [ ! -f .taxi-cache/taxi-$TAXI_VERSION ]; then wget $TAXI_URL -O .taxi-cache/taxi-$TAXI_VERSION; fi; cp .taxi-cache/taxi-$TAXI_VERSION ./taxi; chmod +x ./taxi
 - ./taxi install

script:
 - ./taxi test "pip install flake8; flake8 /code"
 - ./taxi test "python -m compileall /code"

after_script:
 - ./taxi cleanup
```

