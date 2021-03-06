SHELL=/bin/bash
VERSION=$(shell cat VERSION)
MAJOR_VERSION=$(shell awk -F"." '{ print $$1 }' VERSION)
MINOR_VERSION=$(shell awk -F"." '{ print $$2 }' VERSION)
REVISION_VERSION=$(shell awk -F"." '{ print $$3 }' VERSION)
NEW_REVISION_VERSION=$$((${REVISION_VERSION} + 1))
NEW_MINOR_VERSION=$$((${MINOR_VERSION} + 1))
NEW_VERSION_MINOR="${MAJOR_VERSION}.${NEW_MINOR_VERSION}.0"
NEW_VERSION_REV:="${MAJOR_VERSION}.${MINOR_VERSION}.${NEW_REVISION_VERSION}"
BUILD_ID=$(shell git log -1 --pretty=format:"%H")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# -s and -w flags to remove unneccesary debug information
LDFLAGS="-s -w -X BuildVersion=${NEW_VERSION_REV} -X BuildTime=${BUILD_TIME} -X BuildID=${BUILD_ID}"

.PHONY : build
build : git-status clean inc-revision build-dir
	go build -ldflags ${LDFLAGS} -v -o build/starchive main.go starchive.go


.PHONY : build-small
build-small : build
# in order to shrink the binary further, we pack the binary via upx
# this can take some time...
	upx --brute build/starchive


.PHONY : build-dir
build-dir : 
# if the build directory doesn't exist, then make it.
	if [ ! -d "build" ]; then mkdir build; fi


.PHONY : clean
clean : 
	$(shell rm -rf build )


.PHONY : inc-revision
inc-revision : 
# increase the revision number and write it out to the VERSION file.
	@echo "${NEW_VERSION_REV}" > VERSION


.PHONY : inc-minor
inc-minor :
# increase the minor version number and write it out to the VERSION file.
	@echo "${NEW_VERSION_MINOR}" > VERSION


.PHONY : git-status
git-status:
	@status=$$(git status --porcelain); \
	if [ ! -z "$${status}" ]; \
	then \
		echo "Error - working directory is dirty. Commit those changes!"; \
		exit 1; \
	fi
