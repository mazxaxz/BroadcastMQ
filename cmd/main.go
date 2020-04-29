package main

func main() {
	// TODO: parse yaml

	http := &Http{}
	http.ServeHTTP(":8080")
}