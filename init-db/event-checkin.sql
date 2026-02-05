-- MySQL dump 10.13  Distrib 8.4.7, for Linux (x86_64)
--
-- Host: localhost    Database: event_checkin
-- ------------------------------------------------------
-- Server version	8.4.7

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `events`
--

DROP TABLE IF EXISTS `events`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `events` (
  `id` varchar(36) NOT NULL,
  `organizer_id` bigint unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `slug` varchar(255) NOT NULL,
  `date` datetime NOT NULL,
  `venue` varchar(500) NOT NULL,
  `participant_count` int NOT NULL DEFAULT '0',
  `total_price` int NOT NULL DEFAULT '0',
  `payment_status` enum('pending','verified','active') NOT NULL DEFAULT 'pending',
  `payment_proof_url` varchar(500) DEFAULT NULL,
  `scanner_pin` char(4) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `slug` (`slug`),
  KEY `idx_events_organizer_id` (`organizer_id`),
  KEY `idx_events_slug` (`slug`),
  KEY `idx_events_payment_status` (`payment_status`),
  CONSTRAINT `events_ibfk_1` FOREIGN KEY (`organizer_id`) REFERENCES `organizers` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `events`
--

LOCK TABLES `events` WRITE;
/*!40000 ALTER TABLE `events` DISABLE KEYS */;
INSERT INTO `events` VALUES ('af6dc11a-ff14-4f2a-bcbb-b3ad2c0809c6',2,'DevOps Learning','devops-learning','2026-02-14 00:00:00','Kampus ITATS',10,50000,'pending','','6702','2026-01-22 12:01:06'),('eb375f58-2d10-49ba-865e-347764ac3c0f',1,'GO Developer','go-developer','2026-01-31 00:00:00','Kampus ITATS Gedung H',10,50000,'pending','','0793','2026-01-14 03:05:46');
/*!40000 ALTER TABLE `events` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `organizers`
--

DROP TABLE IF EXISTS `organizers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `organizers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  KEY `idx_organizers_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `organizers`
--

LOCK TABLES `organizers` WRITE;
/*!40000 ALTER TABLE `organizers` DISABLE KEYS */;
INSERT INTO `organizers` VALUES (1,'test@example.com','fizo','$2a$10$pvKGgeL7A86xbRfcGQIsSOnF7/sg5qNa16Ez5uUDgXC0LetLmXAMC','2026-01-06 02:16:14'),(2,'ytp@example.com','Yantop','$2a$10$tH9thgaLZ8D1b.KbVEDE6uoklSVzvDgfKltta3TuMThpVy0vOxFG.','2026-01-22 11:56:39');
/*!40000 ALTER TABLE `organizers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `participants`
--

DROP TABLE IF EXISTS `participants`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `participants` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `event_id` varchar(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `phone` varchar(20) NOT NULL,
  `qr_token` varchar(100) NOT NULL,
  `qr_sent` tinyint(1) DEFAULT '0',
  `qr_sent_at` timestamp NULL DEFAULT NULL,
  `checked_in` tinyint(1) DEFAULT '0',
  `checked_in_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `qr_token` (`qr_token`),
  KEY `idx_participants_event_id` (`event_id`),
  KEY `idx_participants_qr_token` (`qr_token`),
  KEY `idx_participants_checked_in` (`checked_in`),
  KEY `idx_qr_sent` (`qr_sent`),
  CONSTRAINT `participants_ibfk_1` FOREIGN KEY (`event_id`) REFERENCES `events` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=111 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `participants`
--

LOCK TABLES `participants` WRITE;
/*!40000 ALTER TABLE `participants` DISABLE KEYS */;
INSERT INTO `participants` VALUES (91,'eb375f58-2d10-49ba-865e-347764ac3c0f','Dea Nurinda','deanurindasalsadil@gmail.com','0896-1020-3040','ccf8c1761b01e7d12ca8ab9bf8abda74',1,'2026-01-14 03:06:45',0,NULL,'2026-01-14 03:06:17'),(92,'eb375f58-2d10-49ba-865e-347764ac3c0f','Tari','deanurindasalsadila@gmail.com','0812-3456-7890','e581279ec500a844712fe54b04e154ad',1,'2026-01-14 03:07:01',1,'2026-01-14 04:30:08','2026-01-14 03:06:17'),(93,'eb375f58-2d10-49ba-865e-347764ac3c0f','Sunarto','sunarto1114@gmail.com','0878-1122-3344','720ee358bf0789f2cc31996e3c09baa4',1,'2026-01-14 03:07:16',0,NULL,'2026-01-14 03:06:17'),(94,'eb375f58-2d10-49ba-865e-347764ac3c0f','oca','fizonendapocasondaya@gmail.com','0856-9876-5432','0af6cf8304b6b006ee043a764d6f1d19',1,'2026-01-14 03:07:31',1,'2026-01-14 13:42:12','2026-01-14 03:06:17'),(95,'eb375f58-2d10-49ba-865e-347764ac3c0f','Fajar Nugroho','fajar.nugroho@gmail.com','0821-5566-7788','094891589822764cd73fd4e48e22d453',1,'2026-01-14 03:07:47',0,NULL,'2026-01-14 03:06:17'),(96,'eb375f58-2d10-49ba-865e-347764ac3c0f','Siti Hajar','siti.hajar@gmail.com','0813-2244-6688','29cb6ab7f33db4303f90ceb7b7cc31b8',1,'2026-01-14 03:08:02',0,NULL,'2026-01-14 03:06:17'),(97,'eb375f58-2d10-49ba-865e-347764ac3c0f','Ahmad Yani','ahmad.yani@gmail.com','0877-0011-2233','7824c6859eafef83624e3822508c6380',1,'2026-01-14 03:08:17',0,NULL,'2026-01-14 03:06:17'),(98,'eb375f58-2d10-49ba-865e-347764ac3c0f','Maya Sari Dewi','maya.sari@gmail.com','0852-4433-2211','3ce08d08a1e4e81e6c70e3f40fe80b8c',1,'2026-01-14 03:08:33',0,NULL,'2026-01-14 03:06:17'),(99,'eb375f58-2d10-49ba-865e-347764ac3c0f','Eko Prasetya','eko.prasetya@gmail.com','0815-6789-0123','9353a456c9de82174ea7d4f52da2ceb9',1,'2026-01-14 03:08:48',0,NULL,'2026-01-14 03:06:17'),(100,'eb375f58-2d10-49ba-865e-347764ac3c0f','Fizonenda','fizonenda18@gmail.com','0822-9988-7766','0f616ac5b9bfe25eb45d5fb0e42e08f3',1,'2026-01-14 03:09:03',1,'2026-01-14 13:42:38','2026-01-14 03:06:17'),(101,'af6dc11a-ff14-4f2a-bcbb-b3ad2c0809c6','Dea Nurinda','deanurindasalsadil@gmail.com','0896-1020-3040','a54e3800cf311376f74bd04e488831da',1,'2026-01-22 12:02:41',0,NULL,'2026-01-22 12:01:59'),(102,'af6dc11a-ff14-4f2a-bcbb-b3ad2c0809c6','Tari','deanurindasalsadila@gmail.com','0812-3456-7890','c5d43fe57875015a2cdaa929ac540941',1,'2026-01-22 12:02:59',0,NULL,'2026-01-22 12:01:59'),(103,'af6dc11a-ff14-4f2a-bcbb-b3ad2c0809c6','Sunarto','sunarto1114@gmail.com','0878-1122-3344','dff1fbed4b7e225fb1498e4799c90312',1,'2026-01-22 12:03:19',1,'2026-01-22 12:13:01','2026-01-22 12:01:59'),(104,'af6dc11a-ff14-4f2a-bcbb-b3ad2c0809c6','oca','fizonendapocasondaya@gmail.com','0856-9876-5432','94161a01236fc522b3a6e4f8f2f206a8',1,'2026-01-22 12:03:37',0,NULL,'2026-01-22 12:01:59'),(105,'af6dc11a-ff14-4f2a-bcbb-b3ad2c0809c6','Fajar Nugroho','fajar.nugroho@gmail.com','0821-5566-7788','728baf15f73fe43ebb2673e95a028980',1,'2026-01-22 12:03:59',0,NULL,'2026-01-22 12:01:59'),(106,'af6dc11a-ff14-4f2a-bcbb-b3ad2c0809c6','Siti Hajar','siti.hajar@gmail.com','0813-2244-6688','b249bf2da63a57c7cc9f1e74c1a8e1ef',1,'2026-01-22 12:04:17',0,NULL,'2026-01-22 12:01:59'),(107,'af6dc11a-ff14-4f2a-bcbb-b3ad2c0809c6','Ahmad Yani','ahmad.yani@gmail.com','0877-0011-2233','a21168cb11e5eb590472f83b03cffc15',1,'2026-01-22 12:04:34',0,NULL,'2026-01-22 12:01:59'),(108,'af6dc11a-ff14-4f2a-bcbb-b3ad2c0809c6','Maya Sari Dewi','maya.sari@gmail.com','0852-4433-2211','9b4710e18b8e149fda2f8a77714b1166',1,'2026-01-22 12:04:54',0,NULL,'2026-01-22 12:01:59'),(109,'af6dc11a-ff14-4f2a-bcbb-b3ad2c0809c6','Eko Prasetya','eko.prasetya@gmail.com','0815-6789-0123','5946569ef166c56a24e583676b0589c4',1,'2026-01-22 12:05:12',0,NULL,'2026-01-22 12:01:59'),(110,'af6dc11a-ff14-4f2a-bcbb-b3ad2c0809c6','Fizonenda','fizonenda18@gmail.com','0822-9988-7766','2240e6bb78852a5bbd6dc7bc3d671b24',1,'2026-01-22 12:05:31',0,NULL,'2026-01-22 12:01:59');
/*!40000 ALTER TABLE `participants` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `schema_migrations`
--

DROP TABLE IF EXISTS `schema_migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `schema_migrations` (
  `version` bigint NOT NULL,
  `dirty` tinyint(1) NOT NULL,
  PRIMARY KEY (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `schema_migrations`
--

LOCK TABLES `schema_migrations` WRITE;
/*!40000 ALTER TABLE `schema_migrations` DISABLE KEYS */;
INSERT INTO `schema_migrations` VALUES (7,0);
/*!40000 ALTER TABLE `schema_migrations` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-02-01 18:02:46
