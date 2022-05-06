package hero

// import (
// 	"io/ioutil"

// 	"github.com/gin-gonic/gin"
// )

// func Router(router *gin.Engine) {
// 	router.POST("/hero", routerHeroPatch)
// 	router.PATCH("/hero", routerHeroPatch)
// 	router.GET("/heros", routerHeroList)
// 	router.GET("/hero/:uuid", routerHeroGet)
// }

// // 新增或更新
// func routerHeroPatch(ctx *gin.Context) {
// 	v, err := ioutil.ReadAll(ctx.Request.Body)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 從middleware取得相關資訊
// 	tags := q.Tags{"OpType": "test", "OpID": "123-123-123"}

// 	h := &Hero{}
// 	q.STRUCT(v, h)

// 	result := HeroPatch(tags, h)
// 	ctx.String(200, result.String())
// }

// // 列表
// func routerHeroList(ctx *gin.Context) {
// 	hq := &HeroQuery{}
// 	ctx.BindQuery(hq)
// 	result := HeroList(hq)
// 	ctx.String(200, result.String())
// }

// // 取得一個
// func routerHeroGet(ctx *gin.Context) {
// 	uuid := ctx.Param("uuid")
// 	h := &Hero{UUID: uuid}
// 	result := HeroGet(h)
// 	ctx.String(200, result.String())
// }
