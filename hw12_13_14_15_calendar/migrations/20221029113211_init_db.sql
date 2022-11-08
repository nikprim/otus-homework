-- +goose Up
create
extension if not exists "uuid-ossp";

CREATE TABLE events
(
    guid          uuid      NOT NULL,
    title         text      NOT NULL,
    start_at      timestamp NOT NULL,
    end_at        timestamp NOT NULL,
    description   text,
    user_guid     uuid      NOT NULL,
    notify_before interval second (0),
    CONSTRAINT unq_guid UNIQUE (guid),
    CONSTRAINT pk_event_guid PRIMARY KEY (guid)
);

-- +goose Down
drop
extension if exists "uuid-ossp";

DROP TABLE events;