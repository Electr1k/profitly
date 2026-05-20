package model

import "time"

// Константы для статусов
const (
	OperStatusActive   = 1
	OperStatusInactive = 0

	CurrStatusActual    = 0
	CurrStatusNotActual = 1
)

// House - запись о доме (таблица AS_HOUSE)
type House struct {
	HouseGUID  string    `db:"house_guid" xml:"HOUSEGUID,attr"`
	AOGUID     string    `db:"aoguid" xml:"AOGUID,attr"`
	HouseNum   string    `db:"house_num" xml:"HOUSENUM,attr"`
	BuildNum   *string   `db:"build_num" xml:"BUILDNUM,attr"`
	StrucNum   *string   `db:"struc_num" xml:"STRUCNUM,attr"`
	PostalCode *string   `db:"postalcode" xml:"POSTALCODE,attr"`
	OKATO      *string   `db:"okato" xml:"OKATO,attr"`
	OKTMO      *string   `db:"oktmo" xml:"OKTMO,attr"`
	HouseType  int       `db:"house_type" xml:"HOUSETYPE,attr"`
	StartDate  time.Time `db:"start_date" xml:"STARTDATE,attr"`
	EndDate    time.Time `db:"end_date" xml:"ENDDATE,attr"`
	UpdateDate time.Time `db:"update_date" xml:"UPDATEDATE,attr"`
	Counter    int       `db:"counter" xml:"COUNTER,attr"`
	OperStatus int       `db:"oper_status" xml:"OPERSTATUS,attr"`
	CurrStatus int       `db:"curr_status" xml:"CURRSTATUS,attr"`
	DivType    *int      `db:"div_type" xml:"DIVTYPE,attr"`
	NorMDoc    *string   `db:"norm_doc" xml:"NORMDOC,attr"`
	TerIfNSI   *string   `db:"ter_if_nsi" xml:"TERIFNSI,attr"`
	IfNSI      *string   `db:"if_nsi" xml:"IFNSI,attr"`
	OKATOI     *string   `db:"okato_i" xml:"OKATOI,attr"`
	OKTMOI     *string   `db:"oktmo_i" xml:"OKTMOI,attr"`
}

// Room - помещение/квартира (таблица AS_ROOMS)
type Room struct {
	RoomGUID   string    `db:"room_guid" xml:"ROOMGUID,attr"`
	HouseGUID  string    `db:"house_guid" xml:"HOUSEGUID,attr"`
	RoomNum    string    `db:"room_num" xml:"ROOMNUM,attr"`
	RoomType   int       `db:"room_type" xml:"ROOMTYPE,attr"`
	RoomArea   *float64  `db:"room_area" xml:"ROOMAREA,attr"`
	PostalCode *string   `db:"postalcode" xml:"POSTALCODE,attr"`
	StartDate  time.Time `db:"start_date" xml:"STARTDATE,attr"`
	EndDate    time.Time `db:"end_date" xml:"ENDDATE,attr"`
	UpdateDate time.Time `db:"update_date" xml:"UPDATEDATE,attr"`
	OperStatus int       `db:"oper_status" xml:"OPERSTATUS,attr"`
	CurrStatus int       `db:"curr_status" xml:"CURRSTATUS,attr"`
}

// CarPlace - машино-место (таблица AS_CARPLACES)
type CarPlace struct {
	CarPlaceGUID string    `db:"carplace_guid" xml:"CARPLACEGUID,attr"`
	HouseGUID    string    `db:"house_guid" xml:"HOUSEGUID,attr"`
	CarPlaceNum  string    `db:"carplace_num" xml:"CARPLACENUM,attr"`
	CarPlaceType int       `db:"carplace_type" xml:"CARPLACETYPE,attr"`
	StartDate    time.Time `db:"start_date" xml:"STARTDATE,attr"`
	EndDate      time.Time `db:"end_date" xml:"ENDDATE,attr"`
	UpdateDate   time.Time `db:"update_date" xml:"UPDATEDATE,attr"`
	OperStatus   int       `db:"oper_status" xml:"OPERSTATUS,attr"`
	CurrStatus   int       `db:"curr_status" xml:"CURRSTATUS,attr"`
}
