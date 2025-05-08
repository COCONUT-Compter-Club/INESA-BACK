-- MySQL dump 10.13  Distrib 8.0.41, for Linux (x86_64)
--
-- Host: yamabiko.proxy.rlwy.net    Database: railway
-- ------------------------------------------------------
-- Server version	9.2.0

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
-- Table structure for table `Admin`
--

DROP TABLE IF EXISTS `Admin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Admin` (
  `id` varchar(65) COLLATE utf8mb4_general_ci NOT NULL,
  `nik` varchar(20) COLLATE utf8mb4_general_ci NOT NULL,
  `username` varchar(20) COLLATE utf8mb4_general_ci NOT NULL,
  `password` text COLLATE utf8mb4_general_ci NOT NULL,
  `role` enum('superAdmin','bendahara','guest') COLLATE utf8mb4_general_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Admin`
--

LOCK TABLES `Admin` WRITE;
/*!40000 ALTER TABLE `Admin` DISABLE KEYS */;
INSERT INTO `Admin` VALUES ('4433c69f-2003-42a7-9676-ea9b9dbc9f33','123','admin','$2a$10$BVl7TJ1A8Yefr1hmXAsRdeilnHozYjUzplbpAH9fOMPvbiFEcOFwm','superAdmin');
/*!40000 ALTER TABLE `Admin` ENABLE KEYS */;
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
INSERT INTO `history_transaksi` VALUES ('62694789-0edc-4c23-8bab-e5f2a2c408b8','2025-05-03 20:31:00','fgh','Pengeluaran',50000),('e7327991-89f6-4a4f-aa22-bf5f244f7b5d','2025-05-03 20:30:00','vghfg','Pemasukan',100000);
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
INSERT INTO `laporan_keuangan` VALUES ('280ae801-abac-47b0-91d4-41004fa31caa','2025-05-03 20:30:00','vghfg',100000,0,100000,'e7327991-89f6-4a4f-aa22-bf5f244f7b5d'),('d80ee1c8-7e5a-4176-94f8-731a85fa29e8','2025-05-03 20:31:00','fgh',0,50000,50000,'62694789-0edc-4c23-8bab-e5f2a2c408b8');
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
INSERT INTO `pemasukan` VALUES ('a9accb33-978d-4e2b-b635-c386fb880068','2025-05-03 20:30:00','Dana Desa','vghfg',100000,'','e7327991-89f6-4a4f-aa22-bf5f244f7b5d');
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
INSERT INTO `pengeluaran` VALUES ('1f0b9dad-0418-4f72-8c03-339b3d54373f','2025-05-03 20:31:00','2025-05-03-20-31-fe42aa1b-24fe-44d9-a5ca-21704f902662.jpeg',50000,'fgh','62694789-0edc-4c23-8bab-e5f2a2c408b8');
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

-- Dump completed on 2025-05-03 21:47:45
