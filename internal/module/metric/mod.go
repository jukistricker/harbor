package metric

import (
	"database/sql"
	"harbor/internal/module/metric/service"
)

type Module struct {
	Service *service.Processor
	//Thêm Repository hoặc Handler vào đây sau này
}

func NewModule(db *sql.DB) *Module {
	// Khởi tạo Processor(xử lý metric)
	proc := service.NewProcessor(db)

	return &Module{
		Service: proc,
	}
}