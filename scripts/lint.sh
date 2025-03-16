#!/bin/sh

if [ -z ${PACKAGE+x} ]; then echo "PACKAGE is not set"; exit 1; fi
if [ -z ${ROOT_DIR+x} ]; then echo "ROOT_DIR is not set"; exit 1; fi

echo "gofmt:"
OUT=$(gofmt -l $ROOT_DIR)
if [ -n "$OUT" ]; then echo "$OUT"; PROBLEM=1; fi

echo "errcheck:"
OUT=$(errcheck $PACKAGE)
if [ -n "$OUT" ]; then echo "$OUT"; PROBLEM=1; fi

echo "go vet:"
OUT=$(go vet -all=true $PACKAGE 2>&1 | grep --invert-match -E "(Checking file|\%p of wrong type|can't check non-constant format|could not import C)")
if [ -n "$OUT" ]; then echo "$OUT"; PROBLEM=1; fi

if [ -n "$PROBLEM" ]; then exit 1; fi
