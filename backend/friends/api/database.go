package api

import(
  "github.com/schwartzmx/gremtune"
)

var gremlinClient *gremtune.Client 

func InitDB() {
  dialer := gremtune.NewDialer("https://neptune-endpoint:8182/gremlin")
  gremlinClient, err := gremtune.Dial(dialer, nil)
  if err != nil {
		panic(err)
	}
  _, err = gremlinClient.Execute("g.V()")
	if err != nil {
		panic(err)
	}
}
