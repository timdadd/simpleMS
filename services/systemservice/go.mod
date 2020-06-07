module systemservice

go 1.14

replace lib/common v0.0.0 => ./lib/common

require (
	github.com/stretchr/testify v1.4.0
	lib/common v0.0.0
)
