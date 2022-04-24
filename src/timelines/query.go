package timelines

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/orm"
)

type ConnAttempSimple struct {
	tableName   struct{} `pg:",discard_unknown_columns"`
	ID          uint
	Time        time.Time
	Port        string
	IP          string
	CountryCode string
	ClientPort  string
}

type CountResult struct {
	Count int
}

type MapDataEntry struct {
	CountryCode string
	Count       int
}

type PortCount struct {
	Port  string
	Count int64
}

type BytesList struct {
	tableName struct{} `pg:",discard_unknown_columns"`
	Time      time.Time
	Bytes     string
}

func (t *timelines) InsertConnAttemp(connAttemp *ConnAttemp) error {
	err := t.db.Insert(connAttemp)

	if err != nil {
		return err
	}
	return nil
}

func (t *timelines) GetTotalConsumptions(ctx context.Context, rangeValue string) (*CountResult, error) {
	var err error
	query := t.db.Model(&ConnAttemp{})

	query, err = addRange(query, "conn_attemp.time", rangeValue)
	if err != nil {
		return nil, err
	}

	count, err := query.Count()
	if err != nil {
		return nil, err
	}

	return &CountResult{Count: count}, nil
}

func (t *timelines) GetMapData(ctx context.Context, rangeValue string) ([]*MapDataEntry, error) {
	var res []*MapDataEntry

	query, err := t.makeCountQuery("country_code", rangeValue)
	if err != nil {
		return nil, err
	}

	if err := query.Select(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (t *timelines) GetConnAttemps(ctx context.Context, rangeValue string) ([]*ConnAttempSimple, error) {
	var err error
	var res []*ConnAttempSimple

	query := t.db.Model((*ConnAttemp)(nil))

	query, err = addRange(query, "conn_attemp.time", rangeValue)
	if err != nil {
		return nil, err
	}

	query = query.Order("conn_attemp.time DESC")

	err = query.Select(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *timelines) GetTopConsumers(ctx context.Context, rangeValue string) ([]*MapDataEntry, error) {
	var res []*MapDataEntry

	query, err := t.makeCountQuery("country_code", rangeValue)
	if err != nil {
		return nil, err
	}

	if err := query.Select(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (t *timelines) GetTopFlavours(ctx context.Context, rangeValue string) ([]*PortCount, error) {
	var res []*PortCount

	query, err := t.makeCountQuery("port", rangeValue)
	if err != nil {
		return nil, err
	}

	if err := query.Select(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (t *timelines) GetBytes(ctx context.Context, rangeValue, port string) ([]*BytesList, error) {
	var res []*BytesList

	err := t.db.Model((*ConnAttemp)(nil)).
		Where("conn_attemp.port = ?", port).
		Select(&res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *timelines) makeCountQuery(col, rangeValue string) (*orm.Query, error) {
	query := t.db.Model((*ConnAttemp)(nil)).
		Column(col).
		ColumnExpr("count(*) AS count").
		Group(col).
		Order(col)

	return addRange(query, "conn_attemp.time", rangeValue)
}

func addRange(query *orm.Query, col, rangeValue string) (*orm.Query, error) {
	rangeTime, err := rangeValueToTime(rangeValue)
	if err != nil {
		return nil, err
	}
	return query.Where(col+" >= ?", rangeTime), nil
}

func rangeValueToTime(rangeValue string) (time.Time, error) {
	now := time.Now()
	var diff time.Duration

	switch rangeValue {
	case "y":
		// Approximation only
		diff = time.Hour * 24 * 365
	case "mo":
		// Approximation only
		diff = time.Hour * 24 * 31
	case "w":
		diff = time.Hour * 24 * 7
	case "d":
		diff = time.Hour * 24
	case "h":
		diff = time.Hour
	default:
		return time.Time{}, fmt.Errorf("Invalid rangeValue")
	}

	return now.Add(diff * -1), nil
}
