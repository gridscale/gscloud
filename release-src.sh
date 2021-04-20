#!/bin/sh

REPO=https://github.com/gridscale/gscloud.git
VERSION=`git describe --tags`
GIT_COMMIT=`git rev-list -1 HEAD`

CURDIR=`pwd`
RELDIR=$(mktemp -d)

cd $RELDIR

git clone $REPO
cd gscloud
git checkout $1

sed -e "s/VERSION=\$\$(git describe --tags)/VERSION=${VERSION}/g" -e "s/GIT_COMMIT=\$\$(git rev-list -1 HEAD)/GIT_COMMIT=${GIT_COMMIT}/g"  Makefile > Makefile.tmp
mv Makefile.tmp Makefile

rm -rf .git

mkdir ${RELDIR}/gscloud_${VERSION}

cp -R . ${RELDIR}/gscloud_${VERSION}

cd ${RELDIR} && tar czfv gscloud_${VERSION}.tgz gscloud_${VERSION}/
cd ${RELDIR} && zip -r gscloud_${VERSION}.zip gscloud_${VERSION}/

cd ${CURDIR}
mkdir -p release

cp ${RELDIR}/gscloud_${VERSION}.tgz release/
cp ${RELDIR}/gscloud_${VERSION}.zip release/

rm -r ${RELDIR}

