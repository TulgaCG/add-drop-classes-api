package database

import _ "embed"

//go:embed schema.sql
var Schema string

//go:embed mockdata.sql
var Mockdata string
