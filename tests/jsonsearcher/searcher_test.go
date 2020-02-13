package searchertest

import (
	"fmt"
	"testing"

	"github.com/markity/goutils/jsonsearcher"
)

var jsonString = `
{
	"name":"Markity",
	"age":16,
	"friends":[
		{
			"name":"Jack",
			"age":17
		},
		{
			"name":"Mary",
			"age":18,
			"email":"3402002560@qq.com"
		}
	],
	"details":{
		"interests":["golang","python"]
	},
	"phone":null
}
`

func TestA(t *testing.T) {
	_, err := jsonsearcher.New([]byte(jsonString))
	if err != nil {
		t.Fatalf("err is %v, expected nil", err)
	}
	_, err = jsonsearcher.New([]byte(jsonString + "}"))
	if err == nil {
		t.Fatalf("err is nil, expected not nil")
	}
	_, err = jsonsearcher.New([]byte("null"))
	if err != nil {
		t.Fatalf("err is nil, expected not nil")
	}
}

func TestB(t *testing.T) {
	s, _ := jsonsearcher.New([]byte(jsonString))

	// exists
	r1 := s.Query("name")
	if !r1.Exists() {
		t.Fatalf("r1.Exists() is false, expected true")
	}
	if r1.Type() != jsonsearcher.TypeString {
		t.Fatalf("r1.Type is not TypeString, expected TypeString")
	}
	if r1.GetString() != "Markity" {
		t.Fatalf("r1.GetString() is %v, expected Markity", r1.GetString())
	}

	r2 := s.Query("age")
	if !r2.Exists() {
		t.Fatalf("r2.Exists() is false, expected true")
	}
	if r2.Type() != jsonsearcher.TypeNumber {
		t.Fatalf("r2.Type is not TypeNumber, expected TypeNumber")
	}
	if r2.GetInt64() != 16 {
		t.Fatalf("r2.GetInt64() is %v, expected 16", r2.GetInt64())
	}

	r3 := s.Query(("friends"))
	if !r3.Exists() {
		t.Fatalf("r3.Exists() is false, expected true")
	}
	if r3.Type() != jsonsearcher.TypeArray {
		t.Fatalf("r3.Type is not TypeArray, expected TypeArray")
	}
	if fmt.Sprintf("%v", r3.GetArray()) != "[map[age:17 name:Jack] map[age:18 email:3402002560@qq.com name:Mary]]" {
		t.Fatalf("r3.GetArray() is %v, expected [map[age:17 name:Jack] map[age:18 email:3402002560@qq.com name:Mary]]", r3.GetArray())
	}

	r4 := s.Query("friends", 0)
	if !r4.Exists() {
		t.Fatalf("r4.Exists() is false, expected true")
	}
	if r4.Type() != jsonsearcher.TypeObject {
		t.Fatalf("r4.Type is not TypeObject, expected TypeObject")
	}
	if fmt.Sprintf("%v", r4.GetObject()) != "map[age:17 name:Jack]" {
		t.Fatalf("r4.GetObject() is %v, expected map[age:17 name:Jack]", r4.GetObject())
	}

	r5 := s.Query("friends", 0, "name")
	if !r5.Exists() {
		t.Fatalf("r5.Exists() is false, expected true")
	}
	if r5.Type() != jsonsearcher.TypeString {
		t.Fatalf("r5.Type is not TypeString, expected TypeString")
	}
	if r5.GetString() != "Jack" {
		t.Fatalf("r5.GetString() is %v, expected Jack", r5.GetString())
	}

	r6 := s.Query("friends", 0, "age")
	if !r6.Exists() {
		t.Fatalf("r6.Exists() is false, expected true")
	}
	if r6.Type() != jsonsearcher.TypeNumber {
		t.Fatalf("r6.Type is not TypeNumber, expected TypeNumber")
	}
	if r6.GetInt64() != 17 {
		t.Fatalf("r6.GetInt64() is %v, expected 17", r6.GetInt64())
	}

	r7 := s.Query("friends", 1, "name")
	if !r7.Exists() {
		t.Fatalf("r7.Exists() is false, expected true")
	}
	if r7.Type() != jsonsearcher.TypeString {
		t.Fatalf("r7.Type is not TypeString, expected TypeString")
	}
	if r7.GetString() != "Mary" {
		t.Fatalf("r7.GetString() is %v, expected Mary", r7.GetString())
	}

	r8 := s.Query("friends", 1, "age")
	if !r8.Exists() {
		t.Fatalf("r8.Exists() is false, expected true")
	}
	if r8.Type() != jsonsearcher.TypeNumber {
		t.Fatalf("r8.Type is not TypeNumber, expected TypeNumber")
	}
	if r8.GetInt64() != 18 {
		t.Fatalf("r8.GetInt64() is %v, expected 18", r8.GetInt64())
	}

	r9 := s.Query("friends", 1, "email")
	if !r9.Exists() {
		t.Fatalf("r9.Exists() is false, expected true")
	}
	if r9.Type() != jsonsearcher.TypeString {
		t.Fatalf("r9.Type is not TypeString, expected TypeString")
	}
	if r9.GetString() != "3402002560@qq.com" {
		t.Fatalf("r9.GetString() is %v, expected 3402002560@qq.com", r9.GetString())
	}

	r10 := s.Query("details")
	if !r10.Exists() {
		t.Fatalf("r10.Exists() is false, expected true")
	}
	if r10.Type() != jsonsearcher.TypeObject {
		t.Fatalf("r10.Type is not TypeObject, expected TypeObject")
	}
	if fmt.Sprintf("%v", r10.GetObject()) != "map[interests:[golang python]]" {
		t.Fatalf("r10.GetObject() is %v, expected map[interests:[golang python]]", r10.GetObject())
	}

	r11 := s.Query("details", "interests")
	if !r11.Exists() {
		t.Fatalf("r11.Exists() is false, expected true")
	}
	if r11.Type() != jsonsearcher.TypeArray {
		t.Fatalf("r11.Type is not TypeArray, expected TypeArray")
	}
	if fmt.Sprintf("%v", r11.GetArray()) != "[golang python]" {
		t.Fatalf("r11.GetArray() is %v, expected [golang python]", r11.GetArray())
	}

	r12 := s.Query("details", "interests", 0)
	if !r12.Exists() {
		t.Fatalf("r12.Exists() is false, expected true")
	}
	if r12.Type() != jsonsearcher.TypeString {
		t.Fatalf("r12.Type is not TypeString, expected TypeString")
	}
	if r12.GetString() != "golang" {
		t.Fatalf("r12.GetString() is %v, expected golang", r12.GetString())
	}

	r13 := s.Query("details", "interests", 1)
	if !r13.Exists() {
		t.Fatalf("r13.Exists() is false, expected true")
	}
	if r13.Type() != jsonsearcher.TypeString {
		t.Fatalf("r13.Type is not TypeString, expected TypeString")
	}
	if r13.GetString() != "python" {
		t.Fatalf("r13.GetString() is %v, expected python", r13.GetString())
	}

	r14 := s.Query("phone")
	if !r14.Exists() {
		t.Fatalf("r14.Exists() is false, expected true")
	}
	if r14.Type() != jsonsearcher.TypeNull {
		t.Fatalf("r14.Type is not TypeNull, expected TypeNull")
	}

	r15 := s.Query()
	if !r15.Exists() {
		t.Fatalf("r15.Exists() is false, expected true")
	}
	if r15.Type() != jsonsearcher.TypeObject {
		t.Fatalf("r15.Type is not TypeObject, expected TypeObject")
	}

	// not exists
	if s.Query("undefined").Exists() {
		t.Fatalf("the value exists, expected not exist")
	}
	if s.Query(0).Exists() {
		t.Fatalf("the value exists, expected not exist")
	}
	if s.Query(-1).Exists() {
		t.Fatalf("the value exists, expected not exist")
	}
	if s.Query("undef").Exists() {
		t.Fatalf("the value exists, expected not exist")
	}
	if s.Query("friends", -1).Exists() {
		t.Fatalf("the value exists, expected not exist")
	}
	if s.Query("friends", 2).Exists() {
		t.Fatalf("the value exists, expected not exist")
	}
	if s.Query("friends", 0, "email").Exists() {
		t.Fatalf("the value exists, expected not exist")
	}
	if s.Query("phone", 0).Exists() {
		t.Fatalf("the value exists, expected not exist")
	}
	if s.Query("phone", "undefined").Exists() {
		t.Fatalf("the value exists, expected not exist")
	}
}

func TestC(t *testing.T) {
	s, _ := jsonsearcher.New([]byte("null"))
	if !s.Query().Exists() {
		t.Fatalf("the value does not exist, expected exist")
	}
	if s.Query().Type() != jsonsearcher.TypeObject {
		t.Fatalf("the value type is not TypeObject, expected TypeObject")
	}
	if s.Query(0).Exists() {
		t.Fatalf("the value exists, expected not exist")
	}
	if s.Query("undefined").Exists() {
		t.Fatalf("the value exists, expected not exist")
	}
}
