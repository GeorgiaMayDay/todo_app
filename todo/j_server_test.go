package todo

// func TestServer(t *testing.T) {
// 	t.Run("happy path test", func(t *testing.T) {
// 		data := "hello, world"
// 		store := &TodoList{t, data}
// 		svr := Server(store)

// 		request := httptest.NewRequest(http.MethodGet, "/", nil)
// 		response := httptest.NewRecorder()

// 		svr.ServeHTTP(response, request)

// 		if response.Body.String() != data {
// 			t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
// 		}
// 	})
// }
