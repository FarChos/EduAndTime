USE EduAndTime;

INSERT INTO usuarios(nombre, correo, contrasena, direcImgPerfil, fechaOrigen) VALUES ("yo", "yo@gmail.com", "1234", "en algun lugar", CURDATE())

INSERT INTO documentos(titulo, autor, categoria, idUsuario, formato, descripcion, fechaOrigen, DirecDoc) VALUES ("don quijote", "miguel de cervantes", "arte", 1, "PDF", "un caballero de cuento", CURDATE(), "en algun lugar de la mancha")

INSERT INTO chats (nombreChat, direcIcono, direcFondo, descripcion, fechaOrigen, esPublico, esVisible, esAbierto) VALUES ("orticultura", "entre las matas", "tambien entre las matas", "Un chat dedicado a las bellas plantas de jardin", CURDATE(), FALSE, FALSE, FALSE)