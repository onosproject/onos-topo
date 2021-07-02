// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package utils

import (
	"fmt"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
)

const (
	DefaultServicePort = 5150
	DefaultServiceHost = "onos-topo"
)

// Options topo SDK options
type Options struct {
	// Service service options
	Service ServiceOptions
}

// WithTopoAddress sets the address for the topo service
func WithTopoAddress(host string, port int) Option {
	return newOption(func(options *Options) {
		options.Service.Host = host
		options.Service.Port = port
	})
}

// WithTopoHost sets the host for the topo service
func WithTopoHost(host string) Option {
	return newOption(func(options *Options) {
		options.Service.Host = host
	})
}

// WithTopoPort sets the port for the topo service
func WithTopoPort(port int) Option {
	return newOption(func(options *Options) {
		options.Service.Port = port
	})
}

// ServiceOptions are the options for a service
type ServiceOptions struct {
	// Host is the service host
	Host string
	// Port is the service port
	Port int

	Insecure bool
}

// GetHost gets the service host
func (o ServiceOptions) GetHost() string {
	return o.Host
}

// GetPort gets the service port
func (o ServiceOptions) GetPort() int {
	if o.Port == 0 {
		return DefaultServicePort
	}
	return o.Port
}

// IsInsecure is topo connection secure
func (o ServiceOptions) IsInsecure() bool {
	return o.Insecure
}

// GetAddress gets the service address
func (o ServiceOptions) GetAddress() string {
	return fmt.Sprintf("%s:%d", o.GetHost(), o.GetPort())
}

// Option topo client
type Option interface {
	apply(*Options)
}

type funcOption struct {
	f func(*Options)
}

func (f funcOption) apply(options *Options) {
	f.f(options)
}

func newOption(f func(*Options)) Option {
	return funcOption{
		f: f,
	}
}

// WithOptions sets the client options
func WithOptions(opts Options) Option {
	return newOption(func(options *Options) {
		*options = opts
	})
}

// WatchOptions topo client watch method options
type WatchOptions struct {
	filters *topoapi.Filters
}

// GetFilters get filters
func (w WatchOptions) GetFilters() *topoapi.Filters {
	return w.filters
}

// WatchOption topo client watch option
type WatchOption interface {
	apply(*WatchOptions)
}

type funcWatchOption struct {
	f func(*WatchOptions)
}

func (f funcWatchOption) apply(options *WatchOptions) {
	f.f(options)
}

func newWatchOption(f func(*WatchOptions)) WatchOption {
	return funcWatchOption{
		f: f,
	}
}

// WithWatchFilters sets filters for watch method
func WithWatchFilters(filters *topoapi.Filters) WatchOption {
	return newWatchOption(func(o *WatchOptions) {
		o.filters = filters
	})
}

// ListOptions topo client get method options
type ListOptions struct {
	filters *topoapi.Filters
}

// ListOption topo client list option
type ListOption interface {
	apply(*ListOptions)
}

type funcListOption struct {
	f func(options *ListOptions)
}

func (f funcListOption) apply(options *ListOptions) {
	f.f(options)
}

func newListOption(f func(options *ListOptions)) ListOption {
	return funcListOption{
		f: f,
	}
}

// GetFilters get filters
func (l ListOptions) GetFilters() *topoapi.Filters {
	return l.filters
}

// WithListFilters sets filters for list method
func WithListFilters(filters *topoapi.Filters) ListOption {
	return newListOption(func(o *ListOptions) {
		o.filters = filters
	})

}
