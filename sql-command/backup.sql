-- MySQL dump 10.16  Distrib 10.1.38-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: DeviceBooking
-- ------------------------------------------------------
-- Server version	10.1.38-MariaDB-0ubuntu0.18.04.2

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `Bookings`
--

DROP TABLE IF EXISTS `Bookings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Bookings` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `User` int(11) NOT NULL,
  `LendingTime` datetime NOT NULL,
  `ReturnTime` datetime NOT NULL,
  `Teacher` int(11) DEFAULT '0',
  `Student` int(11) DEFAULT '0',
  `Chromebook` int(11) DEFAULT '0',
  `WAP` int(11) DEFAULT '0',
  `Projector` int(11) DEFAULT '0',
  PRIMARY KEY (`ID`),
  KEY `User` (`User`),
  CONSTRAINT `Bookings_ibfk_1` FOREIGN KEY (`User`) REFERENCES `Users` (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Bookings`
--

LOCK TABLES `Bookings` WRITE;
/*!40000 ALTER TABLE `Bookings` DISABLE KEYS */;
INSERT INTO `Bookings` VALUES (1,1,'2019-05-13 09:15:00','2019-05-13 10:10:00',1,30,4,1,1),(2,1,'2019-05-14 12:30:00','2019-05-14 13:00:00',1,10,0,0,0),(3,1,'2019-06-10 13:05:00','2019-06-10 13:50:00',1,30,10,0,0),(4,1,'2019-05-16 10:10:00','2019-05-16 10:55:00',0,20,20,1,0),(5,1,'2019-05-15 14:00:00','2019-05-15 14:45:00',0,5,15,1,0),(6,1,'2019-05-13 11:05:00','2019-05-13 11:50:00',1,20,13,2,1),(7,2,'2019-05-23 13:05:00','2019-05-23 13:50:00',1,30,15,1,1),(8,2,'2019-05-11 07:30:00','2019-05-11 08:10:00',1,50,100,0,0),(9,2,'2019-05-11 07:30:00','2019-05-11 08:10:00',0,20,22,0,0),(10,2,'2019-05-11 07:30:00','2019-05-11 08:10:00',2,1,3,1,0),(11,2,'2019-05-23 10:10:00','2019-05-23 10:55:00',0,50,0,0,0),(12,2,'2019-05-23 10:10:00','2019-05-23 10:55:00',0,20,0,0,0),(13,1,'2019-05-22 09:15:00','2019-05-22 10:10:00',3,50,100,0,0),(14,1,'2019-05-22 09:15:00','2019-05-22 10:10:00',0,20,22,0,0),(15,1,'2019-05-22 10:10:00','2019-05-22 10:55:00',10,50,28,1,1);
/*!40000 ALTER TABLE `Bookings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Devices`
--

DROP TABLE IF EXISTS `Devices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Devices` (
  `ID` varchar(8) NOT NULL,
  `Type` enum('Student-iPad','Teacher-iPad','Chromebook','WAP','WirelessProjector') NOT NULL,
  `JoinDate` datetime DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Devices`
--

LOCK TABLES `Devices` WRITE;
/*!40000 ALTER TABLE `Devices` DISABLE KEYS */;
INSERT INTO `Devices` VALUES ('AP01','WAP','2019-04-20 11:01:37'),('AP02','WAP','2019-04-20 11:01:37'),('AP03','WAP','2019-04-20 11:01:37'),('C001','Chromebook','2019-04-20 11:01:36'),('C002','Chromebook','2019-04-20 11:01:36'),('C003','Chromebook','2019-04-20 11:01:36'),('C004','Chromebook','2019-04-20 11:01:36'),('C005','Chromebook','2019-04-20 11:01:36'),('C006','Chromebook','2019-04-20 11:01:36'),('C007','Chromebook','2019-04-20 11:01:36'),('C008','Chromebook','2019-04-20 11:01:36'),('C009','Chromebook','2019-04-20 11:01:36'),('C010','Chromebook','2019-04-20 11:01:36'),('C011','Chromebook','2019-04-20 11:01:36'),('C012','Chromebook','2019-04-20 11:01:36'),('C013','Chromebook','2019-04-20 11:01:36'),('C014','Chromebook','2019-04-20 11:01:36'),('C015','Chromebook','2019-04-20 11:01:36'),('C016','Chromebook','2019-04-20 11:01:36'),('C017','Chromebook','2019-04-20 11:01:36'),('C018','Chromebook','2019-04-20 11:01:36'),('C019','Chromebook','2019-04-20 11:01:36'),('C020','Chromebook','2019-04-20 11:01:36'),('C021','Chromebook','2019-04-20 11:01:36'),('C022','Chromebook','2019-04-20 11:01:36'),('C023','Chromebook','2019-04-20 11:01:36'),('C024','Chromebook','2019-04-20 11:01:36'),('C025','Chromebook','2019-04-20 11:01:36'),('C026','Chromebook','2019-04-20 11:01:36'),('C027','Chromebook','2019-04-20 11:01:36'),('C028','Chromebook','2019-04-20 11:01:36'),('C029','Chromebook','2019-04-20 11:01:36'),('C030','Chromebook','2019-04-20 11:01:36'),('C031','Chromebook','2019-04-20 11:01:37'),('C032','Chromebook','2019-04-20 11:01:37'),('C033','Chromebook','2019-04-20 11:01:37'),('C034','Chromebook','2019-04-20 11:01:37'),('C035','Chromebook','2019-04-20 11:01:37'),('C036','Chromebook','2019-04-20 11:01:37'),('C037','Chromebook','2019-04-20 11:01:37'),('C038','Chromebook','2019-04-20 11:01:37'),('C039','Chromebook','2019-04-20 11:01:37'),('C040','Chromebook','2019-04-20 11:01:37'),('C041','Chromebook','2019-04-20 11:01:37'),('C042','Chromebook','2019-04-20 11:01:37'),('C043','Chromebook','2019-04-20 11:01:37'),('C044','Chromebook','2019-04-20 11:01:37'),('C045','Chromebook','2019-04-20 11:01:37'),('C046','Chromebook','2019-04-20 11:01:37'),('C047','Chromebook','2019-04-20 11:01:37'),('C048','Chromebook','2019-04-20 11:01:37'),('C049','Chromebook','2019-04-20 11:01:37'),('C050','Chromebook','2019-04-20 11:01:37'),('C051','Chromebook','2019-04-20 11:01:37'),('C052','Chromebook','2019-04-20 11:01:37'),('C053','Chromebook','2019-04-20 11:01:37'),('C054','Chromebook','2019-04-20 11:01:37'),('C055','Chromebook','2019-04-20 11:01:37'),('C056','Chromebook','2019-04-20 11:01:37'),('C057','Chromebook','2019-04-20 11:01:37'),('C058','Chromebook','2019-04-20 11:01:37'),('C059','Chromebook','2019-04-20 11:01:37'),('C060','Chromebook','2019-04-20 11:01:37'),('C061','Chromebook','2019-04-20 11:01:37'),('C062','Chromebook','2019-04-20 11:01:37'),('C063','Chromebook','2019-04-20 11:01:37'),('C064','Chromebook','2019-04-20 11:01:37'),('C065','Chromebook','2019-04-20 11:01:37'),('C066','Chromebook','2019-04-20 11:01:37'),('C067','Chromebook','2019-04-20 11:01:37'),('C068','Chromebook','2019-04-20 11:01:37'),('C069','Chromebook','2019-04-20 11:01:37'),('C070','Chromebook','2019-04-20 11:01:37'),('C071','Chromebook','2019-04-20 11:01:37'),('C072','Chromebook','2019-04-20 11:01:37'),('C073','Chromebook','2019-04-20 11:01:37'),('C074','Chromebook','2019-04-20 11:01:37'),('C075','Chromebook','2019-04-20 11:01:37'),('C076','Chromebook','2019-04-20 11:01:37'),('C077','Chromebook','2019-04-20 11:01:37'),('C078','Chromebook','2019-04-20 11:01:37'),('C079','Chromebook','2019-04-20 11:01:37'),('C080','Chromebook','2019-04-20 11:01:37'),('C081','Chromebook','2019-04-20 11:01:37'),('C082','Chromebook','2019-04-20 11:01:37'),('C083','Chromebook','2019-04-20 11:01:37'),('C084','Chromebook','2019-04-20 11:01:37'),('C085','Chromebook','2019-04-20 11:01:37'),('C086','Chromebook','2019-04-20 11:01:37'),('C087','Chromebook','2019-04-20 11:01:37'),('C088','Chromebook','2019-04-20 11:01:37'),('C089','Chromebook','2019-04-20 11:01:37'),('C090','Chromebook','2019-04-20 11:01:37'),('CB001','Chromebook','2019-04-20 11:01:36'),('CB002','Chromebook','2019-04-20 11:01:36'),('CB003','Chromebook','2019-04-20 11:01:36'),('CB004','Chromebook','2019-04-20 11:01:36'),('CB005','Chromebook','2019-04-20 11:01:36'),('CB006','Chromebook','2019-04-20 11:01:36'),('CB007','Chromebook','2019-04-20 11:01:36'),('CB008','Chromebook','2019-04-20 11:01:36'),('CB009','Chromebook','2019-04-20 11:01:36'),('CB010','Chromebook','2019-04-20 11:01:36'),('CB011','Chromebook','2019-04-20 11:01:36'),('CB012','Chromebook','2019-04-20 11:01:36'),('CB013','Chromebook','2019-04-20 11:01:36'),('CB014','Chromebook','2019-04-20 11:01:36'),('CB015','Chromebook','2019-04-20 11:01:36'),('CB016','Chromebook','2019-04-20 11:01:36'),('CB017','Chromebook','2019-04-20 11:01:36'),('CB018','Chromebook','2019-04-20 11:01:36'),('CB019','Chromebook','2019-04-20 11:01:36'),('CB020','Chromebook','2019-04-20 11:01:36'),('CB021','Chromebook','2019-04-20 11:01:36'),('CB022','Chromebook','2019-04-20 11:01:36'),('CB023','Chromebook','2019-04-20 11:01:36'),('CB024','Chromebook','2019-04-20 11:01:36'),('CB025','Chromebook','2019-04-20 11:01:36'),('CB026','Chromebook','2019-04-20 11:01:36'),('CB027','Chromebook','2019-04-20 11:01:36'),('CB028','Chromebook','2019-04-20 11:01:36'),('CB029','Chromebook','2019-04-20 11:01:36'),('CB030','Chromebook','2019-04-20 11:01:36'),('CB031','Chromebook','2019-04-20 11:01:36'),('CB032','Chromebook','2019-04-20 11:01:36'),('EZ01','WirelessProjector','2019-04-20 11:01:37'),('EZ02','WirelessProjector','2019-04-20 11:01:37'),('EZ03','WirelessProjector','2019-04-20 11:01:37'),('ST001','Student-iPad','2019-04-20 11:01:36'),('ST002','Student-iPad','2019-04-20 11:01:36'),('ST003','Student-iPad','2019-04-20 11:01:36'),('ST004','Student-iPad','2019-04-20 11:01:36'),('ST005','Student-iPad','2019-04-20 11:01:36'),('ST006','Student-iPad','2019-04-20 11:01:36'),('ST007','Student-iPad','2019-04-20 11:01:36'),('ST008','Student-iPad','2019-04-20 11:01:36'),('ST009','Student-iPad','2019-04-20 11:01:36'),('ST010','Student-iPad','2019-04-20 11:01:36'),('ST011','Student-iPad','2019-04-20 11:01:36'),('ST012','Student-iPad','2019-04-20 11:01:36'),('ST013','Student-iPad','2019-04-20 11:01:36'),('ST014','Student-iPad','2019-04-20 11:01:36'),('ST015','Student-iPad','2019-04-20 11:01:36'),('ST016','Student-iPad','2019-04-20 11:01:36'),('ST017','Student-iPad','2019-04-20 11:01:36'),('ST018','Student-iPad','2019-04-20 11:01:36'),('ST019','Student-iPad','2019-04-20 11:01:36'),('ST020','Student-iPad','2019-04-20 11:01:36'),('ST021','Student-iPad','2019-04-20 11:01:36'),('ST022','Student-iPad','2019-04-20 11:01:36'),('ST023','Student-iPad','2019-04-20 11:01:36'),('ST024','Student-iPad','2019-04-20 11:01:36'),('ST025','Student-iPad','2019-04-20 11:01:36'),('ST026','Student-iPad','2019-04-20 11:01:36'),('ST027','Student-iPad','2019-04-20 11:01:36'),('ST028','Student-iPad','2019-04-20 11:01:36'),('ST029','Student-iPad','2019-04-20 11:01:36'),('ST030','Student-iPad','2019-04-20 11:01:36'),('ST031','Student-iPad','2019-04-20 11:01:36'),('ST032','Student-iPad','2019-04-20 11:01:36'),('ST033','Student-iPad','2019-04-20 11:01:36'),('ST034','Student-iPad','2019-04-20 11:01:36'),('ST035','Student-iPad','2019-04-20 11:01:36'),('ST036','Student-iPad','2019-04-20 11:01:36'),('ST037','Student-iPad','2019-04-20 11:01:36'),('ST038','Student-iPad','2019-04-20 11:01:36'),('ST039','Student-iPad','2019-04-20 11:01:36'),('ST040','Student-iPad','2019-04-20 11:01:36'),('ST041','Student-iPad','2019-04-20 11:01:36'),('ST042','Student-iPad','2019-04-20 11:01:36'),('ST043','Student-iPad','2019-04-20 11:01:36'),('ST044','Student-iPad','2019-04-20 11:01:36'),('ST045','Student-iPad','2019-04-20 11:01:36'),('ST046','Student-iPad','2019-04-20 11:01:36'),('ST047','Student-iPad','2019-04-20 11:01:36'),('ST048','Student-iPad','2019-04-20 11:01:36'),('ST049','Student-iPad','2019-04-20 11:01:36'),('ST050','Student-iPad','2019-04-20 11:01:36'),('ST051','Student-iPad','2019-04-20 11:01:36'),('ST052','Student-iPad','2019-04-20 11:01:36'),('ST053','Student-iPad','2019-04-20 11:01:36'),('ST054','Student-iPad','2019-04-20 11:01:36'),('ST055','Student-iPad','2019-04-20 11:01:36'),('ST056','Student-iPad','2019-04-20 11:01:36'),('ST057','Student-iPad','2019-04-20 11:01:36'),('ST058','Student-iPad','2019-04-20 11:01:36'),('ST059','Student-iPad','2019-04-20 11:01:36'),('ST060','Student-iPad','2019-04-20 11:01:36'),('ST061','Student-iPad','2019-04-20 11:01:36'),('ST062','Student-iPad','2019-04-20 11:01:36'),('ST063','Student-iPad','2019-04-20 11:01:36'),('ST064','Student-iPad','2019-04-20 11:01:36'),('ST065','Student-iPad','2019-04-20 11:01:36'),('ST066','Student-iPad','2019-04-20 11:01:36'),('ST067','Student-iPad','2019-04-20 11:01:36'),('ST068','Student-iPad','2019-04-20 11:01:36'),('ST069','Student-iPad','2019-04-20 11:01:36'),('ST070','Student-iPad','2019-04-20 11:01:36'),('T01','Teacher-iPad','2019-04-20 11:01:36'),('T02','Teacher-iPad','2019-04-20 11:01:36'),('T03','Teacher-iPad','2019-04-20 11:01:36'),('T04','Teacher-iPad','2019-04-20 11:01:36'),('T05','Teacher-iPad','2019-04-20 11:01:36'),('T06','Teacher-iPad','2019-04-20 11:01:36'),('T07','Teacher-iPad','2019-04-20 11:01:36'),('T08','Teacher-iPad','2019-04-20 11:01:36'),('T09','Teacher-iPad','2019-04-20 11:01:36'),('T10','Teacher-iPad','2019-04-20 11:01:36'),('T11','Teacher-iPad','2019-04-20 11:01:36'),('T12','Teacher-iPad','2019-04-20 11:01:36'),('T13','Teacher-iPad','2019-04-20 11:01:36');
/*!40000 ALTER TABLE `Devices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Records`
--

DROP TABLE IF EXISTS `Records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Records` (
  `Booking` int(11) NOT NULL,
  `Device` varchar(8) NOT NULL,
  `LentFrom` datetime NOT NULL,
  `LentUntil` datetime DEFAULT NULL,
  PRIMARY KEY (`Booking`,`Device`),
  KEY `Device` (`Device`),
  CONSTRAINT `Records_ibfk_1` FOREIGN KEY (`Booking`) REFERENCES `Bookings` (`ID`),
  CONSTRAINT `Records_ibfk_2` FOREIGN KEY (`Device`) REFERENCES `Devices` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Records`
--

LOCK TABLES `Records` WRITE;
/*!40000 ALTER TABLE `Records` DISABLE KEYS */;
/*!40000 ALTER TABLE `Records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Sessions`
--

DROP TABLE IF EXISTS `Sessions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Sessions` (
  `ID` char(36) NOT NULL,
  `User` int(11) NOT NULL,
  `LastUsed` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  KEY `User` (`User`),
  CONSTRAINT `Sessions_ibfk_1` FOREIGN KEY (`User`) REFERENCES `Users` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Sessions`
--

LOCK TABLES `Sessions` WRITE;
/*!40000 ALTER TABLE `Sessions` DISABLE KEYS */;
INSERT INTO `Sessions` VALUES ('5438b0b9-2b3e-4eb9-a46f-6c8ae92dea6c',2,'2019-05-22 06:29:38'),('6967afd9-f92c-4eb2-a8df-f774e2cf1aa1',1,'2019-05-05 03:02:53'),('75309c9e-20a3-43ae-9be4-063c8a8549cd',1,'2019-05-04 07:48:25');
/*!40000 ALTER TABLE `Sessions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Users`
--

DROP TABLE IF EXISTS `Users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Users` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `Email` varchar(64) NOT NULL,
  `Name` varchar(32) NOT NULL,
  `Type` enum('Admin','Teacher') DEFAULT 'Teacher',
  `Password` blob NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `Email` (`Email`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Users`
--

LOCK TABLES `Users` WRITE;
/*!40000 ALTER TABLE `Users` DISABLE KEYS */;
INSERT INTO `Users` VALUES (1,'lancatlin@pm.me','lancat','Admin','$2a$10$uw8gXSy8mVD2jFJoC3XXDeZYsHsPl8oK0juQvVXzoIOnqbXSAe81W'),(2,'test@email.com','testman','Teacher','$2a$10$wmpNd3NtHk1volsA1XG/1uEgHKc/1Ro.Pq.cWM9up0yNOK02gjQNW');
/*!40000 ALTER TABLE `Users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-05-22 15:22:35
