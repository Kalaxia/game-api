#!/bin/sh
npm install
npm run prod

exec "$@"