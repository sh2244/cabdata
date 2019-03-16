package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"log"
	"strings"
	"time"
)

// connect to the MySQL database. Exit the program if connection fails
func connectDB() *sqlx.DB {
	// the dsn would be normally be provided externally eg via service discovery (consul, etcd)
	mydsn := "api:secret@tcp(127.0.0.1:3306)/ny_cab_data?parseTime=true"
	dbx, err := sqlx.Connect("mysql", mydsn)
	if err != nil {
		log.Fatal(err)
	}
	return dbx
}

// Trips implements the TripService interface
type Trips struct {
	dbx   *sqlx.DB     // connection to the mysql database
	cache *cache.Cache // connection to go-cache
}

var _ TripService = (*Trips)(nil) // assert that Trips implements the TripService interface

// MedallionsCount contains results of cache & database lookups
type MedallionsCount struct {
	Medallion string `db:"medallion" json:"medallion"`
	Count     int    `db:"mcount" json:"count"`
}

// NewTrips creates a new Trips struct
func NewTrips() Trips {
	dbx := connectDB()
	// default expiration duration and cleanup interval would be normally be provided externally eg via service discovery (consul, etcd)
	return Trips{dbx, cache.New(5*time.Minute, 10*time.Minute)}
}

// FlushCache flushes the cache
func (t Trips) FlushCache() {
	t.cache.Flush()
}

// CountByMedallions retrieves medallion counts from the cache. Any cache misses are retrieved directly
// from the database
func (t Trips) CountByMedallions(medallions []string, date string) (medallionCounts []MedallionsCount) {
	var bypassMedallions []string

	for _, medallion := range medallions {
		key := keyify(medallion, date)
		if count, found := t.cache.Get(key); found {
			if count, ok := count.(int); ok {
				medallionCounts = append(medallionCounts, MedallionsCount{medallion, count})
				continue
			}
		}
		bypassMedallions = append(bypassMedallions, medallion)
	}
	bypassCounts := t.CountByMedallionsBypass(bypassMedallions, date)
	return append(medallionCounts, bypassCounts...)
}

// CountByMedallionsBypass bypasses the cache to retrieve medallion counts directly from the database,
// then stores the counts in the cache
func (t Trips) CountByMedallionsBypass(medallions []string, date string) (medallionCounts []MedallionsCount) {
	if len(medallions) == 0 || len(date) == 0 {
		return []MedallionsCount{}
	}

	// sqlx.In() doco is unclear on how to combine IN and where fields; manually build IN statement
	sql := "select medallion, count(medallion) as mcount from cab_trip_data where medallion in ( "
	for i, medallion := range medallions {
		medallions[i] = fmt.Sprintf("'%s'", medallion)
	}
	sql += strings.Join(medallions, " , ")
	sql += " ) and pickup_datetime >= ? and pickup_datetime < date_add( ? , INTERVAL 1 DAY ) group by medallion"

	err := t.dbx.Select(&medallionCounts, sql, date, date)
	if err != nil {
		log.Fatal(err)
	}

	t.addToCache(medallionCounts, date)
	return medallionCounts
}

// addToCache adds fresh copies of medallian counts to the cache
func (t Trips) addToCache(counts []MedallionsCount, date string) {
	for _, row := range counts {
		key := keyify(row.Medallion, date)
		t.cache.Delete(key)
		if err := t.cache.Add(key, row.Count, cache.DefaultExpiration); err != nil {
			log.Fatal(err)
		}
	}
}

// keyify produces a key like "2B1A06E9228B7278227621EF1B879A1D_2013-12-01" for storing
// counts in the cache
func keyify(medallion string, date string) string {
	return fmt.Sprintf("%s_%s", medallion, date)
}
