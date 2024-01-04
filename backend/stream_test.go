package stream_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"stream/models"
	"stream/routes"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"nhooyr.io/websocket"
)

func SetupRouterForTest() *echo.Echo {
	e := echo.New()

	// registers all the available routes
	routes.Init(e)
	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetOutput(os.Stdout)

	// starts echo server

	return e
}

const EAddr = "127.0.0.1:8082"

func TestParallelWebscket(t *testing.T) {
	// t.Skip()
	// Create test server with the echo handler.
	router := SetupRouterForTest()
	s := httptest.NewServer(router.Server.Handler)
	//go router.Start(EAddr)
	//defer router.Close()

	<-time.After(time.Second * 2)
	StepGetCpu(t, s)
}

func StepGetCpu(t *testing.T, s *httptest.Server) {
	// Connect to the server
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	//u := "ws://" + EAddr + "/stream"
	u := s.URL + "/stream"
	t.Logf("websocket-URL:%v", u)

	opt := &websocket.DialOptions{}
	c, _, err := websocket.Dial(ctx, u, opt)
	if err != nil {
		t.Fatal(err)
	}
	defer c.CloseNow()

	// mt, r, err := c.Reader(ctx)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Logf("Message type is %v", mt)

	for {
		// var reply models.Cpu
		// dec := json.NewDecoder(r)
		// err := dec.Decode(&reply)
		_, buf, err := c.Read(context.Background())
		if err == nil {
			var reply models.Cpu
			err := json.Unmarshal(buf, &reply)
			if err != nil {
				t.Logf("Unmarshal error: %v", err)
			} else {
				t.Logf("Get Reply: %v", reply)
			}
		}
		//<-time.After(time.Second)
		if err == io.EOF {
			t.Logf("Received EOF: %v", err)
			return
		} else if err != nil {
			t.Fatalf("FAIL: Invalid Response: %v", err)
		}
	}

}

func TestParallelWebscketReader(t *testing.T) {
	// t.Skip()
	// Create test server with the echo handler.
	router := SetupRouterForTest()
	s := httptest.NewServer(router.Server.Handler)
	//go router.Start(EAddr)
	//defer router.Close()

	<-time.After(time.Second * 2)
	StepGetCpuReader(t, s)
}

func StepGetCpuReader(t *testing.T, s *httptest.Server) {
	// Connect to the server
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	//u := "ws://" + EAddr + "/stream"
	u := s.URL + "/stream"
	t.Logf("websocket-URL:%v", u)

	opt := &websocket.DialOptions{}
	c, _, err := websocket.Dial(ctx, u, opt)
	if err != nil {
		t.Fatal(err)
	}
	defer c.CloseNow()

	mt, r, err := c.Reader(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Message type is %v", mt)

	for {
		var reply models.Cpu
		dec := json.NewDecoder(r)
		err := dec.Decode(&reply)
		if err == nil {
			if err != nil {
				t.Logf("Unmarshal error: %v", err)
			} else {
				t.Logf("Get Reply: %v", reply)
			}
		}
		//<-time.After(time.Second)
		if err == io.EOF {
			t.Logf("Received EOF: %v", err)
			return
		} else if err != nil {
			t.Fatalf("FAIL: Invalid Response: %v", err)
		}
	}
}
