
CREATE TABLE IF NOT EXISTS Users (
  ID INT PRIMARY KEY,
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

CREATE TABLE IF NOT EXISTS Bookings (
  ID INT PRIMARY KEY,
  User INT NOT NULL,
  LendingTime DATETIME NOT NULL,
  ReturnTime DATETIME NOT NULL,
  Teacher INT DEFAULT 0,
  Student INT DEFAULT 0,
  Chromebook INT DEFAULT 0,
  WAP INT DEFAULT 0,
  Projector INT DEFAULT 0,
  FOREIGN KEY (User) REFERENCES Users (ID)
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