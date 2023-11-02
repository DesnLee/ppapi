package pkg

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func CreatePgTimeTZ(s string) (pgtype.Timestamptz, error) {
	tm, err := time.Parse("2006-01-02", s)
	if err != nil {
		return pgtype.Timestamptz{}, fmt.Errorf("时间格式错误")
	}

	var pgTm pgtype.Timestamptz
	err = pgTm.Scan(tm)
	if err != nil {
		return pgtype.Timestamptz{}, fmt.Errorf("时间格式错误")
	}

	return pgTm, nil
}
