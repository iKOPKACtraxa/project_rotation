-- +goose Up

CREATE TABLE Slots (
            ID SERIAL PRIMARY KEY,
            Descr TEXT
);

INSERT INTO Slots (ID, Descr)
VALUES
    ('1', 'top slot'),
    ('2', 'side slot'),
    ('3', 'bottom slot');

CREATE TABLE Banners (
            ID SERIAL PRIMARY KEY,
            Descr TEXT
);

CREATE TABLE BannersInSlots (
    BannerID SERIAL REFERENCES Banners (ID),
    SlotID SERIAL REFERENCES Slots (ID),
    PRIMARY KEY (BannerID, SlotID)
);

CREATE TABLE SocGroups (
            ID SERIAL PRIMARY KEY,
            Descr TEXT
);
INSERT INTO SocGroups (ID, Descr)
VALUES
    ('1', 'Young'),
    ('2', 'Old');

CREATE TABLE Statistic (
            SlotID SERIAL REFERENCES Slots (ID),
            BannerID SERIAL REFERENCES Banners (ID),
            SocGroupID SERIAL REFERENCES SocGroups (ID),
            PRIMARY KEY (SlotID, BannerID, SocGroupID),
            Impressions BIGINT,
            Clicks BIGINT
);

-- +goose Down
drop table Statistic;
drop table BannersInSlots;
drop table Slots;
drop table Banners;
drop table SocGroups;

