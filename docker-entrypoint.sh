#!/bin/sh

if [[ -z "${TRAVIS_JOB_ID}" ]]; then
    make migrate-latest
fi

exec "$@"