package main

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Read inventory from "inventario.txt"

	products, errorInventario := leerInventario("inventario.txt")
	if errorInventario != nil {
		// Handle error
		return
	}

	// Read transactions from "transacciones.txt"
	transactions, errorTransacciones := leerTransacciones("transacciones.txt")
	if errorTransacciones != nil {
		// Handle error
		return
	}

	// Update products
	salida := procesarTransacciones(products,transactions)
	for i  := range salida {
		println(salida[i])
	}

	// Update the log
	err := escribirLog(salida, "log.txt")
	if err != nil {
		println("Ha habido algun error al actualizar el log")
		println(err)
	}

	// Write new inventory
	escribirInventario(products, "inventario_actualizado.txt")

	// Generate a report with low stock products
	generarReporteBajoStock(products, 10)
}


type Producto struct {
ID string
Nombre string
Categoria string
Precio float64
Stock int
}
type Transaccion struct {
Tipo string
IDProducto string
Cantidad int
Fecha string
}
func leerInventario(nombreArchivo string) ([]Producto, error) {
	buffer, err := os.ReadFile(nombreArchivo)
	arrayProducts := make([]Producto, 0)
	if (err != nil) {
		return nil, err
	}
	content := string(buffer)
	newContent := strings.Split(content, "\n")
	// According to the structure I have to jump the first line (0)
	for i := 1; i < len(newContent); i++ {

		line:= newContent[i]
		lineMembers := strings.Split(line, ",")
		// I have to control 2 values to string to other format
		Precio, err :=strconv.ParseFloat(strings.TrimSpace(lineMembers[3]), 64)

		// If format give an error, I return a nil
		if err != nil {
			return nil, err
		}

		// Here the same for int
		Stock, err := strconv.Atoi(strings.TrimSpace(lineMembers[4]))
		if err != nil {
			return nil, err
		}

		// I build a product and append it to the array
		arrayProducts = append(arrayProducts, Producto{
			ID: lineMembers[0],
			Nombre: lineMembers[1],
			Categoria: lineMembers[2],
			Precio: Precio,
			Stock: Stock,
		})
	}
	return arrayProducts, nil
}
func leerTransacciones(nombreArchivo string) ([]Transaccion, error) {
	buffer, err := os.ReadFile(nombreArchivo)
	arrayTransaction := make([]Transaccion, 0)
	if (err != nil) {
		return nil, err
	}
	content := string(buffer)
	newContent := strings.Split(content, "\n")

	// According to the structure I have to jump the first line (0)
	for i := 1; i < len(newContent); i++ {

		line:= newContent[i]
		lineMembers := strings.Split(line, ",")

		// I have to control 1 values to string to other format
		Cantidad, err := strconv.Atoi(strings.TrimSpace(lineMembers[2]))

		// If format give an error, I return a nil
		if err != nil {
			return nil, err
		}


		// I build a product and append it to the array
		arrayTransaction = append(arrayTransaction, Transaccion{
			Tipo: lineMembers[0],
			IDProducto: lineMembers[1],
			Cantidad: Cantidad,
			Fecha: lineMembers[3],
			
		})
	}
	return arrayTransaction, nil
}
func procesarTransacciones(productos []Producto, transacciones []Transaccion) []string {
	var resultados []string

	// Consider transaction as a parent
	for i := 0; i < len(transacciones); i++ {
		transacion := transacciones[i]

		// Consider productos as a child
		for j := 0; j < len(productos); j++ {
			// First condition the ID exist (the same ID)
			if productos[j].ID == transacion.IDProducto {
				// Check the stock, in addition I have to controlate the case when is COMPRA or DEVOLUCION
				if (transacion.Tipo == "VENTA" && productos[j].Stock >= transacion.Cantidad) ||
				transacion.Tipo == "COMPRA" || transacion.Tipo == "DEVOLUCION" {
					if transacion.Tipo == "VENTA" && productos[j].Stock >= transacion.Cantidad {
						productos[i].Stock = productos[i].Stock - transacion.Cantidad
					} else {
						productos[i].Stock = productos[i].Stock + transacion.Cantidad
					}
				} else {
					// Error for the log when you don't have enough stock
					ts := time.Now().Format("2006-01-02 15:04:05")
					resultados = append(resultados, "["+ts+"] Error: Stock insufficient for sale of product "+transacion.IDProducto+" not found in transaction of type "+transacion.Tipo)
					resultados = append(resultados, "actual: "+strconv.Itoa(productos[j].Stock)+", Cantidad solicitada: "+strconv.Itoa(transacion.Cantidad))
				}
				break
			} else {
				if j+1==len(productos) {
					// Error for the log when the ID isn't equals
					ts := time.Now().Format("2006-01-02 15:04:05")
					resultados = append(resultados, "["+ts+"] Error: Product "+transacion.IDProducto+" not found in transaction of type "+transacion.Tipo)
				}
			}
		}
	}

	return resultados
}
func escribirInventario(productos []Producto, nombreArchivo string) error {
	

	// Add the first line and concat all array
	contenido := "ID,Nombre,Categoría,Precio,Stock\n"
	for i := 0; i < len(productos); i++ {
		contenido+=productos[i].ID+","+productos[i].Nombre+","+productos[i].Categoria+","+strconv.FormatFloat(productos[i].Precio, 'f', 2, 64)+","+strconv.Itoa(productos[i].Stock)+"\n"
	}
	// I write the file
	err := os.WriteFile(nombreArchivo, []byte(contenido), 0644)
	return err
}
func generarReporteBajoStock(productos []Producto, limite int) error {
	//  We open the file and control the success
	fichero, err := os.OpenFile("productos_bajo_stock.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fichero.Close()

	// First line of the report
	contenido := "ALERTA: PRODUCTOS CON BAJO STOCK\n================================\n\n"

	// Count all alerts
	contador := 0
	for i := 0; i < len(productos); i++ {
		if productos[i].Stock < limite {
			contenido+=productos[i].ID+" | "+productos[i].Nombre+" | Stock actual: "+strconv.Itoa(productos[i].Stock)+" unidades\n"
			contador++
		}
	}

	// Add Final line with the total of products with low stock
	contenido += "\nTotal productos con bajo stock: "+strconv.Itoa(contador)+"\n"
	_, err = fichero.WriteString(contenido)
	return err
}
func escribirLog(errores []string, nombreArchivo string) error {
	// Concat the string to add the file
	contenido := ""
	for i := 0; i < len(errores); i++ {
		contenido+=errores[i]
	}

	// Open the file and control the success
	fichero, err := os.OpenFile(nombreArchivo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil
	}
	defer fichero.Close()
	_, err  = fichero.WriteString(contenido)
	return err 
}