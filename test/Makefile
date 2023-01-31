# SPDX-FileCopyrightText: 2023-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0


SHELL = bash -e -o pipefail

export CGO_ENABLED=1
export GO111MODULE=on

helmit-topo: # @HELP run helmit topo tests locally
	helmit test -n test ./onos-topo-tests --suite topo --no-teardown --timeout 40m