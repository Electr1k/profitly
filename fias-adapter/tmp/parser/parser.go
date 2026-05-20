package parser

import (
	"archive/zip"
	"context"
	"encoding/xml"
	"fias-adapter/tmp/model"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const batchSize = 1000

func ParseArchive(db *pgxpool.Pool, archivePath string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}
	defer r.Close()

	for _, file := range r.File {
		name := file.Name
		switch {
		case strings.Contains(name, "AS_HOUSE"):
			log.Printf("Parsing houses from %s", name)
			if err := parseHouses(db, file); err != nil {
				return fmt.Errorf("parse houses: %w", err)
			}
		case strings.Contains(name, "AS_ROOMS"):
			log.Printf("Parsing rooms from %s", name)
			if err := parseRooms(db, file); err != nil {
				return fmt.Errorf("parse rooms: %w", err)
			}
		case strings.Contains(name, "AS_CARPLACES"):
			log.Printf("Parsing carplaces from %s", name)
			if err := parseCarplaces(db, file); err != nil {
				return fmt.Errorf("parse carplaces: %w", err)
			}
		}
	}
	return nil
}

func parseHouses(db *pgxpool.Pool, zipFile *zip.File) error {
	rc, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	decoder := xml.NewDecoder(rc)
	batch := &pgx.Batch{}
	total := 0

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "House" {
				var house model.House
				if err := decoder.DecodeElement(&house, &se); err != nil {
					log.Printf("Warning: decode house failed: %v", err)
					continue
				}
				// фильтр: только актуальные и действующие
				if house.CurrStatus != model.CurrStatusActual || house.OperStatus != model.OperStatusActive {
					continue
				}

				batch.Queue(`
                    INSERT INTO houses (
                        house_guid, aoguid, house_num, build_num, struc_num,
                        postalcode, okato, oktmo, house_type, start_date, end_date,
                        update_date, counter, oper_status, curr_status, div_type,
                        norm_doc, ter_if_nsi, if_nsi, okato_i, oktmo_i
                    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
                    ON CONFLICT (house_guid) DO NOTHING
                `,
					house.HouseGUID, house.AOGUID, house.HouseNum, house.BuildNum,
					house.StrucNum, house.PostalCode, house.OKATO, house.OKTMO,
					house.HouseType, house.StartDate, house.EndDate, house.UpdateDate,
					house.Counter, house.OperStatus, house.CurrStatus, house.DivType,
					house.NorMDoc, house.TerIfNSI, house.IfNSI, house.OKATOI, house.OKTMOI,
				)
				total++

				if total%batchSize == 0 {
					if err := executeBatch(db, batch); err != nil {
						return err
					}
					batch = &pgx.Batch{}
					log.Printf("Houses inserted: %d", total)
				}
			}
		}
	}

	if batch.Len() > 0 {
		if err := executeBatch(db, batch); err != nil {
			return err
		}
	}
	log.Printf("Finished houses: total %d records", total)
	return nil
}

func parseRooms(db *pgxpool.Pool, zipFile *zip.File) error {
	rc, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	decoder := xml.NewDecoder(rc)
	batch := &pgx.Batch{}
	total := 0

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "Room" {
				var room model.Room
				if err := decoder.DecodeElement(&room, &se); err != nil {
					log.Printf("Warning: decode room failed: %v", err)
					continue
				}
				if room.CurrStatus != model.CurrStatusActual || room.OperStatus != model.OperStatusActive {
					continue
				}

				batch.Queue(`
                    INSERT INTO rooms (
                        room_guid, house_guid, room_num, room_type, room_area,
                        postalcode, start_date, end_date, update_date, oper_status, curr_status
                    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
                    ON CONFLICT (room_guid) DO NOTHING
                `,
					room.RoomGUID, room.HouseGUID, room.RoomNum, room.RoomType,
					room.RoomArea, room.PostalCode, room.StartDate, room.EndDate,
					room.UpdateDate, room.OperStatus, room.CurrStatus,
				)
				total++

				if total%batchSize == 0 {
					if err := executeBatch(db, batch); err != nil {
						return err
					}
					batch = &pgx.Batch{}
					log.Printf("Rooms inserted: %d", total)
				}
			}
		}
	}

	if batch.Len() > 0 {
		if err := executeBatch(db, batch); err != nil {
			return err
		}
	}
	log.Printf("Finished rooms: total %d records", total)
	return nil
}

func parseCarplaces(db *pgxpool.Pool, zipFile *zip.File) error {
	rc, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	decoder := xml.NewDecoder(rc)
	batch := &pgx.Batch{}
	total := 0

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "Carplace" {
				var cp model.CarPlace
				if err := decoder.DecodeElement(&cp, &se); err != nil {
					log.Printf("Warning: decode carplace failed: %v", err)
					continue
				}
				if cp.CurrStatus != model.CurrStatusActual || cp.OperStatus != model.OperStatusActive {
					continue
				}

				batch.Queue(`
                    INSERT INTO carplaces (
                        carplace_guid, house_guid, carplace_num, carplace_type,
                        start_date, end_date, update_date, oper_status, curr_status
                    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
                    ON CONFLICT (carplace_guid) DO NOTHING
                `,
					cp.CarPlaceGUID, cp.HouseGUID, cp.CarPlaceNum, cp.CarPlaceType,
					cp.StartDate, cp.EndDate, cp.UpdateDate, cp.OperStatus, cp.CurrStatus,
				)
				total++

				if total%batchSize == 0 {
					if err := executeBatch(db, batch); err != nil {
						return err
					}
					batch = &pgx.Batch{}
					log.Printf("Carplaces inserted: %d", total)
				}
			}
		}
	}

	if batch.Len() > 0 {
		if err := executeBatch(db, batch); err != nil {
			return err
		}
	}
	log.Printf("Finished carplaces: total %d records", total)
	return nil
}

func executeBatch(db *pgxpool.Pool, batch *pgx.Batch) error {
	ctx := context.Background()
	br := db.SendBatch(ctx, batch)
	defer br.Close()

	for range batch.QueuedQueries {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("batch exec: %w", err)
		}
	}
	return nil
}
