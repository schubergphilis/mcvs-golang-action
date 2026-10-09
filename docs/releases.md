# Releases

In some cases, you may want the executable binary to be built and released
automatically. This action will build the binary which could then be used
as a release asset.

Create a `.github/workflows/golang-releases.yml` file with the following
content:

```yml
---
name: golang-releases
"on": push
permissions:
  contents: write
  packages: read
jobs:
  mcvs-golang-action:
    strategy:
      matrix:
        args:
          - release-application-name: mcvs-image-downloader
            release-architecture: amd64
            release-dir: cmd/mcvs-image-downloader
            release-type: binary
          - release-application-name: mcvs-image-downloader
            release-architecture: arm64
            release-dir: cmd/mcvs-image-downloader
            release-os: darwin
            release-type: binary
    runs-on: ubuntu-24.04
    steps:
      - uses: actions/checkout@v4.2.2
      - uses: schubergphilis/mcvs-golang-action@v3
        with:
          release-application-name: ${{ matrix.args.release-application-name }}
          release-architecture: ${{ matrix.args.release-architecture }}
          release-build-tags: ${{ matrix.args.release-build-tags }}
          release-dir: ${{ matrix.args.release-dir }}
          release-os: ${{ matrix.args.release-os }}
          release-type: ${{ matrix.args.release-type }}
          token: ${{ secrets.GITHUB_TOKEN }}
```

## Downloading released assets from another private repository

You will need a personal access token (PAT) with the `repo` scope. To download
releases from a private repository. You can simply use the gh command or curl
to download the release assets. Please read the
[GitHub documentation](https://docs.github.com/en/rest/releases/assets)
for more information.
