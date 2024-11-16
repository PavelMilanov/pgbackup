name: build pgbackup images

on:
  push:
      paths-ignore:
        # - '.github/**'
        - 'README.md'
        - 'Dockerfile'
        - docker-compose.yaml
      branches:
      - dev

jobs:
  push_container:
    runs-on: ubuntu-latest
    services:
      docker:
        image: docker:dind
        options: --privileged --shm-size=2g
        volumes:
          - /var/run/docker.sock:/var/run/docker.sock:ro
    container:
      image: ubuntu:latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Install Docker
        run: |
          apt-get update
          apt-get install -y docker.io          

      - name: Test Docker
        run: |
          docker version
          docker info          