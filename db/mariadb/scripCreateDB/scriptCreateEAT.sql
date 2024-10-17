CREATE DATABASE IF NOT EXISTS EduAndTime;

USE EduAndTime;

CREATE TABLE usuarios (
  idUsuario INT AUTO_INCREMENT NOT NULL,
  nombre VARCHAR(255) UNIQUE NOT NULL,
  correo VARCHAR(255) UNIQUE NOT NULL,
  contrasena VARCHAR(255) NOT NULL,
  DirecImgPerfil VARCHAR(255) NOT NULL, 
  fechaOrigen DATE NOT NULL,
    PRIMARY KEY (idUsuario)
);

CREATE TABLE documentos (
  idDoc INT AUTO_INCREMENT NOT NULL,
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

CREATE TABLE chat (
  idChat INT AUTO_INCREMENT NOT NULL,
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

CREATE INDEX inx_nombreChat
ON chat(nombreChat);

CREATE INDEX inx_titulo
ON documentos (titulo);
CREATE INDEX inx_autor
ON documentos (autor);
CREATE INDEX inx_categoria
ON documentos (categoria);
CREATE INDEX inx_formato
ON documentos (formato);

CREATE UNIQUE INDEX idx_correo
ON usuarios (correo);
