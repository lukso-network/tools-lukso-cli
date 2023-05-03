#!/usr/bin/env sh

set -eu
printf '\n'

BOLD="$(tput bold 2>/dev/null || printf '')"
GREY="$(tput setaf 0 2>/dev/null || printf '')"
UNDERLINE="$(tput smul 2>/dev/null || printf '')"
RED="$(tput setaf 1 2>/dev/null || printf '')"
GREEN="$(tput setaf 2 2>/dev/null || printf '')"
YELLOW="$(tput setaf 3 2>/dev/null || printf '')"
BLUE="$(tput setaf 4 2>/dev/null || printf '')"
MAGENTA="$(tput setaf 5 2>/dev/null || printf '')"
NO_COLOR="$(tput sgr0 2>/dev/null || printf '')"

SUPPORTED_TARGETS="x86_64-unknown-linux-gnu x86_64-unknown-linux-musl \
                  i686-unknown-linux-musl aarch64-unknown-linux-musl \
                  arm-unknown-linux-musleabihf x86_64-apple-darwin \
                  aarch64-apple-darwin x86_64-unknown-freebsd"


info() {
  printf '%s\n' "${BOLD}${GREY}>${NO_COLOR} $*"
}

warn() {
  printf '%s\n' "${YELLOW}! $*${NO_COLOR}"
}

error() {
  printf '%s\n' "${RED}x $*${NO_COLOR}" >&2
}

completed() {
  printf '%s\n' "${GREEN}✓${NO_COLOR} $*"
}

has() {
  command -v "$1" 1>/dev/null 2>&1
}

# Make sure user is not using zsh or non-POSIX-mode bash, which can cause issues
verify_shell_is_posix_or_exit() {
  if [ -n "${ZSH_VERSION+x}" ]; then
    error "Running installation script with \`zsh\` is known to cause errors."
    error "Please use \`sh\` instead."
    exit 1
  elif [ -n "${BASH_VERSION+x}" ] && [ -z "${POSIXLY_CORRECT+x}" ]; then
    error "Running installation script with non-POSIX \`bash\` may cause errors."
    error "Please use \`sh\` instead."
    exit 1
  else
    true  # No-op: no issues detected
  fi
}

# Gets path to a temporary file, even if
get_tmpfile() {
  suffix="$1"
  if has mktemp; then
    printf "%s.%s" "$(mktemp)" "${suffix}"
  else
    # No really good options here--let's pick a default + hope
    printf "/tmp/lukso.%s" "${suffix}"
  fi
}

# Test if a location is writeable by trying to write to it. Windows does not let
# you test writeability other than by writing: https://stackoverflow.com/q/1999988
test_writeable() {
  path="${1:-}/test.txt"
  if touch "${path}" 2>/dev/null; then
    rm "${path}"
    return 0
  else
    return 1
  fi
}

download() {
  file="$1"
  url="$2"

  if has curl; then
    cmd="curl --fail --silent --location --output $file $url"
  elif has wget; then
    cmd="wget --quiet --output-document=$file $url"
  elif has fetch; then
    cmd="fetch --quiet --output=$file $url"
  else
    error "No HTTP download program (curl, wget, fetch) found, exiting…"
    return 1
  fi

  $cmd && return 0 || rc=$?

  error "Command failed (exit code $rc): ${BLUE}${cmd}${NO_COLOR}"
  printf "\n" >&2
  info "This is likely due to lukso not yet supporting your configuration."
  info "If you would like to see a build for your configuration,"
  info "please create an issue requesting a build for ${MAGENTA}${TARGET}${NO_COLOR}:"
  info "${BOLD}${UNDERLINE}https://github.com/lukso-network/tools-lukso-cli/issues/new/${NO_COLOR}"
  return $rc
}

unpack() {
  archive=$1
  bin_dir=$2
  sudo=${3-}

  case "$archive" in
    *.tar.gz)
      flags=$(test -n "${VERBOSE-}" && echo "-xzvof" || echo "-xzof")
      ${sudo} tar "${flags}" "${archive}" -C "${bin_dir}"
      return 0
      ;;
    *.zip)
      flags=$(test -z "${VERBOSE-}" && echo "-qqo" || echo "-o")
      UNZIP="${flags}" ${sudo} unzip "${archive}" -d "${bin_dir}"
      return 0
      ;;
  esac

  error "Unknown package extension."
  printf "\n"
  info "This almost certainly results from a bug in this script--please file a"
  info "bug report at https://github.com/lukso-network/tools-lukso-cli/issues"
  return 1
}

usage() {
  printf "%s\n" \
    "install.sh [option]" \
    "" \
    "Fetch and install the latest version of lukso, if lukso is already" \
    "installed it will be updated to the latest version."

  printf "\n%s\n" "Options"
  printf "\t%s\n\t\t%s\n\n" \
    "-V, --verbose" "Enable verbose output for the installer" \
    "-f, -y, --force, --yes" "Skip the confirmation prompt during installation" \
    "-p, --platform" "Override the platform identified by the installer [default: ${PLATFORM}]" \
    "-b, --bin-dir" "Override the bin installation directory [default: ${BIN_DIR}]" \
    "-a, --arch" "Override the architecture identified by the installer [default: ${ARCH}]" \
    "-B, --base-url" "Override the base URL used for downloading releases [default: ${BASE_URL}]" \
    "-h, --help" "Display this help message"
}

elevate_priv() {
  if ! has sudo; then
    error 'Could not find the command "sudo", needed to get permissions for install.'
    info "Otherwise, please run this script as root, or install"
    info "sudo."
    exit 1
  fi
  if ! sudo -v; then
    error "Superuser not granted, aborting installation"
    exit 1
  fi
}

install() {
  ext="$1"

  if test_writeable "${BIN_DIR}"; then
    sudo=""
    msg="Installing lukso, please wait…"
  else
    warn "Escalated permissions are required to install to ${BIN_DIR}"
    elevate_priv
    sudo="sudo"
    msg="Installing lukso as root, please wait…"
  fi
  info "$msg"

  archive=$(get_tmpfile "$ext")

  # download to the temp file
  download "${archive}" "${URL}"

  # unpack the temp file to the bin dir, using sudo if required
  unpack "${archive}" "${BIN_DIR}" "${sudo}"
}

# Currently supporting:
#   - win (Git Bash)
#   - darwin
#   - linux
#   - linux_musl (Alpine)
#   - freebsd
detect_platform() {
  platform="$(uname -s | tr '[:upper:]' '[:lower:]')"

  case "${platform}" in
    # use the statically compiled musl bins on linux to avoid linking issues.
    linux) platform="unknown-linux-musl" ;;
    darwin) platform="apple-darwin" ;;
    freebsd) platform="unknown-freebsd" ;;
  esac

  printf '%s' "${platform}"
}

# Currently supporting:
#   - x86_64
#   - i386
#   - arm
#   - arm64
detect_arch() {
  arch="$(uname -m | tr '[:upper:]' '[:lower:]')"

  case "${arch}" in
    amd64) arch="x86_64" ;;
    armv*) arch="arm" ;;
    arm64) arch="aarch64" ;;
  esac

  # `uname -m` in some cases mis-reports 32-bit OS as 64-bit, so double check
  if [ "${arch}" = "x86_64" ] && [ "$(getconf LONG_BIT)" -eq 32 ]; then
    arch=i686
  elif [ "${arch}" = "aarch64" ] && [ "$(getconf LONG_BIT)" -eq 32 ]; then
    arch=arm
  fi

  printf '%s' "${arch}"
}

detect_target() {
  arch="$1"
  platform="$2"
  target="$arch-$platform"

  if [ "${target}" = "arm-unknown-linux-musl" ]; then
    target="${target}eabihf"
  fi

  printf '%s' "${target}"
}


confirm() {
  if [ -z "${FORCE-}" ]; then
    printf "%s " "${MAGENTA}?${NO_COLOR} $* ${BOLD}[y/N]${NO_COLOR}"
    set +e
    read -r yn </dev/tty
    rc=$?
    set -e
    if [ $rc -ne 0 ]; then
      error "Error reading from prompt (please re-run with the '--yes' option)"
      exit 1
    fi
    if [ "$yn" != "y" ] && [ "$yn" != "yes" ]; then
      error 'Aborting installation'
      exit 1
    fi
  fi
}

check_bin_dir() {
  bin_dir="${1%/}"

  if [ ! -d "$BIN_DIR" ]; then
    error "Installation location $BIN_DIR does not appear to be a directory"
    info "Make sure the location exists and is a directory, then try again."
    usage
    exit 1
  fi

  # https://stackoverflow.com/a/11655875
  good=$(
    IFS=:
    for path in $PATH; do
      if [ "${path%/}" = "${bin_dir}" ]; then
        printf 1
        break
      fi
    done
  )

  if [ "${good}" != "1" ]; then
    warn "Bin directory ${bin_dir} is not in your \$PATH"
  fi
}

print_install() {
  echo "To get started do the following:"
  echo "${BOLD}1.${NO_COLOR} 📁  Create the working directory using ${BOLD}'mkdir myLUKSOnode && cd ./myLUKSOnode${NO_COLOR}'"
  echo "${BOLD}2.${NO_COLOR} 🚀  Initialize the folder using ${BOLD}'lukso init${NO_COLOR}'"
  echo "${BOLD}3.${NO_COLOR} 🛠️   Install your desired clients using ${BOLD}'lukso install${NO_COLOR}'"
}

is_build_available() {
  arch="$1"
  platform="$2"
  target="$3"

  good=$(
    IFS=" "
    for t in $SUPPORTED_TARGETS; do
      if [ "${t}" = "${target}" ]; then
        printf 1
        break
      fi
    done
  )

  if [ "${good}" != "1" ]; then
    error "${arch} builds for ${platform} are not yet available for lukso"
    printf "\n" >&2
    info "If you would like to see a build for your configuration,"
    info "please create an issue requesting a build for ${MAGENTA}${target}${NO_COLOR}:"
    info "${BOLD}${UNDERLINE}https://github.com/lukso-network/tools-lukso-cli/issues/new/${NO_COLOR}"
    printf "\n"
    exit 1
  fi
}

# defaults
if [ -z "${PLATFORM-}" ]; then
  PLATFORM="$(detect_platform)"
fi

if [ -z "${BIN_DIR-}" ]; then
  BIN_DIR=/usr/local/bin
fi

if [ -z "${ARCH-}" ]; then
  ARCH="$(detect_arch)"
fi

if [ -z "${BASE_URL-}" ]; then
  BASE_URL="https://storage.googleapis.com/lks-lz-binaries-euw4"
fi

# Non-POSIX shells can break once executing code due to semantic differences
verify_shell_is_posix_or_exit

# parse argv variables
while [ "$#" -gt 0 ]; do
  case "$1" in
  -p | --platform)
    PLATFORM="$2"
    shift 2
    ;;
  -b | --bin-dir)
    BIN_DIR="$2"
    shift 2
    ;;
  -a | --arch)
    ARCH="$2"
    shift 2
    ;;
  -B | --base-url)
    BASE_URL="$2"
    shift 2
    ;;

  -V | --verbose)
    VERBOSE=1
    shift 1
    ;;
  -f | -y | --force | --yes)
    FORCE=1
    shift 1
    ;;
  -h | --help)
    usage
    exit
    ;;

  -p=* | --platform=*)
    PLATFORM="${1#*=}"
    shift 1
    ;;
  -b=* | --bin-dir=*)
    BIN_DIR="${1#*=}"
    shift 1
    ;;
  -a=* | --arch=*)
    ARCH="${1#*=}"
    shift 1
    ;;
  -B=* | --base-url=*)
    BASE_URL="${1#*=}"
    shift 1
    ;;
  -V=* | --verbose=*)
    VERBOSE="${1#*=}"
    shift 1
    ;;
  -f=* | -y=* | --force=* | --yes=*)
    FORCE="${1#*=}"
    shift 1
    ;;

  *)
    error "Unknown option: $1"
    usage
    exit 1
    ;;
  esac
done

TARGET="$(detect_target "${ARCH}" "${PLATFORM}")"

is_build_available "${ARCH}" "${PLATFORM}" "${TARGET}"

if WHERE=$(which lukso)
then
  WHERE=$(dirname "${WHERE}")
  PREVIOUS=""
  if ! PREVIOUS=$(lukso version 2> /dev/null)
  then
    if PREVIOUS=$(lukso -v 2> /dev/null)
    then
      PREVIOUS="${PREVIOUS} (deprecated)"
    else
      PREVIOUS=""    
    fi
  fi
  PREVIOUS=$(echo "${PREVIOUS}" | sed -n -e "s#.*\(v[0-9][0-9]*\..*\)\$#\1#p" -e "s#.*\(develop\)\$#\1#p")
  if [ -n "${PREVIOUS}" ]
  then
    if [ "${PREVIOUS}" = "__VERSION__" ]
    then
      printf "  %s\n" "${UNDERLINE}${GREEN}You currently have the current version installed${NO_COLOR}"
      info "${BOLD}Current bin${NO_COLOR}:      ${GREEN}${WHERE}${NO_COLOR}"
      info "${BOLD}Current version${NO_COLOR}:  ${GREEN}${PREVIOUS}${NO_COLOR}"
      info "${BOLD}Platform${NO_COLOR}:         ${GREEN}${PLATFORM}${NO_COLOR}"
      info "${BOLD}Arch${NO_COLOR}:             ${GREEN}${ARCH}${NO_COLOR}"
      exit
    fi
    printf "  %s\n" "${UNDERLINE}You currently have a previous version installed${NO_COLOR}"
    info "${BOLD}Current bin${NO_COLOR}:      ${GREEN}${WHERE}${NO_COLOR}"
    info "${BOLD}Current version${NO_COLOR}:  ${GREEN}${PREVIOUS}${NO_COLOR}"
    if [ "${BIN_DIR}" != "${WHERE}" ]
    then
      echo ""
      info "${RED}CAUTION: This installer will install the bin into ${GREEN}${BIN_DIR}${RED}"
      info "${RED}while your old executable is in ${GREEN}${WHERE}${RED}"
      info "${RED}make sure your ${GREEN}\$PATH${RED} includes the new one first${NO_COLOR}"
    fi
    echo ""
    info "${RED}Press Ctrl-C or reply ${BOLD}N${NO_COLOR}${RED} if you don't want to upgrade to version __VERSION__${NO_COLOR}"
    echo ""
  fi
fi
printf "  %s\n" "${UNDERLINE}Configuration${NO_COLOR}"
info "${BOLD}Repository${NO_COLOR}:       ${GREEN}https://github.com/lukso-network/tools-lukso-cli"
info "${BOLD}Bin directory${NO_COLOR}:    ${GREEN}${BIN_DIR}${NO_COLOR}"
info "${BOLD}Platform${NO_COLOR}:         ${GREEN}${PLATFORM}${NO_COLOR}"
info "${BOLD}Arch${NO_COLOR}:             ${GREEN}${ARCH}${NO_COLOR}"
info "${BOLD}Version${NO_COLOR}:          ${GREEN}__VERSION__${NO_COLOR}"

# non-empty VERBOSE enables verbose untarring
if [ -n "${VERBOSE-}" ]; then
  VERBOSE=v
  info "${BOLD}Verbose${NO_COLOR}: yes"
else
  VERBOSE=
fi

printf '\n'

EXT=tar.gz

URL="${BASE_URL}/artifacts/lukso-${TARGET}.${EXT}"
info "Tarball URL: ${UNDERLINE}${BLUE}${URL}${NO_COLOR}"
info ""
info "You need to agree to the LUKSO CLI terms before continuing:"
info "${BOLD}https://github.com/lukso-network/tools-lukso-cli/blob/main/LICENSE${NO_COLOR}"
confirm "Agree and install lukso ${GREEN}__VERSION__${NO_COLOR} to ${BOLD}${GREEN}${BIN_DIR}${NO_COLOR}?"
check_bin_dir "${BIN_DIR}"

# if the shell does not fit the default case change the config file
# and or the config cmd variable
info "Installing ${BOLD}LUKSO CLI${NO_COLOR}..."

install "${EXT}"
completed "${GREEN}lukso${NO_COLOR} installed"

printf '\n'

print_install
