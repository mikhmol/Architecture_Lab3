#!/bin/bash

curl -X POST http://localhost:17000 -d "update"
curl -X POST http://localhost:17000 -d "reset"