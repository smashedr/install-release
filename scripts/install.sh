#!/usr/bin/env bash
# https://raw.githubusercontent.com/smashedr/install-release/refs/heads/master/scripts/install.sh

OWNER="smashedr"
REPO="install-release"
EXE="ir"

echo "Installing: ${OWNER}/${REPO} as ${EXE}"

TARGET_BIN="${HOME}/bin"
echo "TARGET_BIN: ${TARGET_BIN}"

set -e

WORK_DIR=$(mktemp -d -t install-release-XXXXXXXXXX 2>&1)
echo "WORK_DIR: ${WORK_DIR}"

function _execution_trap() {
	_ST="$?"
	rm -rf "${WORK_DIR}"
	if [[ "${_ST}" != "0" ]];then
	    echo "⛔ Installation Error! Exit Code: ${_ST}"
	fi
	exit "${_ST}"
}

trap _execution_trap EXIT HUP INT QUIT PIPE TERM

function fail() {
	_ST="$?"
	echo "_ST: ${_ST} - Message: ${1}" 1>&2
	exit ${_ST}
}

# CHECKS
[ ! "$BASH_VERSION" ] && fail "Please use bash instead"
which find > /dev/null || fail "find not installed"
which xargs > /dev/null || fail "xargs not installed"
which sort > /dev/null || fail "sort not installed"
which tail > /dev/null || fail "tail not installed"
which cut > /dev/null || fail "cut not installed"
which du > /dev/null || fail "du not installed"


# PATH
if ! echo "${PATH}" | tr ':' '\n' | grep -qx "${TARGET_BIN}"; then
    echo "⚠️ Target bin NOT in PATH: ${TARGET_BIN}"
    # TODO: Add TARGET_BIN to PATH
fi


# GET
GET=""
if which curl > /dev/null; then
    GET="curl --fail -# -L"
elif which wget > /dev/null; then
    GET="wget -qO-"
else
    fail "neither curl or wget are installed"
fi
if [[ -n "${GITHUB_TOKEN}" ]]; then
    GET="$GET -H 'Authorization: ${GITHUB_TOKEN}'"
fi
echo "GET: ${GET}"  # TODO: Exposes GITHUB_TOKEN


# OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
darwin)
    OS="Darwin"
    FTYPE="tar.gz"
    ;;
linux)
    OS="Linux"
    FTYPE="tar.gz"
    ;;
*)
    fail "unknown os: $(uname -s)"
    ;;
esac
# TODO: Add windows support?
echo "OS: ${OS}"


# ARCH
ARCH="$(uname -m)"
if [[ $OS = "darwin" ]] && sysctl hw.optional.arm64 2>/dev/null | grep -q ': 1'; then
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


# URL
URL="https://github.com/${OWNER}/${REPO}/releases/latest/download/${EXE}_${OS}_${ARCH}.${FTYPE}"
echo "URL: ${URL}"


# DOWNLOAD
echo "cd ${WORK_DIR}"
cd "${WORK_DIR}"
if [[ $FTYPE = "tar.gz" ]] || [[ $FTYPE = ".tgz" ]]; then
    which tar > /dev/null || fail "tar is not installed"
    which gzip > /dev/null || fail "gzip is not installed"
    bash -c "$GET $URL" | tar zxf - || fail "download failed"
elif [[ $FTYPE = "zip" ]]; then
    which unzip > /dev/null || fail "unzip is not installed"
    bash -c "$GET $URL" > tmp.zip || fail "download failed"
    unzip -o -qq tmp.zip || fail "unzip failed"
    rm tmp.zip || fail "cleanup failed"
else
    fail "unknown file type: $FTYPE"
fi


# CHMOD
TMP_BIN="${WORK_DIR}/${EXE}"
echo "TMP_BIN: ${TMP_BIN}"
chmod +x "${TMP_BIN}" || fail "chmod +x failed"
DEST="${TARGET_BIN}/${EXE}"
echo "DEST: ${DEST}"


# MOVE
mv "${TMP_BIN}" "${DEST}"

#OUT=$(mv "${TMP_BIN}" "${DEST}" 2>&1)
#STATUS=$?
#echo "OUT: ${OUT}"
#echo "STATUS: ${STATUS}"
#if [[ "${STATUS}" != "0" ]]; then
#    if [[ $OUT =~ "Permission denied" ]]; then
#        echo "mv with sudo..."
#        sudo mv "${TMP_BIN}" "${DEST}" || fail "sudo mv failed"
#    fi
#fi


echo "✅ Successfully Installed: ${EXE}"
