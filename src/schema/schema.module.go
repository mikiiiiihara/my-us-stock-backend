package schema

import (
	"my-us-stock-backend/src/schema/generated"
	"my-us-stock-backend/src/schema/user"

	"gorm.io/gorm"
)

type SchemaModule struct {
	userModule *user.UserModule
}

func NewSchemaModule(db *gorm.DB) *SchemaModule {
	userModule := user.NewUserModule(db)
	// 他のモジュールの初期化

	return &SchemaModule{
		userModule: userModule,
		// 他のモジュールのインスタンス化
	}
}

func (r *SchemaModule) Query() generated.QueryResolver {
	return r.userModule.Query()
}

func (r *SchemaModule) Mutation() generated.MutationResolver {
	return r.userModule.Mutation()
}
