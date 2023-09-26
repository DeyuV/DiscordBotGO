CREATE TABLE "guilds"
(
    "id" varchar(70) primary key,
    "name" varchar(255) not null
);

CREATE TABLE "commands"
(
    "commandid" serial primary key,
    "name" varchar(50) not null,
    "description" varchar(255) not null,
    "defaultcommand" boolean not null
);

CREATE TABLE "guildcommands"
(
    "guildid" varchar(70) references guilds(id) on delete cascade,
    "commandid" int references commands(commandid) on delete cascade,
    primary key (guildid, commandid)
);

CREATE TABLE "guildlogsp"
(
    "id" varchar(255) primary key,
    "guildid" varchar(70),
    "map" varchar(255) not null,
    "spawntime" varchar(255) not null,
    "winningnation" varchar(255) not null,
    "userspawning" varchar(255) not null,
    "userinteracting" varchar(255) not null,
    "modified" varchar(255),
    "spdate" date not null,
    foreign key (guildid) references guilds(id) on delete cascade
);

CREATE TABLE "guildchannelsid"
(
    "guildid" varchar(70),
    "name" varchar(255) not null,
    "channelid" varchar(255) not null,
    primary key(guildid, name),
    foreign key (guildid) references guilds(id) on delete cascade
);

CREATE TABLE "guildmessagesid"
(
    "guildid" varchar(70),
    "name" varchar(255),
    "messageid" varchar(255),
    primary key(guildid, name),
    foreign key (guildid) references guilds(id) on delete cascade
);

INSERT INTO "commands"
    (name, description, defaultcommand)
VALUES
    ('setup-sp', 'Setup for all sp commands', true),
    ('set-sp-forum-channel', 'Set current channel to be used for forum sp notification', true);