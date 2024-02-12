#!/bin/bash

if [[ -z "$1" ]]; then
  echo "To few arguments"
  exit 1
fi

if [[ -v $1 ]]; then
  printenv $1
else
  echo "Variable not exists"
fi