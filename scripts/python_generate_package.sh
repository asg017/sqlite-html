#!/bin/bash

set -eo pipefail

export PACKAGE_NAME_BASE="sqlite-html"
export EXTENSION_NAME="html0"
export VERSION=$(cat VERSION)

export LOADABLE_PATH=$1
export PYTHON_LOADABLE=$2
export OUTPUT_WHEELS=$3
export RENAME_WHEELS_ARGS=$4

cp $LOADABLE_PATH $PYTHON_LOADABLE
rm $OUTPUT_WHEELS/sqlite_html* || true
pip3 wheel python/sqlite_html/ -w $OUTPUT_WHEELS
python3 scripts/rename-wheels.py $OUTPUT_WHEELS $RENAME_WHEELS_ARGS
echo "✅ generated python wheel"

envsubst < python/version.py.tmpl > python/sqlite_html/sqlite_html/version.py
echo "✅ generated python/sqlite_html/sqlite_html/version.py"

envsubst < python/version.py.tmpl > python/datasette_sqlite_html/datasette_sqlite_html/version.py
echo "✅ generated python/datasette_sqlite_html/datasette_sqlite_html/version.py"