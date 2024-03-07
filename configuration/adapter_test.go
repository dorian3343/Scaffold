package configuration

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
	"github.com/metalim/jsonmap"
	"reflect"
	"service/components/controller"
	model2 "service/components/model"
	"testing"
)

func TestModelAdaptModel(t *testing.T) {
	// Mocking the database connection
	db, err := sql.Open("sqlite", "")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create a sample model
	sampleModel := model{
		QueryTemplate: "SELECT * FROM table_name WHERE column_name = ?",
		JsonTemplate: []JsonSpecSkeleton{
			{Type: "string", Name: "name"},
			{Type: "int", Name: "age"},
		},
		Name: "sample_model",
	}
	x, err := sampleModel.adapt(db)

	// Check if error is nil
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	f := jsonmap.New()
	for i := 0; i < len(sampleModel.JsonTemplate); i++ {
		f.Set(sampleModel.JsonTemplate[i].Name, sampleModel.JsonTemplate[i].Type)
	}

	sampleCase := model2.Create("sample_model", db, "SELECT * FROM table_name WHERE column_name = ?", f)
	if !reflect.DeepEqual(x, sampleCase) {
		t.Errorf("Adapted controller does not match expected controller. Expected: %+v, Got: %+v", sampleCase, x)
	}
}

func TestModelAdaptController(t *testing.T) {

	sample := Controller{Name: "name", Fallback: "ok", Model: "", Cors: "*", Cache: ""}
	x, err := sample.adapt(nil)
	// Check if error is nil
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// Create an expected controller
	expected := controller.Create("name", nil, []byte(`"ok"`), "*", "")

	// Check if the adapted model is equal to the expected model using reflection
	if !reflect.DeepEqual(x, expected) {
		t.Errorf("Adapted controller does not match expected controller. Expected: %+v, Got: %+v", expected, sample)
	}
}
