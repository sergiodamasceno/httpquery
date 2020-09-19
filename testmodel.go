package httpquery

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"gopkg.in/guregu/null.v3"
)

type TestModel struct {
	ID        int               `json:"id,omitempty" db:"id"`                 // id
	UUID      uuid.UUID         `json:"uuid,omitempty" db:"uuid"`             // uuid
	CreatedAt null.Time         `json:"created_at,omitempty" db:"created_at"` // created_at
	UpdatedAt null.Time         `json:"updated_at,omitempty" db:"updated_at"` // updated_at
	Notes     *[]sql.NullString `json:"notes,omitempty" db:"notes"`           // notes
	StartDate null.Time         `json:"start_date,omitempty" db:"start_date"` // start_date
	EndDate   null.Time         `json:"end_date,omitempty" db:"end_date"`     // end_date
	Name      string            `json:"name,omitempty" db:"name"`             // name
}

var TestModelQueryColumns = map[string]struct{}{
	"id":              struct{}{},
	"uuid":            struct{}{},
	"created_at":      struct{}{},
	"updated_at":      struct{}{},
	"notes":           struct{}{},
	"start_date":      struct{}{},
	"end_date":        struct{}{},
	"name":            struct{}{},
	"project_type_id": struct{}{},
}
