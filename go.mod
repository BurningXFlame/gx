module github.com/burningxflame/gx

go 1.18

require (
	github.com/fsnotify/fsnotify v1.7.0
	golang.org/x/exp v0.0.0-20230224173230-c95f2b4c22f2 // v0.0.0-20230224173230-c95f2b4c22f2 for go 1.18
	sigs.k8s.io/yaml v1.4.0
)

// test dependencies
require github.com/stretchr/testify v1.9.0

// indirect dependencies
require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
