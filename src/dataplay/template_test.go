package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	response := httptest.NewRecorder()
	custom := map[string]string{"test" : "test"}


	Convey("On page template render", t, func() {
		RenderTemplate("test", custom, response)
	})

}
