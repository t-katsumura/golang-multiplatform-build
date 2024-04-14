# Golang Multiplatform Build with GitHub Actions

Build test of golang application on multiple platforms.

## Using github hosted runner

Example workflow at [.github/workflows/runner.yaml](.github/workflows/runner.yaml).

```yaml
name: Runner

on: push

jobs:
  unit:
    name: Runner
    runs-on: ${{ matrix.os }}-latest
    timeout-minutes: 30
    strategy:
      matrix:
        os:
          - ubuntu
          - windows
          - macos
    steps:
      - run: git config --global core.autocrlf false
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          check-latest: true
          go-version-file: ./go.mod
      - run: go run ./...
```

## Using QEMU emulator

Example workflow at [.github/workflows/qemu.yaml](.github/workflows/qemu.yaml).

Available architectures are listed below.

- <https://github.com/docker/setup-qemu-action>
- <https://hub.docker.com/_/golang/tags>

```yaml
name: QEMU

on: push

jobs:
  unit:
    name: QEMU
    runs-on: ubuntu-latest
    timeout-minutes: 30
    strategy:
      matrix:
        arch:
          - linux/amd64
          - linux/arm64
          # - linux/riscv64 # Not supported by go official image.
          - linux/ppc64le
          - linux/s390x
          # - linux/386 # Runtime error occurs.
          # - linux/mips64le # Runtime error occurs.
          # - linux/mips64 # Not supported by go official image.
          - linux/arm/v7
          # - linux/arm/v6 # Not supported by go official image.
    steps:
      - run: git config --global core.autocrlf false
      - uses: actions/checkout@v4
      - uses: docker/setup-qemu-action@v3
        with:
          platforms: ${{ matrix.arch }}
      - name: Docker run
        run: |
          docker run \
            --rm \
            -v $(pwd):/${{ github.workspace }} \
            --platform ${{ matrix.arch }} \
            golang:1.22.2 \
            bash -c "go run /${{ github.workspace }}/main.go"
```
