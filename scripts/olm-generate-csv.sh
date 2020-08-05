#!/usr/bin/env sh

# Script to generate the CSV manifest for the operator.

set -eux

DIR=$(cd "$(dirname "$0")"; pwd -P)
. ${DIR}/env.sh

ROOT_DIR="$(dirname $DIR)"

last_bundle_version() {
    # get directories sorted by last modified first
    list=($(ls -dt $ROOT_DIR/$XROOTD_OPERATOR_BUNDLE_MANIFEST_DIR/*/))
    # return the first directory (latest modified)
    latest_dir=${list[0]}
    # print the directory name (i.e. version)
    echo -n "$(basename $latest_dir)"
}

# Generate CSV 
echo "Generating CSV for version ${XROOTD_OPERATOR_VERSION}"
operator-sdk generate csv \
    --operator-name ${XROOTD_OPERATOR_NAME} \
    --csv-version ${XROOTD_OPERATOR_VERSION} \
    --make-manifests=false \
    --from-version "$(last_bundle_version)" \
    --update-crds

# Set the operator image version
sed -i "s|REPLACE_IMAGE|$XROOTD_OPERATOR_FULL_IMAGE|g" \
    "$ROOT_DIR/$XROOTD_OPERATOR_BUNDLE_MANIFEST_VERSION_DIR/$XROOTD_OPERATOR_NAME.v$XROOTD_OPERATOR_VERSION.clusterserviceversion.yaml"
