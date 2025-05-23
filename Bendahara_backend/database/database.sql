-- MySQL dump 10.13  Distrib 8.0.42, for Win64 (x86_64)
--
-- Host: localhost    Database: bendahara
-- ------------------------------------------------------
-- Server version	8.0.42

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
-- Table structure for table `admin`
--

DROP TABLE IF EXISTS `admin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admin` (
  `id` varchar(65) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `nik` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `username` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `password` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `role` enum('superAdmin','bendahara','guest') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admin`
--

LOCK TABLES `admin` WRITE;
/*!40000 ALTER TABLE `admin` DISABLE KEYS */;
INSERT INTO `admin` VALUES ('4433c69f-2003-42a7-9676-ea9b9dbc9f33','123','admin','$2a$10$BVl7TJ1A8Yefr1hmXAsRdeilnHozYjUzplbpAH9fOMPvbiFEcOFwm','superAdmin');
/*!40000 ALTER TABLE `admin` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `history_transaksi`
--

DROP TABLE IF EXISTS `history_transaksi`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `history_transaksi` (
  `id_transaksi` varchar(36) NOT NULL,
  `tanggal` timestamp NOT NULL,
  `keterangan` varchar(255) DEFAULT NULL,
  `jenis_transaksi` enum('Pemasukan','Pengeluaran') NOT NULL,
  `nominal` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id_transaksi`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `history_transaksi`
--

LOCK TABLES `history_transaksi` WRITE;
/*!40000 ALTER TABLE `history_transaksi` DISABLE KEYS */;
INSERT INTO `history_transaksi` VALUES ('282bd4a7-2a9b-43a3-b62c-6b463e13701b','2025-05-21 07:34:00','test','Pengeluaran',10000),('49201e0c-2c11-4147-b3b6-4dd289205c98','2006-01-02 07:04:00','keterangan','Pemasukan',1000000),('ea55e279-d7c6-494e-97ac-cf600e9c210d','2025-05-21 07:32:00','rewrewr','Pemasukan',1000000);
/*!40000 ALTER TABLE `history_transaksi` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `laporan_keuangan`
--

DROP TABLE IF EXISTS `laporan_keuangan`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `laporan_keuangan` (
  `id_laporan` varchar(36) NOT NULL,
  `tanggal` timestamp NOT NULL,
  `keterangan` varchar(255) DEFAULT NULL,
  `pemasukan` bigint NOT NULL,
  `pengeluaran` bigint NOT NULL,
  `nota` varchar(255) DEFAULT 'no data',
  `saldo` bigint NOT NULL,
  `id_transaksi` varchar(36) NOT NULL,
  PRIMARY KEY (`id_laporan`),
  KEY `id_transaksi` (`id_transaksi`),
  CONSTRAINT `laporan_keuangan_ibfk_1` FOREIGN KEY (`id_transaksi`) REFERENCES `history_transaksi` (`id_transaksi`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `laporan_keuangan`
--

LOCK TABLES `laporan_keuangan` WRITE;
/*!40000 ALTER TABLE `laporan_keuangan` DISABLE KEYS */;
INSERT INTO `laporan_keuangan` VALUES ('04451000-1aff-4d78-a0d4-45f078e869f5','2006-01-02 07:04:00','keterangan',1000000,0,'2006-01-02-15-04-6b67684f-64cc-433e-8557-c396fe39a057.jpg',1000000,'49201e0c-2c11-4147-b3b6-4dd289205c98'),('42719ea3-176f-49f3-8cd3-3f97fa5640dd','2025-05-21 07:32:00','rewrewr',1000000,0,'no data',2000000,'ea55e279-d7c6-494e-97ac-cf600e9c210d'),('ba77b216-b77f-45af-9ae7-fc576b5c0429','2025-05-21 07:34:00','test',0,10000,'no data',1990000,'282bd4a7-2a9b-43a3-b62c-6b463e13701b');
/*!40000 ALTER TABLE `laporan_keuangan` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pemasukan`
--

DROP TABLE IF EXISTS `pemasukan`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pemasukan` (
  `id_pemasukan` varchar(36) NOT NULL,
  `tanggal` timestamp NOT NULL,
  `kategori` varchar(255) NOT NULL,
  `keterangan` varchar(255) DEFAULT NULL,
  `nominal` bigint NOT NULL,
  `nota` varchar(255) DEFAULT 'no data',
  `id_transaksi` varchar(36) NOT NULL,
  PRIMARY KEY (`id_pemasukan`),
  KEY `id_transaksi` (`id_transaksi`),
  CONSTRAINT `pemasukan_ibfk_1` FOREIGN KEY (`id_transaksi`) REFERENCES `history_transaksi` (`id_transaksi`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pemasukan`
--

LOCK TABLES `pemasukan` WRITE;
/*!40000 ALTER TABLE `pemasukan` DISABLE KEYS */;
INSERT INTO `pemasukan` VALUES ('adf8e59a-8daa-47f5-92b2-ad84d2b3f0f3','2006-01-02 07:04:00','kategori','keterangan',1000000,'2006-01-02-15-04-6b67684f-64cc-433e-8557-c396fe39a057.jpg','49201e0c-2c11-4147-b3b6-4dd289205c98'),('e0b71be9-ce92-4555-b901-495a01ea274a','2025-05-21 07:32:00','Dana Desa','rewrewr',1000000,'','ea55e279-d7c6-494e-97ac-cf600e9c210d');
/*!40000 ALTER TABLE `pemasukan` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pengeluaran`
--

DROP TABLE IF EXISTS `pengeluaran`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pengeluaran` (
  `id_pengeluaran` varchar(36) NOT NULL,
  `tanggal` timestamp NOT NULL,
  `nota` varchar(255) NOT NULL,
  `nominal` bigint NOT NULL,
  `keterangan` varchar(255) DEFAULT NULL,
  `id_transaksi` varchar(36) NOT NULL,
  PRIMARY KEY (`id_pengeluaran`),
  KEY `id_transaksi` (`id_transaksi`),
  CONSTRAINT `pengeluaran_ibfk_1` FOREIGN KEY (`id_transaksi`) REFERENCES `history_transaksi` (`id_transaksi`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pengeluaran`
--

LOCK TABLES `pengeluaran` WRITE;
/*!40000 ALTER TABLE `pengeluaran` DISABLE KEYS */;
INSERT INTO `pengeluaran` VALUES ('b2cffac1-8d79-4205-9b28-f9e8c5e8698f','2025-05-21 07:34:00','2025-05-21-15-34-fecf2952-6903-4367-a78e-51bf50785e62.jpeg',10000,'test','282bd4a7-2a9b-43a3-b62c-6b463e13701b');
/*!40000 ALTER TABLE `pengeluaran` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-05-23 17:09:30
