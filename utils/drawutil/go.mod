module github.com/jiushengTech/common/utils/drawutil

go 1.24.1

toolchain go1.24.2

require (
	github.com/fogleman/gg v1.3.0
	github.com/jiushengTech/common v0.0.0-20250505175851-7adfce197c0c
	golang.org/x/image v0.28.0
)

require github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect

replace github.com/jiushengTech/common => ../../
replace github.com/jiushengTech/common/utils/drawutil => ./
