#!/usr/bin/env bash

get_minimum_version() {
    # to extract the version from 'version.go'
    MIN_VERSION=$(grep -w "defaultMinimumVersion =" vendor/knative.dev/pkg/version/version.go | sed -E 's/.*"([^"]+)".*/\1/')
    echo "${MIN_VERSION}"
}


check_kubernetes_version() {
    local MINIMUM_VERSION=$(get_minimum_version)
    local CURRENT_VERSION=$(kubectl version --output=json 2>/dev/null | jq -r '.serverVersion.gitVersion')

    if [[ "$CURRENT_VERSION" < "$MINIMUM_VERSION" ]]; then
        echo "Your Kubernetes version ($CURRENT_VERSION) is less than the required minimum version ($MINIMUM_VERSION)."
        exit 1
    else
        echo "Kubernetes version check passed."
    fi
}
