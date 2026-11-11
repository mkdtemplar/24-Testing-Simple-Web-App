package main

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Check(t *testing.T) {
	form := NewForm(nil)

	form.Check(false, "password", "password is required")
	if form.Valid() {
		t.Error("Valid returns false, and should be true when calling Check()")
	}
}

func TestForm_Has(t *testing.T) {
	form := NewForm(nil)

	has := form.Has("some_field")
	if has {
		t.Error("Form shows that have a field but it should not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = NewForm(postedData)

	has = form.Has("a")
	if !has {
		t.Error("The form shows that there is no fields but it should")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/testroute", nil)
	form := NewForm(r.PostForm)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r = httptest.NewRequest("POST", "/testroute", nil)
	r.PostForm = postedData
	form = NewForm(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Form shows that has not the required fields but it does")
	}
}

func Test_errors_Get(t *testing.T) {
	form := NewForm(nil)
	form.Check(false, "password", "password is required")

	s := form.Errors.Get("password")

	if len(s) == 0 {
		t.Error("Should have an error, but do not")
	}

	s = form.Errors.Get("testfield")
	if len(s) != 0 {
		t.Error("There should not be an error, but error returned")
	}
}
