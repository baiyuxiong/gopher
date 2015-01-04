/*
官方板块
*/

package gopher

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jimmykuu/wtforms"
	"gopkg.in/mgo.v2/bson"
)

// URL: /article/new
// 新建文章
func newOfficialHandler(handler Handler) {
	form := wtforms.NewForm(
		wtforms.NewHiddenField("html", ""),
		wtforms.NewTextField("title", "标题", "", wtforms.Required{}),
		wtforms.NewTextField("description", "内容", "", wtforms.Required{}),
		wtforms.NewTextField("code", "代号", "", wtforms.Required{}),
	)

	if handler.Request.Method == "POST" && form.Validate(handler.Request) {
		user, _ := currentUser(handler)

		c := handler.DB.C(OFFICIAL)

		id_ := bson.NewObjectId()

		html := form.Value("html")
		html = strings.Replace(html, "<pre>", `<pre class="prettyprint linenums">`, -1)

		err := c.Insert(&Official{
			Id_: id_,
			Content: Content{
				Id_:       id_,
				Type:      TypeOfficial,
				Title:     form.Value("title"),
				Markdown:  form.Value("description"),
				Html:      template.HTML(html),
				CreatedBy: user.Id_,
				CreatedAt: time.Now(),
			},
			OfficialCode: form.Value("code"),
		})

		if err != nil {
			fmt.Println("newOfficialHandler:", err.Error())
			return
		}

		http.Redirect(handler.ResponseWriter, handler.Request, "/official/show/"+form.Value("code"), http.StatusFound)
		return
	}

	renderTemplate(handler, "official/form.html", BASE, map[string]interface{}{
		"form":   form,
		"title":  "新建",
		"action": "/official/new",
		"active": "official",
	})
}

// URL: /official
// 首页，列出一篇文章
func officialHandler(handler Handler) {
	officialcode := "Overview"

	c := handler.DB.C(OFFICIAL)

	official := Official{}

	err := c.Find(bson.M{"officialcode": officialcode}).One(&official)

	if err != nil {
		fmt.Println("officialHandler:", err.Error())
		return
	}

	renderTemplate(handler, "official/index.html", BASE, map[string]interface{}{
		"official": official,
		"active":   "official",
	})
}

// URL: /official/{code}
// 显示文章
func showOfficialHandler(handler Handler) {
	vars := mux.Vars(handler.Request)
	officialcode := vars["officialcode"]

	c := handler.DB.C(OFFICIAL)

	official := Official{}

	err := c.Find(bson.M{"officialcode": officialcode}).One(&official)

	if err != nil {
		fmt.Println("showArticleHandler:", err.Error())
		return
	}

	renderTemplate(handler, "official/index.html", BASE, map[string]interface{}{
		"official": official,
		"active":   "official",
	})
}

// URL: /a/{articleId}/edit
// 编辑主题
func editOfficialHandler(handler Handler) {

}

// URL: /a/{articleId}/delete
// 删除文章
func deleteOfficialHandler(handler Handler) {

}
