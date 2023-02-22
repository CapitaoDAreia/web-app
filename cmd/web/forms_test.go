package main

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestFormHas(t *testing.T) {

	emptyPostedData := url.Values{}

	tests := []struct {
		testName   string
		fieldValue string
		formValue  url.Values
		inputKey   string
		inputValue string
	}{
		{
			testName:   "When form have fields",
			fieldValue: "",
			formValue:  emptyPostedData,
			inputKey:   "b",
			inputValue: "a",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			//if form haven't a specific field
			form := NewForm(test.formValue)
			has := form.Has(test.fieldValue)
			if has {
				t.Error("There is no fields in form.")
			}

			// if form has a specific field
			postedData := url.Values{}
			form = NewForm(postedData)
			postedData.Add(test.inputKey, test.inputValue)
			has = form.Has(test.inputKey)
			if !has {
				t.Error("There should be fields.")
			}
		})
	}
}

// TODO: Refactor this function
func TestFormRequired(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := NewForm(r.PostForm)

	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("Form shows valid when required fields are missing.")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	nr := httptest.NewRequest("POST", "/whatever", nil)
	nr.PostForm = postedData
	form = NewForm(nr.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Form shows invalid when required fields are there.")
	}
}
