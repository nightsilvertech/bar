package data

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	pb "github.com/nightsilvertech/bar/protoc/api/v1"
	_interface "github.com/nightsilvertech/bar/repository/interface"
	"github.com/nightsilvertech/utl/errwrap"
	"go.opencensus.io/trace"
	"sync"
	"time"
)

var mutex = &sync.RWMutex{}

type dataReadWrite struct {
	tracer trace.Tracer
	db     *sql.DB
}

func (d *dataReadWrite) WriteBar(ctx context.Context, req *pb.Bar) (res *pb.Bar, err error) {
	const funcName = `WriteBar`
	ctx, span := d.tracer.StartSpan(ctx, funcName)
	defer span.End()

	currentTime := time.Now()
	req.CreatedAt = currentTime.Unix()
	req.UpdatedAt = currentTime.Unix()
	stmt, err := d.db.Prepare(`
	INSERT INTO bars(id, name, description, created_at, updated_at) VALUES (?,?,?,?,?)
	`)
	if err != nil {
		return res, errwrap.Wrap(funcName, "db.Prepare", err)
	}
	result, err := stmt.ExecContext(
		ctx,
		req.Id,          // id
		req.Name,        // name
		req.Description, // description
		currentTime,     // created_at
		currentTime,     // updated_at
	)
	if err != nil {
		return res, errwrap.Wrap(funcName, "stmt.ExecContext", err)
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return res, errwrap.Wrap(funcName, "result.RowsAffected", err)
	}
	return req, nil
}

func (d *dataReadWrite) ModifyBar(ctx context.Context, req *pb.Bar) (res *pb.Bar, err error) {
	const funcName = `ModifyBar`
	ctx, span := d.tracer.StartSpan(ctx, funcName)
	defer span.End()

	currentTime := time.Now()
	req.UpdatedAt = currentTime.Unix()
	stmt, err := d.db.Prepare(`
	UPDATE bars
	SET name = ?, description = ?
	WHERE id = ?
	`)
	if err != nil {
		return res, errwrap.Wrap(funcName, "db.Prepare", err)
	}
	result, err := stmt.ExecContext(
		ctx,
		req.Name,        // name
		req.Description, // description
		req.Id,          // id
	)
	if err != nil {
		return res, errwrap.Wrap(funcName, "stmt.ExecContext", err)
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return res, errwrap.Wrap(funcName, "result.RowsAffected", err)
	}
	return req, nil
}

func (d *dataReadWrite) RemoveBar(ctx context.Context, req *pb.Select) (res *pb.Bar, err error) {
	const funcName = `RemoveBar`
	ctx, span := d.tracer.StartSpan(ctx, funcName)
	defer span.End()

	stmt, err := d.db.Prepare(`DELETE FROM bars WHERE id = ?`)
	if err != nil {
		return res, errwrap.Wrap(funcName, "db.Prepare", err)
	}
	result, err := stmt.ExecContext(
		ctx,
		req.Id, // id
	)
	if err != nil {
		return res, errwrap.Wrap(funcName, "stmt.ExecContext", err)
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return res, errwrap.Wrap(funcName, "result.RowsAffected", err)
	}
	return res, nil
}

func (d *dataReadWrite) ReadDetailBar(ctx context.Context, selects *pb.Select) (res *pb.Bar, err error) {
	const funcName = `ReadDetailBar`
	ctx, span := d.tracer.StartSpan(ctx, funcName)
	defer span.End()

	stmt, err := d.db.Prepare(`SELECT * FROM bars WHERE id = ?`)
	if err != nil {
		return res, errwrap.Wrap(funcName, "db.Prepare", err)
	}
	mutex.Lock()
	row := stmt.QueryRowContext(ctx, selects.Id)
	mutex.Unlock()

	var bar pb.Bar
	var createdAt, updatedAt time.Time
	err = row.Scan(
		&bar.Id,          // id
		&bar.Name,        // name
		&bar.Description, // description
		&createdAt,       // created_at
		&updatedAt,       // updated_at
	)
	if err != nil {
		return res, errwrap.Wrap(funcName, "row.Scan", err)
	}
	bar.CreatedAt = createdAt.Unix()
	bar.UpdatedAt = updatedAt.Unix()
	return &bar, nil
}

func (d *dataReadWrite) ReadAllBar(ctx context.Context, req *pb.Pagination) (res *pb.Bars, err error) {
	const funcName = `ReadAllBar`
	ctx, span := d.tracer.StartSpan(ctx, funcName)
	defer span.End()

	stmt, err := d.db.Prepare(`SELECT * FROM bars ORDER BY created_at DESC`)
	if err != nil {
		return res, errwrap.Wrap(funcName, "db.Prepare", err)
	}
	mutex.Lock()
	row, err := stmt.QueryContext(ctx)
	if err != nil {
		return res, errwrap.Wrap(funcName, "stmt.QueryContext", err)
	}
	mutex.Unlock()
	defer row.Close()

	var bars pb.Bars
	var bar pb.Bar
	var createdAt, updatedAt time.Time
	for row.Next() {
		err = row.Scan(
			&bar.Id,          // id
			&bar.Name,        // name
			&bar.Description, // description
			&createdAt,       // created_at
			&updatedAt,       // updated_at
		)
		if err != nil {
			return res, errwrap.Wrap(funcName, "row.Scan", err)
		}
		bars.Bars = append(bars.Bars, &pb.Bar{
			Id:          bar.Id,
			Name:        bar.Name,
			Description: bar.Description,
			CreatedAt:   createdAt.Unix(),
			UpdatedAt:   updatedAt.Unix(),
		})
	}
	return &bars, nil
}

func NewDataReadWriter(username, password, host, port, name string, tracer trace.Tracer) (_interface.DRW, error) {
	const funcName = `NewDataReadWriter`

	databaseUrl := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		username,
		password,
		host,
		port,
		name,
	)
	db, err := sql.Open("mysql", databaseUrl)
	if err != nil {
		return nil, errwrap.Wrap(funcName, "sql.Open", err)
	}

	return &dataReadWrite{
		tracer: tracer,
		db:     db,
	}, nil
}
