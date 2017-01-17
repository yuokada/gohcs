#!/bin/bash -x
WORK_DIR=`pwd`
BUILD_NAME=gohcs
BUILD_SPEC=rpm/${BUILD_NAME}.spec
BUILD_TAR=${BUILD_NAME}.tar.gz
BUILD_DIR=${WORK_DIR}/rpmbuild
BUILD_TAR_PATH=${WORK_DIR}/rpmbuild/SOURCES/${BUILD_TAR}
BUILD_NUMBER=${BUILD_NUMBER:=0}

ARCHIVE_DIR=${WORK_DIR}/${BUILD_NAME}

/bin/rm   -rf ${BUILD_DIR}
/bin/mkdir -p ${BUILD_DIR}/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

/bin/mkdir -p ${ARCHIVE_DIR}{/etc/systemd/system/,/etc/tmpfiles.d/,/etc/gohcs,/var/run/gohcs,/src}
/bin/cp  ./etc/tmpfiles.d/gohcs.conf        ${ARCHIVE_DIR}/etc/tmpfiles.d/
/bin/cp  ./etc/systemd/system/gohcs.service ${ARCHIVE_DIR}/etc/systemd/system/
/bin/cp  ./etc/checklist.json               ${ARCHIVE_DIR}/etc/gohcs/checklist.json
/bin/cp  ./src/server.go                    ${ARCHIVE_DIR}/src
/bin/cp  ./Makefile                         ${ARCHIVE_DIR}/
/bin/tar czvf ${BUILD_TAR} gohcs/
/bin/rm  -rf gohcs

#/usr/bin/wget https://github.com/yuokada/gohcs/archive/master.tar.gz -O ${BUILD_TAR}
/bin/mv ${BUILD_TAR} ${BUILD_TAR_PATH}


/usr/bin/rpmbuild --define "release ${BUILD_NUMBER}%{?dist}" --define "_topdir ${BUILD_DIR}" -bb ${BUILD_SPEC}

