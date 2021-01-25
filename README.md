# onos-topo
[![Build Status](https://travis-ci.com/onosproject/onos-topo.svg?branch=master)](https://travis-ci.com/onosproject/onos-topo)
[![Integration Test Status](https://img.shields.io/travis/onosproject/onos-test?label=integration-tests&logo=integration-tests)](https://travis-ci.com/onosproject/onos-test)
[![Go Report Card](https://goreportcard.com/badge/github.com/onosproject/onos-topo)](https://goreportcard.com/report/github.com/onosproject/onos-topo)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/gojp/goreportcard/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/onosproject/onos-topo?status.svg)](https://godoc.org/github.com/onosproject/onos-topo)

## Overview

The µONOS Topology subsystem provides topology management for µONOS core services and applications. The topology is exposed to services via an implementation of the onos-topo gRPC API as defined by the [onos-api]. The toplogy API is generalized as a set of entities and relations, where `Entity` models objects within the system (e.g. devices, nodes, etc) and `Relation` models the associations between objects.

The topology subsystem is shipped as a [Docker] image and deployed with [Helm]. To build the Docker image, run `make images`.

[onos-api]: https://github.com/onosproject/onos-api
[Docker]: https://www.docker.com/
[Helm]: https://helm.sh
