package main

//driver para conectar a la DB: go get -u github.com/go-sql-driver/mysql e importamos:  _ "github.com/go-sql-driver/mysql"
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contraseña := "escobar2"
	Nombre := "sistema_empleados"

	conexion, err := sql.Open(Driver, Usuario+":"+Contraseña+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

// obtener infomracion dentro de una carpeta
var plantillas = template.Must(template.ParseGlob("templates/*")) //accedemos a la carpeta /templates y lo que esta ahi dentro lo almacenamos en la var: plantillas

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", BuscarParaEditar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Println("servidor corriendo...") //mensaje de que esta corriendo el server con el fecha y hora incluida
	//fmt.Println("servidor corriendo..") --> este mensaje en consola no te dice el dia ni la fecha como el log de arriba
	http.ListenAndServe(":8080", nil)
}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

// go es case sensitive: Primer parametro para poder responder a la solicitud, segundo parametro para devolver la peticion, con r.
func Home(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionDB()
	// insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre,correo) VALUES('ivan','correo@gmail.com') ")

	// if err != nil {
	// 	panic(err.Error()) //si existe un error, mostramoss el error con el metodo: .Error()
	// }
	// //Exec, permite la insersion a la base de datos. El 'Exec' permite ejecutar codigo sql en codigo de go.
	// insertarRegistros.Exec()

	//el Query nos permite ejecutar la consulta sql sin tener que usar el 'Exec'
	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")

	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{} //indicamos que empleados va a ser un arreglo de Empleados, osea de la estructura definida.
	arregloEmpleado := []Empleado{}

	for registros.Next() { //usamos el for para recorrer todos los datos de la sentencia SELECT
		//identificamos que es lo que tenemos que rellenar
		var id int
		var nombre, correo string
		err = registros.Scan(&id, &nombre, &correo) //asignamos los valores que trae de la consulta las variables: &id, &nombre y &correo: se almacena lo que trae--> var id int, var nombre, correo string
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		//asignamos los valores directamente a 'arregloEmpleado' asi despues imprimimos ese arregloEmpleado
		arregloEmpleado = append(arregloEmpleado, empleado)
	}
	//con el Println, mostramos los datos en consola.
	//fmt.Println(arregloEmpleado)

	plantillas.ExecuteTemplate(w, "home", arregloEmpleado) //entramos en inicio.html que se encuentra en la carpeta /templates
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre") //del FormValue traemos el "nombre" que esta como id en el formuladrio de html
		correo := r.FormValue("correo")

		conexionEstablecida := conexionDB()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre,correo) VALUES(?,?) ") //los signos se le agrega en lugares donde se ingresara texto

		if err != nil {
			panic(err.Error())
		}

		insertarRegistros.Exec(nombre, correo) //nombre y correo vienen de los '?' osea en el 'INSERT' de Prepare()

		http.Redirect(w, r, "/", 301) //301 codigo de redireccion a una url
	}
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	//el idEmpleado va a ser igual al Request que me estan enviando
	idEmpleado := r.URL.Query().Get("id")
	//MOSTRAMOS LOS ID POR CONSOLA
	fmt.Println(idEmpleado)

	conexionEstablecida := conexionDB()                                                     //conexion a base de datos
	borrarRegistros, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?") //ejecutamos el borrado con el parametro "id"

	if err != nil {
		panic(err.Error())
	}
	borrarRegistros.Exec(idEmpleado)
	http.Redirect(w, r, "/", 301)
}

func BuscarParaEditar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := conexionDB()
	editarRegistros, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)

	empleado := Empleado{}
	for editarRegistros.Next() {
		var id int
		var nombre, correo string
		err = editarRegistros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}
	fmt.Println(empleado)
	plantillas.ExecuteTemplate(w, "editar", empleado)
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionDB()

		actualizarRegistros, err := conexionEstablecida.Prepare("UPDATE empleados SET nombre = ?, correo = ? WHERE id = ? ") //los signos se le agrega en lugares donde se ingresara texto

		if err != nil {
			panic(err.Error())
		}

		actualizarRegistros.Exec(nombre, correo, id) //nombre y correo vienen de los '?'
		http.Redirect(w, r, "/", 301)                //301 codigo de redireccion a una url
	}
}
