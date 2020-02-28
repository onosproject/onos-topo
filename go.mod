module github.com/onosproject/onos-topo

go 1.13

require (
	github.com/atomix/api v0.0.0-20200211005812-591fe8b07ea8
	github.com/atomix/go-client v0.0.0-20200211010855-927b10345735
	github.com/atomix/go-framework v0.0.0-20200211010411-ae512dcee9ad
	github.com/atomix/go-local v0.0.0-20200211010611-c99e53e4c653
	github.com/gogo/protobuf v1.3.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onosproject/onos-lib-go v0.0.0-20200224171112-d46f89d458c0
	github.com/onosproject/onos-test v0.0.0-20200212201952-fb8d2ac644a0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.5.0
	github.com/stretchr/testify v1.4.0
	google.golang.org/grpc v1.27.1
	gotest.tools v2.2.0+incompatible
	k8s.io/klog v1.0.0
)

replace github.com/onosproject/onos-lib-go => ../onos-lib-go
