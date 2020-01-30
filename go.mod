module github.com/rhdedgar/pleg-watcher

go 1.12

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

require (
	github.com/coreos/go-systemd v0.0.0-00010101000000-000000000000
	github.com/coreos/go-systemd/v22 v22.0.0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/openshift/clam-scanner v0.0.0-20170918135446-9f39c23ef966
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	golang.org/x/sys v0.0.0-20200124204421-9fbb57f87de9
)
