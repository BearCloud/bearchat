package api

import(
  "github.com/go-gremlin/gremlin"
)
var DB *Client

func InitDB() {
  auth := gremlin.OptAuthUserPass("root", "root")
	DB, err = gremlin.NewClient("IP:80/gremlin", auth)
	data, err = client.ExecQuery(`g.V()`)
	if err != nil {
		panic(err)
	}
}
