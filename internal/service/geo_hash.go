package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gansidui/geohash"
	"github.com/yemingfeng/sdb/internal/collection"
	"github.com/yemingfeng/sdb/internal/pb"
	"google.golang.org/protobuf/proto"
	"math"
	"sort"
	"strconv"
)

var NotFoundGeoHashError = errors.New("not found geo hash, please create it")
var GeoHashExistError = errors.New("geo hash exist, please delete it or change other")
var GeoHashInvalidIdError = errors.New("id can not be equals to key, please change other id")

var geoHashCollection = collection.NewCollection(pb.DataType_GEO_HASH)

func newGeoHashIndexes(hash []byte, box []byte) []collection.Index {
	return []collection.Index{
		{Name: []byte("hash"), Value: hash},
		{Name: []byte("box"), Value: box},
	}
}

func GHCreate(key []byte, precision int32) error {
	lock(LGeoHash, key)
	defer unlock(LGeoHash, key)

	batch := collection.NewBatch()
	defer batch.Close()

	exist, err := geoHashCollection.ExistRowById(key, key)
	if err != nil {
		return err
	}
	if exist {
		return GeoHashExistError
	}

	if err := geoHashCollection.UpsertRow(&collection.Row{
		Key:   key,
		Id:    key,
		Value: []byte(strconv.Itoa(int(precision))),
	}, batch); err != nil {
		return err
	}
	return batch.Commit()
}

func GHDel(key []byte) error {
	lock(LGeoHash, key)
	defer unlock(LGeoHash, key)

	batch := collection.NewBatch()
	defer batch.Close()

	rows, err := geoHashCollection.Page(key, 0, math.MaxUint32)
	if err != nil {
		return err
	}
	for i := range rows {
		if err := geoHashCollection.DelRowById(key, rows[i].Id, batch); err != nil {
			return err
		}
	}
	return batch.Commit()
}

func GHAdd(key []byte, points []*pb.Point) error {
	lock(LGeoHash, key)
	defer unlock(LGeoHash, key)

	batch := collection.NewBatch()
	defer batch.Close()

	row, err := geoHashCollection.GetRowById(key, key)
	if err != nil {
		return err
	}
	if row == nil || len(row.Value) == 0 {
		return NotFoundGeoHashError
	}
	precision, err := strconv.ParseInt(string(row.Value), 10, 32)
	if err != nil {
		return err
	}

	for i := range points {
		point := points[i]
		value, err := proto.Marshal(point)
		if err != nil {
			return err
		}
		if bytes.Equal(point.Id, key) {
			return GeoHashInvalidIdError
		}
		hash, box := geohash.Encode(point.Latitude, point.Longitude, int(precision))
		if err := geoHashCollection.UpsertRow(&collection.Row{
			Key:     key,
			Id:      point.Id,
			Value:   value,
			Indexes: newGeoHashIndexes([]byte(hash), []byte(marshalBox(box))),
		}, batch); err != nil {
			return err
		}
	}
	return batch.Commit()
}

func GHRem(key []byte, ids [][]byte) error {
	lock(LGeoHash, key)
	defer unlock(LGeoHash, key)

	batch := collection.NewBatch()
	defer batch.Close()

	exist, err := geoHashCollection.ExistRowById(key, key)
	if err != nil {
		return err
	}
	if !exist {
		return NotFoundGeoHashError
	}

	for i := range ids {
		if err := geoHashCollection.DelRowById(key, ids[i], batch); err != nil {
			return err
		}
	}
	return batch.Commit()
}

func GHGetBoxes(key []byte, latitude float64, longitude float64) ([]*pb.Point, error) {
	row, err := geoHashCollection.GetRowById(key, key)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, NotFoundGeoHashError
	}
	precision, err := strconv.ParseInt(string(row.Value), 10, 32)
	if err != nil {
		return nil, err
	}

	_, box := geohash.Encode(latitude, longitude, int(precision))

	rows, err := geoHashCollection.IndexValuePage(key, []byte("box"), []byte(marshalBox(box)), 0, math.MaxUint32)
	if err != nil {
		return nil, err
	}

	res := make([]*pb.Point, len(rows))
	for i := range rows {
		var point pb.Point
		if err := proto.Unmarshal(rows[i].Value, &point); err != nil {
			return nil, err
		}
		point.Distance = distance(&point, &pb.Point{Latitude: latitude, Longitude: longitude})
		res[i] = &point
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Distance < res[j].Distance
	})

	return res, nil
}

func GHGetNeighbors(key []byte, latitude float64, longitude float64) ([]*pb.Point, error) {
	row, err := geoHashCollection.GetRowById(key, key)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, NotFoundGeoHashError
	}
	precision, err := strconv.ParseInt(string(row.Value), 10, 32)
	if err != nil {
		return nil, err
	}

	res := make([]*pb.Point, 0)
	neighbors := geohash.GetNeighbors(latitude, longitude, int(precision))
	for i := range neighbors {
		rows, err := geoHashCollection.IndexValuePage(key, []byte("hash"), []byte(neighbors[i]), 0, math.MaxUint32)
		if err != nil {
			return nil, err
		}
		for i := range rows {
			var point pb.Point
			if err := proto.Unmarshal(rows[i].Value, &point); err != nil {
				return nil, err
			}
			point.Distance = distance(&point, &pb.Point{Latitude: latitude, Longitude: longitude})
			res = append(res, &point)
		}
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Distance < res[j].Distance
	})

	return res, nil
}

func GHCount(key []byte) (uint32, error) {
	// delete for meta
	count, err := geoHashCollection.Count(key)
	if err != nil {
		return 0, err
	}
	return count - 1, err
}

func GHMembers(key []byte) ([]*pb.Point, error) {
	rows, err := geoHashCollection.IndexPage(key, []byte("hash"), 0, math.MaxUint32)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Point, len(rows))
	for i := range rows {
		var point pb.Point
		if err := proto.Unmarshal(rows[i].Value, &point); err != nil {
			return nil, err
		}
		res[i] = &point
	}
	return res, nil
}

func marshalBox(box *geohash.Box) string {
	return fmt.Sprintf("%32.32f:%32.32f:%32.32f:%32.32f", box.MinLat, box.MaxLat, box.MinLng, box.MaxLng)
}

func distance(one *pb.Point, two *pb.Point) (meter uint64) {
	earthRadius := 6378.1370
	d2r := math.Pi / 180
	dLong := (one.Longitude - two.Longitude) * d2r
	dLat := (one.Latitude - two.Latitude) * d2r
	a := math.Pow(math.Sin(dLat/2.0), 2.0) + math.Cos(one.Latitude*d2r)*math.Cos(two.Latitude*d2r)*math.Pow(math.Sin(dLong/2.0), 2.0)
	c := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1.0-a))
	d := earthRadius * c
	meter = uint64(d * 1000)
	return meter
}
