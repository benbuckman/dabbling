#!/usr/bin/env bash

# Deploy python code to Lambda endpoint.
# See: http://docs.aws.amazon.com/lambda/latest/dg/lambda-python-how-to-create-deployment-package.html

set -exo pipefail

main() {
    test -n "$LAMBDA_PROFILE" || { echo "Missing LAMBDA_PROFILE"; exit 1; }
    test -n "$LAMBDA_FUNCTION_NAME" || { echo "Missing LAMBDA_FUNCTION_NAME"; exit 1; }

    local CUR_DIR=$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)
    local SRC_DIR="${CUR_DIR}/lambda"
    local BUILD_DIR="${CUR_DIR}/_build"
    local ZIP_PATH="${BUILD_DIR}/slack_pybot.zip"

    echo "Zipping from $SRC_DIR to ZIP_PATH ..."

    rm -rf "${BUILD_DIR}" > /dev/null
    mkdir -p "${BUILD_DIR}" > /dev/null

    cd "$SRC_DIR"
    zip -r "${ZIP_PATH}" .

    # install dependencies
    #pip install ?? -t "$SRC_DIR"

    # Docs: http://docs.aws.amazon.com/cli/latest/reference/lambda/update-function-code.html
    aws lambda --profile "$LAMBDA_PROFILE" update-function-code \
        --function-name $LAMBDA_FUNCTION_NAME \
        --zip-file "fileb://$ZIP_PATH" \
        --no-publish
}

main "$@"