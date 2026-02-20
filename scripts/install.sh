#!/usr/bin/env bash
# https://raw.githubusercontent.com/smashedr/install-release/refs/heads/master/scripts/install.sh

set -e


OWNER="smashedr"
REPO="install-release"
EXE="ir"
TARGET_BIN="${HOME}/bin"


echo "Installing: ${OWNER}/${REPO} as ${EXE}"


function fail() {
	_ST="$?"
	echo "fail _ST: ${_ST}"
	echo "⛔ ${1}" 1>&2
    if [ "$_ST" -eq 0 ]; then
        exit 1
    fi
	exit "${_ST}"
}


# OS
# TODO: Confirm setting EXE to .exe is not required
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
echo "OS: ${OS}"
case "${OS}" in
darwin)
    OS="Darwin"
    FTYPE="tar.gz"
    ;;
linux)
    OS="Linux"
    FTYPE="tar.gz"
    ;;
mingw* | msys* | cygwin* | windows*)
    OS="Windows"
    FTYPE="zip"
    ;;
*)
    fail "unknown os: $(uname -s)"
    ;;
esac
echo "OS: ${OS}"


# ARCH
ARCH="$(uname -m)"
echo "ARCH: ${ARCH}"
if [[ ${OS} = "Darwin" ]] && sysctl hw.optional.arm64 2>/dev/null | grep -q ': 1'; then
    ARCH="arm64"
elif uname -m | grep -E '(aarch64|arm64)' > /dev/null; then
    ARCH="arm64"
elif uname -m | grep 64 > /dev/null; then
    ARCH="x86_64"
elif uname -m | grep 386 > /dev/null; then
    ARCH="i386"
else
    fail "unknown arch: $(uname -m)"
fi
echo "ARCH: ${ARCH}"


# FTYPE
if [[ ${FTYPE} = "tar.gz" ]]; then
    which tar > /dev/null || fail "tar is not installed"
    which gzip > /dev/null || fail "gzip is not installed"
elif [[ ${FTYPE} = "zip" ]]; then
    which unzip > /dev/null || fail "unzip is not installed"
else
    fail "unknown file type: ${FTYPE}"
fi
echo "FTYPE: ${FTYPE}"


# URL
URL="https://github.com/${OWNER}/${REPO}/releases/latest/download/${EXE}_${OS}_${ARCH}.${FTYPE}"
echo "URL: ${URL}"


# GET
GET=""
if which curl > /dev/null; then
    GET="curl --fail -# -L"
elif which wget > /dev/null; then
    GET="wget -qO-"
else
    fail "neither curl or wget are installed"
fi
#if [[ -n "${GITHUB_TOKEN}" ]]; then
#    GET="${GET} -H 'Authorization: ${GITHUB_TOKEN}'"
#fi
echo "GET: ${GET}"


# BIN
echo "Target Directory: $TARGET_BIN"
echo -n "Enter Path [press <enter> to accept]: "
read -r input </dev/tty
if [[ -n "${input}" ]]; then
    # TODO: Sanatize TARGET_BIN
    TARGET_BIN="${input}"
fi
echo "TARGET_BIN: $TARGET_BIN"


# PATH
if ! echo "${PATH}" | tr ':' '\n' | grep -qx "${TARGET_BIN}"; then
    echo "⚠️ Target bin NOT in PATH: ${TARGET_BIN}"
    # TODO: Add TARGET_BIN to PATH
fi
if [[ ! -d "${TARGET_BIN}" ]]; then
    echo "⚠️ Creating bin directory: ${TARGET_BIN}"
    mkdir -p "${TARGET_BIN}" || fail "mkdir failed on: ${TARGET_BIN}"
fi


# TEMP
TEMP_DIR=$(mktemp -d -t install-release-XXXXXXXXXX 2>&1)
echo "TEMP_DIR: ${TEMP_DIR}"

function _execution_trap() {
	_ST="$?"
	rm -rf "${TEMP_DIR}"
	if [[ "${_ST}" != "0" ]]; then
	    echo "⛔ Installation Error! Exit Code: ${_ST}" 1>&2
	fi
	exit "${_ST}"
}
trap _execution_trap EXIT HUP INT QUIT PIPE TERM


# DOWNLOAD
cd "${TEMP_DIR}"
if [[ ${FTYPE} = "tar.gz" ]]; then
    bash -c "${GET} ${URL}" | tar zxf - || fail "download failed"
elif [[ ${FTYPE} = "zip" ]]; then
    bash -c "${GET} ${URL}" > tmp.zip || fail "download failed"
    unzip -o -qq tmp.zip || fail "unzip failed"
    rm tmp.zip || fail "cleanup failed"
fi


# CHMOD
TEMP_BIN="${TEMP_DIR}/${EXE}"
echo "TEMP_BIN: ${TEMP_BIN}"
chmod +x "${TEMP_BIN}" || fail "chmod +x failed"
DEST="${TARGET_BIN}/${EXE}"
echo "DEST: ${DEST}"


# MOVE
MOVE_OUT=$(mv "${TEMP_BIN}" "${DEST}" 2>&1)
MOVE_ST=$?
echo "MOVE_ST: ${MOVE_ST}"
echo "MOVE_OUT: ${MOVE_OUT}"
if [[ "${MOVE_ST}" != "0" ]]; then
    if [[ ${MOVE_OUT} =~ "Permission denied" ]]; then
        echo "mv with sudo..."
        sudo mv "${TEMP_BIN}" "${DEST}" || fail "sudo mv failed"
    fi
fi


echo "✅ Successfully Installed: ${EXE}"
