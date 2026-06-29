#!/bin/sh
set -e
cd "$(dirname "$0")"
source fly-secrets.env
flyctl deploy --remote-only
