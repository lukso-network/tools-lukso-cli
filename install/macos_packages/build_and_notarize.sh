#!/bin/bash

set -euo pipefail

# Environmental variables that need to be set. These are sane defaults
# KEYCHAIN_ENTRY=AC_PASSWORD   # Or whatever you picked for <AUTH_ITEM_NAME>
# RUNNER_TEMP=~/Library/Keychains/
# KEYCHAIN_FILENAME=login.keychain-db
#
# Environmental variables that can be set if needed. Else they will default to
# values selected for the CI
#
# The identifier for the application signing key. Can be a name or a fingerprint
# APPLICATION_KEY_IDENT=E03290CABE09E9E42341C8FC82608E91241FAD4A
# The identifier for the installer signing key. Can be a name or a fingerprint
# INSTALLATION_KEY_IDENT=E525359D0B5AE97B7B6F5BB465FEC872C117D681

usage() {
  echo "Builds, signs, and notarizes lukso."
  echo "Read readme.md in the script directory to see the assumptions the script makes."
  echo "Usage: $0 <path-to-lukso-binary> <path-to-docs-directory> <arch> [pkgname]"
  echo "  Example: $0 target/release/lukso docs/ x64"
  echo "  Example: $0 target/debug/lukso docs/ arm64 lukso-1.2.1-arm64.pkg"
  echo ""
  echo "If no pkgname is provided, the package will be named lukso-<version>-<arch>.pkg"
}

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
source "$script_dir/common.sh"

if [[ -z ${KEYCHAIN_ENTRY+x} ]]; then
  error "Environmental variable KEYCHAIN_ENTRY must be set."
fi

if [[ -z ${RUNNER_TEMP+x} ]]; then
  error "Environmental variable RUNNER_TEMP must be set."
fi

if [[ -z ${KEYCHAIN_FILENAME+x} ]]; then
  error "Environmental variable KEYCHAIN_FILENAME must be set."
fi

keychain_path="$RUNNER_TEMP/$KEYCHAIN_FILENAME"
if [[ ! -f "$keychain_path" ]]; then
  error "Could not find keychain at $keychain_path"
fi

if [[ -z ${APPLICATION_KEY_IDENT+x} ]]; then
  APPLICATION_KEY_IDENT=82E5F57AD929C8D626C4410EFAD27E0BC2C0ED97
  echo "APPLICATION_KEY_IDENT not set. Using default value of $APPLICATION_KEY_IDENT"
fi

if [[ -z ${INSTALLATION_KEY_IDENT+x} ]]; then
  INSTALLATION_KEY_IDENT=02F4E6CCCCF6A26097D5E18F938531BC8237C7D4
  echo "INSTALLATION_KEY_IDENT not set. Using default value of $INSTALLATION_KEY_IDENT"
fi

if [[ -z ${3+x} ]]; then
  usage
  exit 1
fi

lukso_binary="$1"
lukso_docs_dir="$2"
arch="$3"
pkgname="${4:-}"

echo ">>>> Signing binary"
codesign --timestamp --keychain "$keychain_path" --sign "$APPLICATION_KEY_IDENT" --verbose -f -o runtime "$lukso_binary"

# Make ZIP file to notarize binary
if [ "$lukso_binary" != "lukso" ]; then
  cp "$lukso_binary" lukso
fi
zip lukso.zip lukso

echo ">>>> Submitting binary for notarization"
xcrun notarytool submit lukso.zip --keychain-profile "$KEYCHAIN_ENTRY" --wait

# Don't think this is actually necessary, but not costly so why not
rm lukso
unzip lukso.zip

# Create the component package
echo ">>>> Building Component Package"
bash "$script_dir/build_component_package.sh" "lukso" "$lukso_docs_dir"

# Create the distribution package
echo ">>>> Building Distribution Package"
resources_path="$script_dir/pkg_resources"
bash "$script_dir/build_distribution_package.sh" "lukso-component.pkg" "$resources_path" "$arch"

# Codesign the package installer
productsign --timestamp --sign "$INSTALLATION_KEY_IDENT" lukso-unsigned.pkg lukso.pkg

# Notarize the package installer
echo ">>>> Submitting .pkg for notarization"
xcrun notarytool submit lukso.pkg --keychain-profile "$KEYCHAIN_ENTRY" --wait

# Staple things
echo ">>>> Running final steps"
xcrun stapler staple lukso.pkg

# Rename to expected name
if [ "$pkgname" = "" ]; then
  version="$(lukso_version "$lukso_binary")"
  pkgname="lukso-$version-$arch.pkg"
fi

echo ">>>> Placing final output at $pkgname"
mv lukso.pkg "$pkgname"
