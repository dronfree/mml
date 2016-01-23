#!/bin/bash
cd postfix && docker build -f Dockerfile.$1 -t postfix .
docker run -d -p 25:25 -t postfix