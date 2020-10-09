
[![Docker Stars](https://img.shields.io/docker/stars/avsr/appver-resource.svg?style=plastic&logo=docker&label=stars)](https://registry.hub.docker.com/v2/repositories/avsr/appver-resource/stars/count/)
[![Docker pulls](https://img.shields.io/docker/pulls/avsr/appver-resource.svg?style=plastic&logo=docker&label=pulls)](https://registry.hub.docker.com/v2/repositories/avsr/appver-resource)
[![Docker build status](https://img.shields.io/docker/cloud/build/avsr/appver-resource.svg?logo=docker&style=plastic&label=build)](https://github.com/aedavelli/appver-resource)
[![Docker Automated build](https://img.shields.io/docker/cloud/automated/avsr/appver-resource.svg?logo=docker&label=build)](https://github.com/aedavelli/appver-resource)
[![Docker Size](https://img.shields.io/docker/image-size/avsr/appver-resource/latest?label=size&logo=docker&style=plastic)](https://hub.docker.com/r/avsr/appver-resource/)

[![dockeri.co](http://dockeri.co/image/avsr/appver-resource)](https://hub.docker.com/r/avsr/appver-resource/)

# Concourse CI Web Application Version Resource

Implements a resource that passes to a task the version and different meta data. Can be used to trigger a job based on specified version field change.

## Source Configuration

#### Parameters

* `url`: *Required.* URL which produces version info in one of the JSON/XML/TEXT formats.

* `version_field`: *Optional.* in general. *Required* if response is JSON/XML. The level is separated by *::* 2 colons.
```json
{
  "version_info" : {
    "git" : {
      "version" : "2a3b4e",
      "branch" : "release",
      "tag" : "v1.1.1",
      "url" : "https://github.com/aedavelli/appver-resource"
    }
  }
}
 ```
 *version_info::git::version* can be used as version_field. All concrete fields under version_info.git will be produced as metadata

* `accept`: *Optional.* in general. If not present defaults to *application/json*

* `username`: *Optional.* in general. Used for basic auth if the URL is protected.

* `password`: *Optional.* in general. Used for basic auth if the URL is protected.

#### Example

``` YAML
resource_types:
  - name: appver
    type: docker-image
    source:
      repository: avsr/appver-resource

resources:
  - name: appver
    type: appver
    source:
      url: http://192.168.1.13:8282/version
      version_field: "version_info::git::version"
      accept: "application/json"
      username: ((webapp.username))
      password: ((webapp.password))
      insecure: true


jobs:
- name: some-job
  plan:
  - get: appver
    trigger: true
  - task: version-changed-task
    config:
      platform: linux
      image_resource:
        type: docker-image
        source: {repository: busybox}
      # appver directory contains key(filename), value(data inside file)
      inputs:
        - name: appver
      run:
        # do your work here
```

## Behavior

### `check`: Check for latest version

When no version/latest version  specified, latest version will be returned. If old version specified in request, both old and latest versions will be returned.

### `in`: Write the metadata to the destination dir

For above mentioned json response below files and content in appver directory

File|Content|
--- | ---
appver/version | 2a3b4e
appver/branch | release
appver/tag | v1.1.1
appver/url | https://github.com/aedavelli/appver-resource

#### Parameters

*None.*

### `out`: Unused

Unused

#### Parameters

*None.*

## Use in a task

```sh

version=`cat "${ROOT_FOLDER}/appver/version"`
echo "Using web APP with version $version"

```

## Development

### Prerequisites

* golang is *required* - version 1.14.x is tested; earlier versions may also
  work.
* docker is *required* - version 17.06.x is tested; earlier versions may also
  work.

### Running the tests

```bash
go test -v github.com/aedavelli/appver-resource/... -ginkgo.v
```
The tests have been embedded with the `Dockerfile`; ensuring that the testing
environment is consistent across any `docker` enabled platform. When the docker
image builds, the test are run inside the docker container, on failure they
will stop the build.

Run the tests with the following command:

```sh
docker build -t appver-resource .
```

### Contributing

Please make all pull requests to the `master` branch and ensure tests pass
locally.
