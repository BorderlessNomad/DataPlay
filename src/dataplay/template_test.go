package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"testing"
)

func TestApplyTemplate(t *testing.T) {
	response := httptest.NewRecorder()

	Convey("On page template apply", t, func() {
		ApplyTemplate("test", "test", response)
	})
}

func TestRenderTemplate(t *testing.T) {
	response := httptest.NewRecorder()
	custom := map[string]string{"test" : "test"}


	Convey("On page template render", t, func() {
		RenderTemplate("test", custom, response)
	})

}
