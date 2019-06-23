CREATE TABLE IF NOT EXISTS Users (
  ID INT PRIMARY KEY AUTO_INCREMENT,
  Email VARCHAR(64) UNIQUE NOT NULL,
  Name VARCHAR(32) NOT NULL,
  Type ENUM('Admin', 'Teacher') DEFAULT 'Teacher',
  Password BLOB NOT NULL
);

CREATE TABLE IF NOT EXISTS Devices (
  ID VARCHAR(8) PRIMARY KEY,
  Type ENUM('Student-iPad', 'Teacher-iPad', 'Chromebook', 'WAP', 'WirelessProjector') NOT NULL,
  JoinDate DATETIME
);

INSERT INTO `Devices` VALUES ('AP01','WAP','2019-04-20 11:01:37'),('AP02','WAP','2019-04-20 11:01:37'),('AP03','WAP','2019-04-20 11:01:37'),('C001','Chromebook','2019-04-20 11:01:36'),('C002','Chromebook','2019-04-20 11:01:36'),('C003','Chromebook','2019-04-20 11:01:36'),('C004','Chromebook','2019-04-20 11:01:36'),('C005','Chromebook','2019-04-20 11:01:36'),('C006','Chromebook','2019-04-20 11:01:36'),('C007','Chromebook','2019-04-20 11:01:36'),('C008','Chromebook','2019-04-20 11:01:36'),('C009','Chromebook','2019-04-20 11:01:36'),('C010','Chromebook','2019-04-20 11:01:36'),('C011','Chromebook','2019-04-20 11:01:36'),('C012','Chromebook','2019-04-20 11:01:36'),('C013','Chromebook','2019-04-20 11:01:36'),('C014','Chromebook','2019-04-20 11:01:36'),('C015','Chromebook','2019-04-20 11:01:36'),('C016','Chromebook','2019-04-20 11:01:36'),('C017','Chromebook','2019-04-20 11:01:36'),('C018','Chromebook','2019-04-20 11:01:36'),('C019','Chromebook','2019-04-20 11:01:36'),('C020','Chromebook','2019-04-20 11:01:36'),('C021','Chromebook','2019-04-20 11:01:36'),('C022','Chromebook','2019-04-20 11:01:36'),('C023','Chromebook','2019-04-20 11:01:36'),('C024','Chromebook','2019-04-20 11:01:36'),('C025','Chromebook','2019-04-20 11:01:36'),('C026','Chromebook','2019-04-20 11:01:36'),('C027','Chromebook','2019-04-20 11:01:36'),('C028','Chromebook','2019-04-20 11:01:36'),('C029','Chromebook','2019-04-20 11:01:36'),('C030','Chromebook','2019-04-20 11:01:36'),('C031','Chromebook','2019-04-20 11:01:37'),('C032','Chromebook','2019-04-20 11:01:37'),('C033','Chromebook','2019-04-20 11:01:37'),('C034','Chromebook','2019-04-20 11:01:37'),('C035','Chromebook','2019-04-20 11:01:37'),('C036','Chromebook','2019-04-20 11:01:37'),('C037','Chromebook','2019-04-20 11:01:37'),('C038','Chromebook','2019-04-20 11:01:37'),('C039','Chromebook','2019-04-20 11:01:37'),('C040','Chromebook','2019-04-20 11:01:37'),('C041','Chromebook','2019-04-20 11:01:37'),('C042','Chromebook','2019-04-20 11:01:37'),('C043','Chromebook','2019-04-20 11:01:37'),('C044','Chromebook','2019-04-20 11:01:37'),('C045','Chromebook','2019-04-20 11:01:37'),('C046','Chromebook','2019-04-20 11:01:37'),('C047','Chromebook','2019-04-20 11:01:37'),('C048','Chromebook','2019-04-20 11:01:37'),('C049','Chromebook','2019-04-20 11:01:37'),('C050','Chromebook','2019-04-20 11:01:37'),('C051','Chromebook','2019-04-20 11:01:37'),('C052','Chromebook','2019-04-20 11:01:37'),('C053','Chromebook','2019-04-20 11:01:37'),('C054','Chromebook','2019-04-20 11:01:37'),('C055','Chromebook','2019-04-20 11:01:37'),('C056','Chromebook','2019-04-20 11:01:37'),('C057','Chromebook','2019-04-20 11:01:37'),('C058','Chromebook','2019-04-20 11:01:37'),('C059','Chromebook','2019-04-20 11:01:37'),('C060','Chromebook','2019-04-20 11:01:37'),('C061','Chromebook','2019-04-20 11:01:37'),('C062','Chromebook','2019-04-20 11:01:37'),('C063','Chromebook','2019-04-20 11:01:37'),('C064','Chromebook','2019-04-20 11:01:37'),('C065','Chromebook','2019-04-20 11:01:37'),('C066','Chromebook','2019-04-20 11:01:37'),('C067','Chromebook','2019-04-20 11:01:37'),('C068','Chromebook','2019-04-20 11:01:37'),('C069','Chromebook','2019-04-20 11:01:37'),('C070','Chromebook','2019-04-20 11:01:37'),('C071','Chromebook','2019-04-20 11:01:37'),('C072','Chromebook','2019-04-20 11:01:37'),('C073','Chromebook','2019-04-20 11:01:37'),('C074','Chromebook','2019-04-20 11:01:37'),('C075','Chromebook','2019-04-20 11:01:37'),('C076','Chromebook','2019-04-20 11:01:37'),('C077','Chromebook','2019-04-20 11:01:37'),('C078','Chromebook','2019-04-20 11:01:37'),('C079','Chromebook','2019-04-20 11:01:37'),('C080','Chromebook','2019-04-20 11:01:37'),('C081','Chromebook','2019-04-20 11:01:37'),('C082','Chromebook','2019-04-20 11:01:37'),('C083','Chromebook','2019-04-20 11:01:37'),('C084','Chromebook','2019-04-20 11:01:37'),('C085','Chromebook','2019-04-20 11:01:37'),('C086','Chromebook','2019-04-20 11:01:37'),('C087','Chromebook','2019-04-20 11:01:37'),('C088','Chromebook','2019-04-20 11:01:37'),('C089','Chromebook','2019-04-20 11:01:37'),('C090','Chromebook','2019-04-20 11:01:37'),('CB001','Chromebook','2019-04-20 11:01:36'),('CB002','Chromebook','2019-04-20 11:01:36'),('CB003','Chromebook','2019-04-20 11:01:36'),('CB004','Chromebook','2019-04-20 11:01:36'),('CB005','Chromebook','2019-04-20 11:01:36'),('CB006','Chromebook','2019-04-20 11:01:36'),('CB007','Chromebook','2019-04-20 11:01:36'),('CB008','Chromebook','2019-04-20 11:01:36'),('CB009','Chromebook','2019-04-20 11:01:36'),('CB010','Chromebook','2019-04-20 11:01:36'),('CB011','Chromebook','2019-04-20 11:01:36'),('CB012','Chromebook','2019-04-20 11:01:36'),('CB013','Chromebook','2019-04-20 11:01:36'),('CB014','Chromebook','2019-04-20 11:01:36'),('CB015','Chromebook','2019-04-20 11:01:36'),('CB016','Chromebook','2019-04-20 11:01:36'),('CB017','Chromebook','2019-04-20 11:01:36'),('CB018','Chromebook','2019-04-20 11:01:36'),('CB019','Chromebook','2019-04-20 11:01:36'),('CB020','Chromebook','2019-04-20 11:01:36'),('CB021','Chromebook','2019-04-20 11:01:36'),('CB022','Chromebook','2019-04-20 11:01:36'),('CB023','Chromebook','2019-04-20 11:01:36'),('CB024','Chromebook','2019-04-20 11:01:36'),('CB025','Chromebook','2019-04-20 11:01:36'),('CB026','Chromebook','2019-04-20 11:01:36'),('CB027','Chromebook','2019-04-20 11:01:36'),('CB028','Chromebook','2019-04-20 11:01:36'),('CB029','Chromebook','2019-04-20 11:01:36'),('CB030','Chromebook','2019-04-20 11:01:36'),('CB031','Chromebook','2019-04-20 11:01:36'),('CB032','Chromebook','2019-04-20 11:01:36'),('EZ01','WirelessProjector','2019-04-20 11:01:37'),('EZ02','WirelessProjector','2019-04-20 11:01:37'),('EZ03','WirelessProjector','2019-04-20 11:01:37'),('ST001','Student-iPad','2019-04-20 11:01:36'),('ST002','Student-iPad','2019-04-20 11:01:36'),('ST003','Student-iPad','2019-04-20 11:01:36'),('ST004','Student-iPad','2019-04-20 11:01:36'),('ST005','Student-iPad','2019-04-20 11:01:36'),('ST006','Student-iPad','2019-04-20 11:01:36'),('ST007','Student-iPad','2019-04-20 11:01:36'),('ST008','Student-iPad','2019-04-20 11:01:36'),('ST009','Student-iPad','2019-04-20 11:01:36'),('ST010','Student-iPad','2019-04-20 11:01:36'),('ST011','Student-iPad','2019-04-20 11:01:36'),('ST012','Student-iPad','2019-04-20 11:01:36'),('ST013','Student-iPad','2019-04-20 11:01:36'),('ST014','Student-iPad','2019-04-20 11:01:36'),('ST015','Student-iPad','2019-04-20 11:01:36'),('ST016','Student-iPad','2019-04-20 11:01:36'),('ST017','Student-iPad','2019-04-20 11:01:36'),('ST018','Student-iPad','2019-04-20 11:01:36'),('ST019','Student-iPad','2019-04-20 11:01:36'),('ST020','Student-iPad','2019-04-20 11:01:36'),('ST021','Student-iPad','2019-04-20 11:01:36'),('ST022','Student-iPad','2019-04-20 11:01:36'),('ST023','Student-iPad','2019-04-20 11:01:36'),('ST024','Student-iPad','2019-04-20 11:01:36'),('ST025','Student-iPad','2019-04-20 11:01:36'),('ST026','Student-iPad','2019-04-20 11:01:36'),('ST027','Student-iPad','2019-04-20 11:01:36'),('ST028','Student-iPad','2019-04-20 11:01:36'),('ST029','Student-iPad','2019-04-20 11:01:36'),('ST030','Student-iPad','2019-04-20 11:01:36'),('ST031','Student-iPad','2019-04-20 11:01:36'),('ST032','Student-iPad','2019-04-20 11:01:36'),('ST033','Student-iPad','2019-04-20 11:01:36'),('ST034','Student-iPad','2019-04-20 11:01:36'),('ST035','Student-iPad','2019-04-20 11:01:36'),('ST036','Student-iPad','2019-04-20 11:01:36'),('ST037','Student-iPad','2019-04-20 11:01:36'),('ST038','Student-iPad','2019-04-20 11:01:36'),('ST039','Student-iPad','2019-04-20 11:01:36'),('ST040','Student-iPad','2019-04-20 11:01:36'),('ST041','Student-iPad','2019-04-20 11:01:36'),('ST042','Student-iPad','2019-04-20 11:01:36'),('ST043','Student-iPad','2019-04-20 11:01:36'),('ST044','Student-iPad','2019-04-20 11:01:36'),('ST045','Student-iPad','2019-04-20 11:01:36'),('ST046','Student-iPad','2019-04-20 11:01:36'),('ST047','Student-iPad','2019-04-20 11:01:36'),('ST048','Student-iPad','2019-04-20 11:01:36'),('ST049','Student-iPad','2019-04-20 11:01:36'),('ST050','Student-iPad','2019-04-20 11:01:36'),('ST051','Student-iPad','2019-04-20 11:01:36'),('ST052','Student-iPad','2019-04-20 11:01:36'),('ST053','Student-iPad','2019-04-20 11:01:36'),('ST054','Student-iPad','2019-04-20 11:01:36'),('ST055','Student-iPad','2019-04-20 11:01:36'),('ST056','Student-iPad','2019-04-20 11:01:36'),('ST057','Student-iPad','2019-04-20 11:01:36'),('ST058','Student-iPad','2019-04-20 11:01:36'),('ST059','Student-iPad','2019-04-20 11:01:36'),('ST060','Student-iPad','2019-04-20 11:01:36'),('ST061','Student-iPad','2019-04-20 11:01:36'),('ST062','Student-iPad','2019-04-20 11:01:36'),('ST063','Student-iPad','2019-04-20 11:01:36'),('ST064','Student-iPad','2019-04-20 11:01:36'),('ST065','Student-iPad','2019-04-20 11:01:36'),('ST066','Student-iPad','2019-04-20 11:01:36'),('ST067','Student-iPad','2019-04-20 11:01:36'),('ST068','Student-iPad','2019-04-20 11:01:36'),('ST069','Student-iPad','2019-04-20 11:01:36'),('ST070','Student-iPad','2019-04-20 11:01:36'),('T01','Teacher-iPad','2019-04-20 11:01:36'),('T02','Teacher-iPad','2019-04-20 11:01:36'),('T03','Teacher-iPad','2019-04-20 11:01:36'),('T04','Teacher-iPad','2019-04-20 11:01:36'),('T05','Teacher-iPad','2019-04-20 11:01:36'),('T06','Teacher-iPad','2019-04-20 11:01:36'),('T07','Teacher-iPad','2019-04-20 11:01:36'),('T08','Teacher-iPad','2019-04-20 11:01:36'),('T09','Teacher-iPad','2019-04-20 11:01:36'),('T10','Teacher-iPad','2019-04-20 11:01:36'),('T11','Teacher-iPad','2019-04-20 11:01:36'),('T12','Teacher-iPad','2019-04-20 11:01:36'),('T13','Teacher-iPad','2019-04-20 11:01:36');

CREATE TABLE IF NOT EXISTS Bookings (
  ID INT PRIMARY KEY AUTO_INCREMENT,
  User INT NOT NULL,
  LendingTime DATETIME NOT NULL,
  ReturnTime DATETIME NOT NULL,
  Done BOOL DEFAULT false,
  FOREIGN KEY (User) REFERENCES Users (ID)
);

CREATE TABLE IF NOT EXISTS BookingDevices (
  BID INT,
  Type ENUM('Student-iPad', 'Teacher-iPad', 'Chromebook', 'WAP', 'WirelessProjector'),
  Amount INT DEFAULT 0,
  FOREIGN KEY (BID) REFERENCES Bookings (ID),
  PRIMARY KEY (BID, Type)
);

CREATE TABLE IF NOT EXISTS Records (
  Booking INT,
  Device VARCHAR(8),
  LentFrom DATETIME NOT NULL,
  LentUntil DATETIME,
  PRIMARY KEY (Booking, Device),
  FOREIGN KEY (Booking) REFERENCES Bookings (ID),
  FOREIGN KEY (Device) REFERENCES Devices (ID)
);

CREATE TABLE IF NOT EXISTS Sessions (
	ID CHAR(36),
	User INT NOT NULL,
	LastUsed TIMESTAMP NOT NULL,
	PRIMARY KEY (ID),
	FOREIGN KEY (User) REFERENCES Users (ID)
);

CREATE OR REPLACE VIEW UnDoneRecords AS 
SELECT Booking, Device
FROM Records
WHERE LentUntil IS NULL;

CREATE OR REPLACE VIEW UnDoneBookings AS 
SELECT ID, COUNT(Device) AS Amount
FROM Bookings B
LEFT JOIN UnDoneRecords R
ON B.ID = R.Booking 
GROUP BY ID;

CREATE OR REPLACE VIEW DevicesStatus AS 
SELECT D.ID, COUNT(Device) AS Status, Name, D.Type
FROM Devices D
LEFT JOIN UnDoneRecords R
ON D.ID = Device
LEFT JOIN Bookings B
ON B.ID = Booking 
LEFT JOIN Users U
ON U.ID = User 
GROUP BY D.ID;