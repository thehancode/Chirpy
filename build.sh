#!/bin/bash

build_and_run() {
  echo "Building and running..."
  go build -o out && chmod +x ./out && ./out
}

build_and_run

fsnotify . | while read event; do
  if [[ $event == *"WRITE"* ]] && [[ $event == *".go"* ]]; then
    echo "File changed. Rebuilding..."
    build_and_run
  fi
done
