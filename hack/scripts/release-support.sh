#!/bin/bash
#
#   Copyright 2015  Xebia Nederland B.V.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
hasChanges() {
	test -n "$(git status -s .)"
}

getRelease() {
	awk -F= '/^release=/{print $2}' .release
}

getBaseTag() {
		sed -n -e "s/^tag=\(.*\)$(getRelease)\$/\1/p" .release
}

getTag() {
	if [ -z "$1" ] ; then
		awk -F= '/^tag/{print $2}' .release
	else
		echo "$(getBaseTag)$1"
	fi
}

setRelease() {
	if [ -n "$1" ] ; then
		sed -i.x -e "s~^tag=.*~tag=$(getTag $1)~" .release
		sed -i.x -e "s~^release=.*~release=$1~g" .release
		rm -f .release.x
		runPreTagCommand "$1"
	else
		echo "ERROR: missing release version parameter " >&2
		return 1
	fi
}

runPreTagCommand() {
	if [ -n "$1" ] ; then
		COMMAND=$(sed -n -e "s/@@RELEASE@@/$1/g" -e 's/^pre_tag_command=\(.*\)/\1/p' .release)
		if [ -n "$COMMAND" ] ; then
			if ! OUTPUT=$(bash -c "$COMMAND" 2>&1) ; then echo $OUTPUT >&2 && exit 1 ; fi
		fi
	else
		echo "ERROR: missing release version parameter " >&2
		return 1
	fi
}

tagExists() {
	tag=${1:-$(getTag)}
	test -n "$tag" && test -n "$(git tag | grep "^$tag\$")"
}

differsFromRelease() {
	tag=$(getTag)
	! tagExists $tag || test -n "$(git diff --shortstat -r $tag .)"
}

getVersion() {
	result=$(getRelease)

	if differsFromRelease; then
		result="$result-$(git log -n 1 --format=%h .)"
	fi

	if hasChanges ; then
		result="$result-dirty"
	fi
	echo $result
}

nextPatchLevel() {
	version=${1:-$(getRelease)}
	major_and_minor=$(echo $version | cut -d. -f1,2)
	patch=$(echo $version | cut -d. -f3)
	version=$(printf "%s.%d" $major_and_minor $(($patch + 1)))
	echo $version
}

nextMinorLevel() {
	version=${1:-$(getRelease)}
	major=$(echo $version | cut -d. -f1);
	minor=$(echo $version | cut -d. -f2);
	version=$(printf "%d.%d.0" $major $(($minor + 1))) ;
	echo $version
}

nextMajorLevel() {
	version=${1:-$(getRelease)}
	major=$(echo $version | cut -d. -f1);
	version=$(printf "%d.0.0" $(($major + 1)))
	echo $version
}
