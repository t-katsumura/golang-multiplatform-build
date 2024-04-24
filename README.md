# Golang Multiplatform Build with GitHub Actions

Build test of golang application on multiple platforms.

## Using github hosted runner

Example workflow at [.github/workflows/runner.yaml](.github/workflows/runner.yaml).

```yaml
name: Runner

on: [push, workflow_dispatch]

jobs:
  unit:
    name: Runner
    runs-on: ${{ matrix.os }}-latest
    timeout-minutes: 15
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

on: [push, workflow_dispatch]

jobs:
  unit:
    name: QEMU
    runs-on: ubuntu-latest
    timeout-minutes: 15
    strategy:
      matrix:
        arch:
          - linux/amd64
          - linux/arm64
          # - linux/riscv64  # Not supported by go official image.
          - linux/ppc64le
          - linux/s390x
          # - linux/386      # Runtime error occurs.
          # - linux/mips64le # Runtime error occurs.
          # - linux/mips64   # Not supported by go official image.
          - linux/arm/v7
          # - linux/arm/v6   # Not supported by go official image.
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

## Using SLSA releaser

Example workflow at [.github/workflows/slsa.yaml](.github/workflows/slsa.yaml).

This workflow uses [slsa-github-generator](https://github.com/slsa-framework/slsa-github-generator/tree/main/internal/builders/go).

Versions are fixed here. Tag names should be used there.

```yaml
name: SLSA

on: push

permissions:
  id-token: write
  contents: write
  actions: read

jobs:
  build:
    uses: slsa-framework/slsa-github-generator/.github/workflows/builder_go_slsa3.yml@v1.10.0
    strategy:
      matrix:
        os:
          - linux
          - windows
          - darwin
          - freebsd
          - openbsd
          - netbsd
        arch:
          - 386
          - amd64
          - arm
          - arm64
        include:
          - os: linux
            arch: ppc64
          - os: linux
            arch: ppc64le
          - os: linux
            arch: riscv64
          - os: linux
            arch: s390x
        exclude:
          - os: darwin
            arch: 386
          - os: darwin
            arch: arm
    with:
      go-version-file: ./go.mod
      config-file: .github/targets/${{matrix.os}}-${{matrix.arch}}.yaml
      evaluated-envs: "VERSION:v0.0.0"
      draft-release: true
      upload-assets: true
      upload-tag-name: v0.0.0
```

## Using ko (with SLSA provenance)

Example workflow at [.github/workflows/container.yaml](.github/workflows/container.yaml).

```yaml
name: Container

on: push

permissions:
  contents: read
  packages: write
  actions: write
  id-token: write

env:
  IMAGE_REGISTRY: ghcr.io
  IMAGE_NAME: hello

jobs:
  build:
    name: build
    runs-on: ubuntu-latest
    timeout-minutes: 15
    outputs:
      image: ${{ steps.build.outputs.image }}
      digest: ${{ steps.build.outputs.digest }}
    steps:
      - run: git config --global core.autocrlf false
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          check-latest: true
          go-version-file: ./go.mod
      - uses: ko-build/setup-ko@v0.6
      - id: build # id cannot be omitted.
        run: |
          export VERSION=v0.0.0
          ko build -v -B --sbom=none --tags=latest --tags=${VERSION} --image-refs .digest ./
          image=$(sed '1q;d' .digest | cut -d'@' -f1 | cut -d':' -f1)
          digest=$(sed '1q;d' .digest | cut -d'@' -f2)
          echo "image=$image" >> "$GITHUB_OUTPUT"
          echo "digest=$digest" >> "$GITHUB_OUTPUT"
      - run: cat .digest

  provenance:
    needs: build
    uses: slsa-framework/slsa-github-generator/.github/workflows/generator_container_slsa3.yml@v1.10.0
    with:
      image: ${{ needs.build.outputs.image }}
      digest: ${{ needs.build.outputs.digest }}
      registry-username: ${{ github.actor }}
      compile-generator: true
    secrets:
      registry-password: ${{ secrets.GITHUB_TOKEN }}
```
