package main

import (
	"os"
	"strconv"
	"strings"
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

	// I check the read functions worked well
	for i := 0; i < len(products); i++ {
		println(products[i].Categoria+" "+products[i].ID+" "+products[i].Nombre+" "+strconv.FormatFloat(products[i].Precio, 'f', 2, 64)+" "+strconv.Itoa(products[i].Stock))
		println(transactions[i].Tipo+" "+transactions[i].IDProducto+" "+strconv.Itoa(transactions[i].Cantidad)+" "+transactions[i].Fecha)
	}
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
	return nil
}
func escribirInventario(productos []Producto, nombreArchivo string) error {
	return nil
}
func generarReporteBajoStock(productos []Producto, limite int) error {
	return nil
}
func escribirLog(errores []string, nombreArchivo string) error {
	return nil
}