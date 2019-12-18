package tasgrid

import (
	"flag"
	"os"
	"testing"
)

var tasGrid *MapGrid
var gridPoint GridPoint
var err error

func TestMain(m *testing.M) {

	tasGrid = NewTasMapGrid()
	gridPoint, err = NewGridPoint("GOO", "545", "519", *tasGrid)

	flag.Parse()
	os.Exit(m.Run())
}

func TestGetFullEasting(t *testing.T) {
	got := gridPoint.GetFullEasting()
	want := "554500"

	if got != want {
		t.Errorf("Function GetFullEasting failed. Expected %v, instead got %v", want, got)
	}
}

func TestGetFullNorthing(t *testing.T) {
	got := gridPoint.GetFullNorthing()
	want := "5551900"

	if got != want {
		t.Errorf("Function GetFullNorthing failed. Expected %v, instead got %v", want, got)
	}
}

func TestGetDecimalLat(t *testing.T) {
	got := gridPoint.GetDecimalLat()
	want := "-40.181512"

	if got != want {
		t.Errorf("Function GetDecimalLat failed. Expected %v, instead got %v", want, got)
	}
}

func TestGetDecimalLong(t *testing.T) {
	got := gridPoint.GetDecimalLong()
	want := "147.640171"

	if got != want {
		t.Errorf("Function GetDecimalLong failed. Expected %v, instead got %v", want, got)
	}
}

func TestGetLatSeconds(t *testing.T) {
	got := gridPoint.GetLatSeconds()
	want := "53.4"

	if got != want {
		t.Errorf("Function GetLatSeconds failed. Expected %v, instead got %v", want, got)
	}
}

func TestGetLongSeconds(t *testing.T) {
	got := gridPoint.GetLongSeconds()
	want := "24.6"

	if got != want {
		t.Errorf("Function GetLongSeconds failed. Expected %v, instead got %v", want, got)
	}
}

func TestGetDistance(t *testing.T) {
	got, _ := gridPoint.GetDistance("-41.432563", "145.234567")
	want := 245878.67573113748

	if got != want {
		t.Errorf("Function GetDistance failed. Expected %v, instead got %v", want, got)
	}
}

func TestDDtoDMS(t *testing.T) {
	got1, got2, got3 := ddToDMS(-42.123456)
	want1, want2, want3 := -42.0, 7.0, 24.441599999990444

	if got1 != want1 || got2 != want2 || got3 != want3 {
		t.Errorf("Function ddToDMS failed. Expected %v, %v, %v, instead got %v, %v, %v", want1, want2, want3, got1, got2, got3)
	}
}
