module github.com/matthewhartstonge/mongo-features/tldr

go 1.13

// use the local code, rather than go'getting the module
replace github.com/matthewhartstonge/mongo-features => ../../../mongo-features

require (
	github.com/matthewhartstonge/mongo-features v0.4.0
	go.mongodb.org/mongo-driver v1.4.0
)
