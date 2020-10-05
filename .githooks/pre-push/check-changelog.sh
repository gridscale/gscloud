#!/bin/sh
############################################################
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin
############################################################
DATE=$(date +%Y-%m-%d)
GIT_CMD="git -P"
GIT_COMMIT_SHORT=$($GIT_CMD rev-parse HEAD | cut -b 1-7)
GIT_TAG=$($GIT_CMD describe --abbrev=0 --tags)
GIT_TAG_COMMIT=$($GIT_CMD describe --abbrev=7 --tags)
############################################################
# ( $GIT_CMD describe )
############################################################
CHANGELOG="CHANGELOG.md"
has_CHANGELOG=$($GIT_CMD ls-tree --full-tree -r --name-only HEAD | grep -i $CHANGELOG)
if [ ! -z $has_CHANGELOG ]; then
  ( grep $GIT_TAG $CHANGELOG && \
    echo $GIT_TAG_COMMIT | grep -v $GIT_COMMIT_SHORT ) >/dev/null && exit 0 || exit 123
  ( grep "${DATE}" $CHANGELOG ) || exit 123
else
  echo "Please tag this release, and rework CHANGELOG.md to add a date."
  exit 666
fi
