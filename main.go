package main

import (
	"fmt"
	"strconv"

	"github.com/golang/leveldb"
	"golang.org/x/crypto/sha3"
)

type DB struct {
	Db *leveldb.DB
}

func concatBytes(b1 []byte, b2 []byte) []byte {
	totalSize := len(b1) + len(b2)

	h := make([]byte, totalSize)

	for i, current := range b1 {
		h[i] = current
	}

	for i, current := range b2 {
		h[i+len(b1)] = current
	}

	return h
}

func (db *DB) hashify(length int) [64]byte {
	h := make([]byte, 12)

	for i := 0; i < length; i++ {
		buf, err := db.Db.Get([]byte(strconv.Itoa(i)), nil)
		if err != nil {
			fmt.Println("panicing: " + err.Error())
			panic(err.Error)
		}

		h = concatBytes(h, buf)
		sha3.Sum512(h)
	}

	return sha3.Sum512(h)
}

func main() {
	db, err := leveldb.Open("db", nil)
	if err != nil {
		fmt.Println("panicing: " + err.Error())
		panic(err.Error)
	}
	defer db.Close()

	data := []string{
		"A,0,5",
		"B,0,7",
		"C,0,3",
		"D,0,1",
		"E,0,9",
		"F,0,6",
	}

	for i := range data {
		db.Set([]byte(strconv.Itoa(i)), []byte(data[i]), nil)
	}

	d := DB{Db: db}

	fmt.Println(d.hashify(len(data)))

	// A ---1----> D
	db.Set([]byte(strconv.Itoa(0)), []byte("A,0,4"), nil)
	db.Set([]byte(strconv.Itoa(3)), []byte("D,0,2"), nil)

	fmt.Println(d.hashify(len(data)))
	fmt.Println(len(d.hashify(len(data))))

	// for test commit.
}
