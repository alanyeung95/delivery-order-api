package tests

import (
	"net/http"
	"testing"

	"github.com/alanyeung95/delivery-order-api/pkg/orders"
	"github.com/gavv/httpexpect"
)

//var testurl string = os.Getenv("API_TEST_DOMAIN")
var testurl string = "http://localhost:8081"

func TestRootAPI(t *testing.T) {
	e := httpexpect.New(t, testurl)
	e.GET("/").
		Expect().
		Status(http.StatusOK)
}

func TestIntegration(t *testing.T) {
	testPlaceOrderAndListOrder(t)
	testTakeOrder(t)
	testListOrder(t)
}

func testPlaceOrderAndListOrder(t *testing.T) {
	e := httpexpect.New(t, testurl)

	// create order
	postdata := map[string]interface{}{
		"origin":      []string{"22.316397", "114.264144"},
		"destination": []string{"22.307588", "114.260881"},
	}
	contentType := "application/json;charset=utf-8"

	obj := e.POST("/orders").
		WithHeader("ContentType", contentType).
		WithJSON(postdata).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("status").
		ValueEqual("status", orders.OrderStatusUnassigned)
	orderOneID := obj.Value("id").String().Raw()

	// test list order
	objs := e.GET("/orders").
		WithQuery("page", 1).
		WithQuery("limit", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	objs.Element(0).Object().ValueEqual("id", orderOneID)
	objs.Element(0).Object().ValueEqual("status", orders.OrderStatusUnassigned)
}

func testTakeOrder(t *testing.T) {
	e := httpexpect.New(t, testurl)

	// create order
	postdata := map[string]interface{}{
		"origin":      []string{"22.316397", "114.264144"},
		"destination": []string{"22.307588", "114.260881"},
	}
	contentType := "application/json;charset=utf-8"

	obj := e.POST("/orders").
		WithHeader("ContentType", contentType).
		WithJSON(postdata).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("status").
		ValueEqual("status", orders.OrderStatusUnassigned)
	orderOneID := obj.Value("id").String().Raw()

	e.PATCH("/orders/"+orderOneID).
		WithHeader("ContentType", contentType).
		WithJSON(postdata).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("status").
		ValueEqual("status", orders.OrderStatusSuccess)

	// test retake order
	e.PATCH("/orders/"+orderOneID).
		WithHeader("ContentType", contentType).
		WithJSON(postdata).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("status").
		ValueEqual("status", orders.OrderStatusTaken)
}

func testListOrder(t *testing.T) {
	e := httpexpect.New(t, testurl)

	// test list order when two records are already on the db
	objs := e.GET("/orders").
		WithQuery("page", 2).
		WithQuery("limit", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	objs.Element(0).Object().ContainsKey("status")
}
