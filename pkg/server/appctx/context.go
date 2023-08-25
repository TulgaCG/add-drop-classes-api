package appctx

import (
	"log/slog"

	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
)

// Context consists application related references to services such as logger, database
// Instead of adding these services individually, one should add every instance into this.
type Context struct {
	DB  *database.DB
	Log *slog.Logger
}
