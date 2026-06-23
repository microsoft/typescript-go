# luchta-tsc-worker

The purpose of this fork is to build the luchta-tsc-worker.

Use the scripts in `scripts/` to install the worker from last published release or to publish a release.

To build and install from source:

    npm ci && npx hereby worker:build && cp ./built/worker/linux-amd64/luchta-tsc-worker ~/.local/bin/

Refer back to the original repo for the main documentation

In addition to adding the luchta worker binary this includes patches to make it work with Yarn PnP.

