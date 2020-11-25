package api

import(
  "github.com/furtiaga/gremlin"
)

func InitDB() {
  if err := gremlin.NewCluster("https://neptune-endpoint:8182/gremlin"); err != nil {
		panic(err)
	}
	_, err = gremlin.Query(`g.V()`).Exec()
	if err != nil {
		panic(err)
	}
}
