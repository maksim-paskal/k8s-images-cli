# List images in your kubernetes cluster

Simple tool to view images in your kubernetes cluster

## Install

MacOS

```bash
brew install maksim-paskal/tap/k8s-images-cli
```

## Usage

```bash
k8s-images-cli -n some_namespace
```

## filter by image

```bash
k8s-images-cli -image some_image
```