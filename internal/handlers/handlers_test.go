package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Sergio-dot/bookings/internal/models"
)

/* type postData struct {
	key   string
	value string
} */

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"non-existent", "/green/eggs/and/ham", "GET", http.StatusNotFound},
	// new routes
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"new res", "/admin/reservations-new", "GET", http.StatusOK},
	{"all res", "/admin/reservations-all", "GET", http.StatusOK},
	{"show res", "/admin/reservations/new/1/show", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Fatalf("Expected status code %d, got %d", e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned %d, expected %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned %d, expected %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	/* 	reqBody := "start_date=01/01/2050"
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02/01/2050")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1") */

	postedData := url.Values{}
	postedData.Add("start_date", "01/01/2050")
	postedData.Add("end_date", "02/01/2050")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "555-555-1234")
	postedData.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid start date
	/* 	reqBody = "start_date=invalid"
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02/01/2050")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1") */

	postedData = url.Values{}
	postedData.Add("start_date", "invalid")
	postedData.Add("end_date", "02/01/2050")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "555-555-1234")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid end date
	/* 	reqBody = "start_date=01/01/2050"
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1") */

	postedData = url.Values{}
	postedData.Add("start_date", "01/01/2050")
	postedData.Add("end_date", "invalid")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "555-555-1234")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid roomID
	/* 	reqBody = "start_date=01/01/2050"
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02/01/2050")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid") */

	postedData = url.Values{}
	postedData.Add("start_date", "01/01/2050")
	postedData.Add("end_date", "02/01/2050")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "555-555-1234")
	postedData.Add("room_id", "invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid data
	/* 	reqBody = "start_date=01/01/2050"
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02/01/2050")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=J")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1") */

	postedData = url.Values{}
	postedData.Add("start_date", "01/01/2050")
	postedData.Add("end_date", "02/01/2050")
	postedData.Add("first_name", "J")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "555-555-1234")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostReservation handler returned %d, expected %d", rr.Code, http.StatusOK)
	}

	// test for failure to insert reservation into database
	/* 	reqBody = "start_date=01/01/2050"
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02/01/2050")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2") */

	postedData = url.Values{}
	postedData.Add("start_date", "01/01/2050")
	postedData.Add("end_date", "02/01/2050")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "555-555-1234")
	postedData.Add("room_id", "2")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// test for failure to insert restriction into database
	/* 	reqBody = "start_date=01/01/2050"
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02/01/2050")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	   	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000") */

	postedData = url.Values{}
	postedData.Add("start_date", "01/01/2050")
	postedData.Add("end_date", "02/01/2050")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "555-555-1234")
	postedData.Add("room_id", "1000")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned %d, expected %d", rr.Code, http.StatusSeeOther)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
