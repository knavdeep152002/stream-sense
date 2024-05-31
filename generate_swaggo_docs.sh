#!/bin/bash

swag fmt && \
swag init --pd -d cmd/server,internal/streamsense,internal/fs,internal/openai,internal/auth
