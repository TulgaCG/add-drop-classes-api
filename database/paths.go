package database

import _ "embed"

//go:embed schema.sql
var Schema string

//go:embed roles.sql
var Roles string

//go:embed lectures.sql
var Lectures string
