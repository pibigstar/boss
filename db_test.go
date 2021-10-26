package main

import (
	"testing"
)

// "gen model -dns "mysql:root:123456@tcp(127.0.0.1:3306)/boss" -t user"
func TestListUser(t *testing.T) {
	users, err := listUserFromDB()
	if err != nil {
		t.Error(err)
	}
	for _, u := range users {
		t.Logf("%+v \n", u)
	}
}
