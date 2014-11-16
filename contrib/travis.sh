
cat <<EOM
env:
  global:
    - DOCKER_HOST="tcp://$HOST:2376"
    - TAXI_VERSION=0.1.0-alpha
    - TAXI_URL=https://github.com/grahamc/taxi/releases/download/v0.1.0-alpha/taxi-0.1.0-alpha.0-prerelease
cache:
  directories:
    - .taxi-cache/

install:
 - if [ ! -f .taxi-cache/taxi-\$TAXI_VERSION ]; then wget \$TAXI_URL -O .taxi-cache/taxi-\$TAXI_VERSION; fi; cp .taxi-cache/taxi-\$TAXI_VERSION ./taxi; chmod +x ./taxi
 - ./taxi install

script:
 - ./taxi test "exit 1"

after_script:
 - ./taxi cleanup

EOM

