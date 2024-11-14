package service

import (
	"ModuloAuth/db"
	"ModuloAuth/graph/model"
	"fmt"
	"time"

	"strconv"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

func BuscarUsuario(correo string, contrasena string) (model.TokenAuth, error) {
	tokenAuth := model.TokenAuth{Exito: false}
	var correoUs, nombreUs, direcImgPerfilUs string
	var contrasenaUs []byte

	// Pasa &direcImgPerfilUs como puntero
	err := db.MariaDB.QueryRow("SELECT correo, nombre, contrasena, COALESCE(direcImgPerfil, '') AS direcImgPerfil FROM usuarios WHERE correo = ?", correo).Scan(&correoUs, &nombreUs, &contrasenaUs, &direcImgPerfilUs)
	if err != nil {
		return tokenAuth, err
	}

	contrasenaAlmacenada := []byte(contrasenaUs)
	contrasenaByte := []byte(contrasena)

	// Compara la contraseña hasheada con la contraseña ingresada
	err = bcrypt.CompareHashAndPassword(contrasenaAlmacenada, contrasenaByte)
	if err != nil {
		return tokenAuth, nil // Contraseña incorrecta
	}

	horaExpiracion := time.Now().Add(1 * time.Hour).Unix()

	token, err := generartoken(horaExpiracion, nombreUs)
	if err != nil {
		return tokenAuth, err
	}

	// Construir y retornar el objeto TokenAuth en caso de éxito
	tokenAuth = model.TokenAuth{
		Token: &token,
		Usuario: &model.Usuario{
			Nombre:  nombreUs,
			Correo:  correoUs,
			ImgPerf: &direcImgPerfilUs,
		},
		Exito:  true,
		Expira: &horaExpiracion,
	}

	return tokenAuth, nil
}

func generartoken(horaExpiracion int64, nombre string) (string, error) {

	claims := &Claims{
		UserID: nombre,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: horaExpiracion,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("viajeAntesQueDestino"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CrearUsuario(usuario *model.Usuario, contrasena string) (bool, error) {
	// Formatear la fecha de creación
	diaCreacion := time.Now().Format("2006-01-02")
	var idUsuario int

	// Generar hash de la contraseña
	contrasenaHash, err := bcrypt.GenerateFromPassword([]byte(contrasena), bcrypt.DefaultCost)
	if err != nil {
		return false, fmt.Errorf("error al hashear la contraseña: %v", err)
	}

	// Insertar en MariaDB
	_, err = db.MariaDB.Exec("INSERT INTO usuarios (nombre, correo, contrasena, direcImgPerfil, fechaOrigen) VALUES (?, ?, ?, ?, ?)", usuario.Nombre, usuario.Correo, contrasenaHash, usuario.ImgPerf, diaCreacion)
	if err != nil {
		return false, fmt.Errorf("error al insertar el usuario en MariaDB: %v", err)
	}

	// Obtener el ID recién insertado
	err = db.MariaDB.QueryRow("SELECT idUsuario FROM usuarios WHERE correo = ?", usuario.Correo).Scan(&idUsuario)
	if err != nil {
		return false, fmt.Errorf("error al buscar el id para Couchbase: %v", err)
	}

	// Crear el documento para Couchbase
	doc := map[string]interface{}{
		"tipo":       "usuario",
		"estaActivo": true,
		"version":    1,
		"idUsuario":  idUsuario,
	}

	// Usar el ID como clave del documento
	docID := "usuario::" + strconv.Itoa(idUsuario)

	// Acceder al bucket y colección de Couchbase
	bucket := db.CouchbaseCluster.Bucket("EduAndTime")
	err = bucket.WaitUntilReady(5*time.Second, nil) // Asegura que el bucket esté listo
	if err != nil {
		return false, fmt.Errorf("error al esperar disponibilidad del bucket en Couchbase: %v", err)
	}

	collection := bucket.Scope("usuarios").Collection("usuario")

	// Insertar en Couchbase
	_, err = collection.Insert(docID, doc, nil)
	if err != nil {
		return false, fmt.Errorf("error al insertar el usuario en Couchbase: %v", err)
	}

	return true, nil
}

func ActualizarUsuario(usuario *model.Usuario, contrasena string) (bool, error) {

	contrasenaHash, err := bcrypt.GenerateFromPassword([]byte(contrasena), bcrypt.DefaultCost)
	if err != nil {
		return false, fmt.Errorf("error al hashear la contraseña: %v", err)
	}

	_, err = db.MariaDB.Exec("UPDATE usuarios SET nombre = ?, contrasena = ? WHERE correo = ?", usuario.Nombre, contrasenaHash, usuario.Correo)
	if err != nil {
		return false, fmt.Errorf("error al actualizar el usuario: %v", err)
	}
	return true, nil
}

func EliminarUsuario(correo string) (bool, error) {
	_, err := db.MariaDB.Exec("DELETE FROM usuarios WHERE correo = ?", correo)
	if err != nil {
		return false, fmt.Errorf("error al eliminar el usuario: %v", err)
	}
	return true, nil
}
