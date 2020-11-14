package api

import(
  "github.com/go-gremlin/gremlin"
)
var DB *gremlin.Client

func InitDB() {
  auth := gremlin.OptAuthUserPass("root", "root")
  var err error
  if err := gremlin.NewCluster("https://neptune-endpoint:8182/gremlin"); err != nil {
		panic(err)
	}
	DB, err = gremlin.NewClient("neptune-endpoint:8182/gremlin", auth)
	_, err = DB.ExecQuery(`g.V()`)
	if err != nil {
		panic(err)
	}
}
