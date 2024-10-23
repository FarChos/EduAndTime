CREATE DATABASE IF NOT EXISTS EduAndTime CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

CREATE USER IF NOT EXISTS 'admin' IDENTIFIED BY 'tearsthemoon';
GRANT ALL PRIVILEGES ON EduAndTime. * TO 'admin'@'%';
FLUSH PRIVILEGES;

USE EduAndTime;

CREATE TABLE IF NOT EXISTS usuarios (
  idUsuario INT AUTO_INCREMENT,
  nombre VARCHAR(255) UNIQUE NOT NULL,
  correo VARCHAR(255) UNIQUE NOT NULL,
  contrasena VARCHAR(255) NOT NULL,
  direcImgPerfil TEXT NOT NULL, 
  fechaOrigen DATE NOT NULL,
    PRIMARY KEY (idUsuario)
);

CREATE TABLE IF NOT EXISTS documentos (
  idDoc INT AUTO_INCREMENT,
  titulo VARCHAR(255) NOT NULL,
  autor VARCHAR(255) NOT NULL,
  categoria ENUM('sociedad', 'geografía', 'tecnologia', 'ciencia', 'economía', 'bienestar', 'política', 'arte', 'filosofia', 'exotica') NOT NULL,
  idUsuario INT NOT NULL,
  formato ENUM('PDF', 'EPUB', 'MOBI', 'TXT') NOT NULL,
  descripción TEXT NOT NULL,
  fechaOrigen DATE NOT NULL,
  DirecDoc VARCHAR(255) NOT NULL,
    PRIMARY KEY (idDoc),
    FOREIGN KEY (idUsuario) REFERENCES usuarios(idUsuario)
);

CREATE TABLE IF NOT EXISTS chat (
  idChat INT AUTO_INCREMENT,
  nombreChat VARCHAR(255) NOT NULL,
  direcIcono VARCHAR(255) NOT NULL,
  direcFondo VARCHAR(255) NOT NULL,
  descripción TEXT NOT NULL,
  fechaOrigen DATE NOT NULL,
  esPublico BOOLEAN NOT NULL,
  esVisible BOOLEAN NOT NULL,
  esAbierto BOOLEAN NOT NULL, 
    PRIMARY KEY (idChat)
);

CREATE INDEX IF NOT EXISTS inx_nombreChat
ON chat(nombreChat);

CREATE INDEX IF NOT EXISTS inx_titulo
ON documentos (titulo);
CREATE INDEX IF NOT EXISTS inx_autor
ON documentos (autor);
CREATE INDEX IF NOT EXISTS inx_categoria
ON documentos (categoria);
CREATE INDEX IF NOT EXISTS inx_formato
ON documentos (formato);

CREATE UNIQUE INDEX IF NOT EXISTS idx_correo
ON usuarios (correo);