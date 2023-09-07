create extension if not exists "uuid-ossp";

create table if not exists compliance_check_requests
(
    id                    uuid primary key default uuid_generate_v4(),
    pwg_entity_guid       varchar(255)             not null,
    pwg_entity_type       varchar(64)              not null,
    request_external_guid varchar(255)             not null,
    provider              varchar(64)              not null,
    check_rules           jsonb,
    raw_request           jsonb,
    status                varchar(64)              not null,
    requested_at          timestamp with time zone not null,
    finished_at           timestamp,
    created_at            timestamp with time zone not null,
    updated_at            timestamp with time zone not null
);

create table if not exists compliance_checks
(
    id                          uuid primary key default uuid_generate_v4(),
    external_guid               varchar(255)             not null unique,
    compliance_check_request_id uuid
        constraint fk_compliance_checks_compliance_check_request
            references compliance_check_requests (id),
    provider                    varchar(64)              not null,
    status                      varchar(64)              not null,
    passed_at                   timestamp with time zone,
    expired_at                  timestamp,
    created_at                  timestamp with time zone not null,
    updated_at                  timestamp with time zone not null
);
