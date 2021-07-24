-- +goose Up
DROP TABLE IF EXISTS Slots;
CREATE TABLE Slots (
            ID SERIAL PRIMARY KEY,
            Descr TEXT
);
INSERT INTO Slots (ID, Descr)
VALUES
    ('1', 'top slot'),
    ('2', 'side slot'),
    ('3', 'bottom slot');

DROP TABLE IF EXISTS Banners;
CREATE TABLE Banners (
            ID SERIAL PRIMARY KEY,
            Descr TEXT
);
INSERT INTO Banners (ID, Descr)
VALUES
    ('1', 'Cartoons'),
    ('2', 'Food'),
    ('3', 'Gardening');

DROP TABLE IF EXISTS BannersInSlots;
CREATE TABLE BannersInSlots (
    BannerID SERIAL REFERENCES Banners (ID),
    SlotID SERIAL REFERENCES Slots (ID),
    PRIMARY KEY (BannerID, SlotID)
);
INSERT INTO BannersInSlots (BannerID, SlotID)
VALUES
    ('1', '1'),
    ('1', '2'),
    ('1', '3'),
    ('2', '1'),
    ('2', '2'),
    ('2', '3'),
    ('3', '1'),
    ('3', '2'),
    ('3', '3');

DROP TABLE IF EXISTS SocGroups;
CREATE TABLE SocGroups (
            ID SERIAL PRIMARY KEY,
            Descr TEXT
);
INSERT INTO SocGroups (ID, Descr)
VALUES
    ('1', 'Young'),
    ('2', 'Old');

DROP TABLE IF EXISTS Statistic;
CREATE TABLE Statistic (
            SlotID SERIAL REFERENCES Slots (ID),
            BannerID SERIAL REFERENCES Banners (ID),
            SocGroupID SERIAL REFERENCES SocGroups (ID),
            PRIMARY KEY (SlotID, BannerID, SocGroupID),
            Impressions BIGINT,
            Clicks BIGINT
);
INSERT INTO Statistic (SlotID, BannerID, SocGroupID, Impressions, Clicks)
VALUES
    ('1', '1', '1', '10', '3'),
    ('2', '1', '1', '10', '3'),
    ('3', '1', '1', '10', '3'),
    ('1', '2', '1', '10', '2'),
    ('2', '2', '1', '10', '2'),
    ('3', '2', '1', '10', '2'),
    ('1', '3', '1', '10', '1'),
    ('2', '3', '1', '10', '1'),
    ('3', '3', '1', '10', '1'),
    ('1', '1', '2', '10', '1'),
    ('2', '1', '2', '10', '1'),
    ('3', '1', '2', '10', '1'),
    ('1', '2', '2', '10', '2'),
    ('2', '2', '2', '10', '2'),
    ('3', '2', '2', '10', '2'),
    ('1', '3', '2', '10', '3'),
    ('2', '3', '2', '10', '3'),
    ('3', '3', '2', '10', '3');

-- +goose Down
drop table Statistic;
drop table BannersInSlots;
drop table Slots;
drop table Banners;
drop table SocGroups;

