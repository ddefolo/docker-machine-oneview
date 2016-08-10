Acceptance test cases
---------------------
This guide will describe steps we're using on running acceptance test cases.

Environment setup
------------------
Typical scenario, `make build` and `make test` should work as normal for the
overall project.

TODO:
We still need to integrate into the normal bats testing done for the overall project.

craete an environment, `drivers/oneview/.oneview.env`, script to export these values:

```bash
cat > "$(git rev-parse --show-toplevel)/drivers/oneview/.oneview.env" << ONEVIEW
ONEVIEW_APIVERSION=120

ONEVIEW_ILO_USER=docker
ONEVIEW_ILO_PASSWORD=password

ONEVIEW_ICSP_ENDPOINT=https://15.x.x.x
ONEVIEW_ICSP_USER=username
ONEVIEW_ICSP_PASSWORD=password
ONEVIEW_ICSP_DOMAIN=LOCAL

ONEVIEW_OV_ENDPOINT=https://15.x.x.x
ONEVIEW_OV_USER=username
ONEVIEW_OV_PASSWORD=password
ONEVIEW_OV_DOMAIN=LOCAL

ONEVIEW_SSLVERIFY=true
ONEVIEW_TEST_DATA=TEST_LAB_NAME

TESTCONFIG_PACKAGE_ROOT_PATH=github.com/docker/machine
TESTCONFIG_JSON_DATA_DIR=test/integration/data/oneview

ONEVIEW

```
NOTE: look in test/integration/data/oneview for test data that is pulled in from ONEVIEW_TEST_DATA

Setup container
---------------
1. setup gotest container, change TEST_CONTAINER_NAME if you want to run multiple test for other parts
```bash
TEST_CONTAINER_NAME=testov
docker run -it \
 --env-file "$(git rev-parse --show-toplevel)/drivers/oneview/.oneview.env" \
 -e ONEVIEW_TEST_ACCEPTANCE=true -e ICSP_TEST_ACCEPTANCE=true \
  -v $(git rev-parse --show-toplevel):/go/src/github.com/docker/machine \
  --name ${TEST_CONTAINER_NAME} docker-machine
# exit the started container
  docker restart ${TEST_CONTAINER_NAME}
```
2. setup alias:
```bash
   alias ${TEST_CONTAINER_NAME}='docker exec '${TEST_CONTAINER_NAME}' godep go test -test.timeout=60m -test.v=true --short'
```
3. to refresh env options, use
```bash
    docker rm -f ${TEST_CONTAINER_NAME}
    # ... repeat docker run commands in previous steps
```

Run test
--------
```bash
   cd "$(git rev-parse  --show-toplevel)"
   testov ./drivers/oneview/ov
```

Running Unit Test cases
-----------------------

```bash
sed -i '' 's/ONEVIEW_TEST_ACCEPTANCE=.*/ONEVIEW_TEST_ACCEPTANCE=false/g' "$(git rev-parse --show-toplevel)/drivers/oneview/.oneview.env"
sed -i '' 's/ICSP_TEST_ACCEPTANCE=.*/ICSP_TEST_ACCEPTANCE=false/g' "$(git rev-parse --show-toplevel)/drivers/oneview/.oneview.env"

docker rm -f ${TEST_CONTAINER_NAME}
docker run -it \
 --env-file "$(git rev-parse --show-toplevel)/drivers/oneview/.oneview.env" \
  -v $(git rev-parse --show-toplevel):/go/src/github.com/docker/machine \
  --name ${TEST_CONTAINER_NAME} docker-machine
  # exit the started container
    docker restart ${TEST_CONTAINER_NAME}
```
Proceed to run test
```bash
   cd "$(git rev-parse  --show-toplevel)"
   testov ./drivers/oneview/ov
```

Running debug log output
-------------------------
Add to the .oneview.env the DEBUG env.  This applies to all docker-machine code using log pacakge.

```bash
echo 'DEBUG=true' >> "$(git rev-parse --show-toplevel)/drivers/oneview/.oneview.env"
```
Follow section in setup container to refresh the env vars for the test container.

Run a single specific test
---------------------------
Sometimes it's usefull to run just a single test case.
```bash
testov ./drivers/oneview/ov -test.run=TestCreateProfileFromTemplate
```

Build one executable example
-----------------------------
Sometimes we need to be able to target building 1 executable
```bash
# checkout the script/build script for different target options, arg 1 and 2
script/build  -os="darwin" -arch="amd64"
alias build_machine='script/build  -os="darwin" -arch="amd64"'
```

Updating external dependencies
------------------------------
This project is no relying on glide to provide reliable & repeatable builds.
To learn more about glide, please visit : https://glide.sh/

Special thanks to Matt Farina for introducing it to us.

Start by installing glide:

```
curl https://glide.sh/get | sh
```

1. Add a dependency by editing the glide.yml in the root directory or run glide,
   glide get <package>#<version>, for example:

   ```
   glide get github.com/docker/machine#0.8.0
   ```

2. To update the existing packages, we use glide install.  Edit the glide.yaml and run:

   ```
   make glide
   ```

3. Run a build in a docker container.

   ```
   make test
   ```

4. Evaluate changes.
   At this point you might have changes to the dependent libraries that have
   to be incorporated into the build process.   Update any additional or
   no longer libraries by editing the file : [glide.yaml](glide.yaml).  
   This file contains all needed packages.
   Whenever adjusting libraries, make sure to re-do steps 1-3 iteratively.

5. Ok, it all test and passes, so it's time to commit your changes.

  ```
  git add --all
  ```
  Use `git status` to review additions, removals, and changes.
  Use `git commit -s -m "library update version X.X"` to commit your changes.

Creating a new Release
----------------------

New releases should be done by maintainers for the project. Steps on performing
these steps can be found in the [release](RELEASE.md) instructions.
