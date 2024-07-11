package callback

import (
	"context"
	"fmt"
	"reflect"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/rs/xid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Callback struct {
}

func New() *Callback {
	return &Callback{}
}

func setField(ctx context.Context, field *schema.Field, rv reflect.Value) {
	if field != nil {
		if _, isZero := field.ValueOf(ctx, rv); isZero {
			var xid = xid.New().String()
			field.Set(ctx, rv, xid)
		}
	}

}

func GenerateXID(db *gorm.DB) {
	if db.Error == nil && db.Statement.Schema != nil {
		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				var rv = db.Statement.ReflectValue.Index(i)
				var field = db.Statement.Schema.LookUpField("ID")
				setField(db.Statement.Context, field, rv)
			}
		case reflect.Struct:
			var field = db.Statement.Schema.LookUpField("ID")
			setField(db.Statement.Context, field, db.Statement.ReflectValue)
		}
	}
}

func GenerateBlurhash(db *gorm.DB) {
	if db.Statement.Schema != nil {
		// crop image fields and upload them to CDN, dummy code
		for _, field := range db.Statement.Schema.Fields {
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					// Get value from field
					if fieldValue, isZero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue.Index(i)); !isZero {
						if attachment, ok := fieldValue.(*models.Attachment); ok {
							var blurhash = attachment.GetBlurhash()

							err := field.Set(db.Statement.Context, db.Statement.ReflectValue.Index(i), attachment)
							if err != nil {
								fmt.Println("callback:GenerateBlurhash set slice value err", i, blurhash, err)
							}
						}
					}
				}
			case reflect.Struct:
				// Get value from field
				if fieldValue, isZero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue); !isZero {
					if attachment, ok := fieldValue.(*models.Attachment); ok {
						var blurhash = attachment.GetBlurhash()
						// Set value to field
						err := field.Set(db.Statement.Context, db.Statement.ReflectValue, attachment)
						if err != nil {
							fmt.Println("callback:GenerateBlurhash set structure value err", blurhash, err)
						}
					}
				}

			}
		}

	}
}

func (c *Callback) Register(db *gorm.DB) {
	c.RegisterGenerateXID(db)
	c.RegisterGenerateBlurhash(db)
}

func (c *Callback) RegisterGenerateXID(db *gorm.DB) {
	db.Callback().Create().Before("gorm:save_before_associations").Register("app:update_xid_when_create", GenerateXID)
}

func (c *Callback) RegisterGenerateBlurhash(db *gorm.DB) {
	db.Callback().Create().Before("gorm:save_before_associations").Register("app:gen_blurhash", GenerateBlurhash)
}
