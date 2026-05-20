CREATE TABLE IF NOT EXISTS houses (
    house_guid UUID PRIMARY KEY,
    aoguid UUID NOT NULL,
    house_num VARCHAR(20) NOT NULL,
    build_num VARCHAR(20),
    struc_num VARCHAR(20),
    postalcode VARCHAR(6),
    okato VARCHAR(11),
    oktmo VARCHAR(11),
    house_type SMALLINT,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    update_date TIMESTAMP,
    counter INTEGER,
    oper_status SMALLINT,
    curr_status SMALLINT,
    div_type SMALLINT,
    norm_doc VARCHAR(36),
    ter_if_nsi VARCHAR(4),
    if_nsi VARCHAR(4),
    okato_i VARCHAR(11),
    oktmo_i VARCHAR(11)
    );

COMMENT ON TABLE houses IS 'Сведения о домах, строениях, сооружениях (ФИАС, AS_HOUSE)';
COMMENT ON COLUMN houses.house_guid IS 'Уникальный идентификатор дома (GUID)';
COMMENT ON COLUMN houses.aoguid IS 'GUID адресного объекта (улица, город), к которому привязан дом';
COMMENT ON COLUMN houses.house_num IS 'Номер дома';
COMMENT ON COLUMN houses.build_num IS 'Номер корпуса';
COMMENT ON COLUMN houses.struc_num IS 'Номер строения';
COMMENT ON COLUMN houses.postalcode IS 'Почтовый индекс';
COMMENT ON COLUMN houses.okato IS 'Код ОКАТО';
COMMENT ON COLUMN houses.oktmo IS 'Код ОКТМО';
COMMENT ON COLUMN houses.house_type IS 'Тип дома (справочник 16)';
COMMENT ON COLUMN houses.start_date IS 'Начало действия записи';
COMMENT ON COLUMN houses.end_date IS 'Окончание действия записи';
COMMENT ON COLUMN houses.update_date IS 'Дата последнего обновления';
COMMENT ON COLUMN houses.counter IS 'Счётчик обновлений';
COMMENT ON COLUMN houses.oper_status IS 'Статус действия (0-не действует, 1-действует)';
COMMENT ON COLUMN houses.curr_status IS 'Статус актуальности (0-актуальный, 1-неактуальный)';
COMMENT ON COLUMN houses.div_type IS 'Тип деления (справочник 104)';
COMMENT ON COLUMN houses.norm_doc IS 'Внешний ключ на нормативный документ';
COMMENT ON COLUMN houses.ter_if_nsi IS 'Код территориального органа ИФНС';
COMMENT ON COLUMN houses.if_nsi IS 'Код ИФНС';
COMMENT ON COLUMN houses.okato_i IS 'Исторический ОКАТО';
COMMENT ON COLUMN houses.oktmo_i IS 'Исторический ОКТМО';

--CREATE INDEX IF NOT EXISTS idx_houses_house_guid ON houses(house_guid);
--CREATE INDEX IF NOT EXISTS idx_houses_aoguid ON houses(aoguid);

-- Таблица квартир/помещений
CREATE TABLE IF NOT EXISTS rooms (
    room_guid UUID PRIMARY KEY,
    house_guid UUID NOT NULL REFERENCES houses(house_guid) ON DELETE CASCADE,
    room_num VARCHAR(20) NOT NULL,
    room_type SMALLINT,
    room_area NUMERIC(15,2),
    postalcode VARCHAR(6),
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    update_date TIMESTAMP,
    oper_status SMALLINT,
    curr_status SMALLINT
    );

COMMENT ON TABLE rooms IS 'Сведения о комнатах, квартирах, офисах (AS_ROOMS)';
COMMENT ON COLUMN rooms.room_guid IS 'Уникальный идентификатор помещения';
COMMENT ON COLUMN rooms.house_guid IS 'Ссылка на дом (внешний ключ)';
COMMENT ON COLUMN rooms.room_num IS 'Номер помещения';
COMMENT ON COLUMN rooms.room_type IS 'Тип помещения (справочник 17)';
COMMENT ON COLUMN rooms.room_area IS 'Площадь помещения (кв. м)';
COMMENT ON COLUMN rooms.postalcode IS 'Почтовый индекс';
COMMENT ON COLUMN rooms.start_date IS 'Начало действия';
COMMENT ON COLUMN rooms.end_date IS 'Окончание действия';
COMMENT ON COLUMN rooms.update_date IS 'Дата обновления';
COMMENT ON COLUMN rooms.oper_status IS 'Статус действия';
COMMENT ON COLUMN rooms.curr_status IS 'Актуальность';

--CREATE INDEX IF NOT EXISTS idx_rooms_house_guid ON rooms(house_guid);
--CREATE INDEX IF NOT EXISTS idx_rooms_curr_status ON rooms(curr_status);

-- Таблица машино-мест
CREATE TABLE IF NOT EXISTS carplaces (
    carplace_guid UUID PRIMARY KEY,
    house_guid UUID NOT NULL REFERENCES houses(house_guid) ON DELETE CASCADE,
    carplace_num VARCHAR(20) NOT NULL,
    carplace_type SMALLINT,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    update_date TIMESTAMP,
    oper_status SMALLINT,
    curr_status SMALLINT
    );

COMMENT ON TABLE carplaces IS 'Сведения о машино-местах (AS_CARPLACES)';
COMMENT ON COLUMN carplaces.carplace_guid IS 'Уникальный идентификатор машино-места';
COMMENT ON COLUMN carplaces.house_guid IS 'Ссылка на дом';
COMMENT ON COLUMN carplaces.carplace_num IS 'Номер машино-места';
COMMENT ON COLUMN carplaces.carplace_type IS 'Тип машино-места (справочник 52)';
COMMENT ON COLUMN carplaces.start_date IS 'Начало действия';
COMMENT ON COLUMN carplaces.end_date IS 'Окончание действия';
COMMENT ON COLUMN carplaces.update_date IS 'Дата обновления';
COMMENT ON COLUMN carplaces.oper_status IS 'Статус действия';
COMMENT ON COLUMN carplaces.curr_status IS 'Актуальность';

--CREATE INDEX IF NOT EXISTS idx_carplaces_house_guid ON carplaces(house_guid);
--CREATE INDEX IF NOT EXISTS idx_carplaces_curr_status ON carplaces(curr_status);
