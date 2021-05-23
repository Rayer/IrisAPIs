create table pbs_traffic_events
(
    region          text        not null,
    source          text        null,
    area            text        null,
    uid             varchar(16) not null,
    direction       varchar(12) null,
    longitude       float       null,
    latitude        float       null,
    entry_timestamp datetime    null,
    constraint pbs_traffic_events_uid_uindex
        unique (uid)
);

alter table pbs_traffic_events
    add primary key (uid);

create table pbs_traffic_history
(
    id               int auto_increment,
    uid              varchar(16) null,
    update_timestamp datetime    not null,
    information      longtext    not null,
    constraint pbs_traffic_history_id_uindex
        unique (id),
    constraint pbs_traffic_history_pbs_traffic_events_uid_fk
        foreign key (uid) references pbs_traffic_events (uid)
            on delete cascade
);

alter table pbs_traffic_history
    add primary key (id);
