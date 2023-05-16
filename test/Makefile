# SPDX-FileCopyrightText: 2023-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0


SHELL = bash -e -o pipefail

export CGO_ENABLED=1
export GO111MODULE=on

helmit-topo: # @HELP run helmit topo tests locally
	kubectl delete ns test 2>/dev/null || exit 0
	helmit test --create-namespace --namespace test --no-teardown --context ../ ./topo/... --set global.image.registry=mirror.registry.opennetworking.org --image mirror.registry.opennetworking.org/onosproject/helmit-runner:latest --set onos-topo.global.image.tag=latest --timeout 15m
