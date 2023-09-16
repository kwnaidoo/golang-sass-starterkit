package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/noirbizarre/gonja"
	"gorm.io/gorm"
	"plexcorp.tech/gosass/models"
)

type NinjaRender string

func (n NinjaRender) Render(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write([]byte(n))
	return err
}

func (n NinjaRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

type Controller struct {
}

func (c *Controller) GetDB(gctx *gin.Context) *gorm.DB {
	return gctx.MustGet("db").(*gorm.DB)
}

func (c *Controller) RenderHtml(tpl_name string, ctx gonja.Context, gctx *gin.Context, layoutTpl string) {

	session := sessions.Default(gctx)
	flashes := session.Flashes("success")
	if len(flashes) == 1 {
		ctx["successMsg"] = flashes[0].(string)
	}

	errorMessage := session.Flashes("error")
	if len(errorMessage) == 1 {
		errors, _ := ctx["errrors"].([]string)
		ctx["errors"] = append(errors, errorMessage[0].(string))
	}

	if len(errorMessage) > 0 || len(flashes) > 0 {
		session.Save()
	}

	_, ok := ctx["highlight"]
	if !ok {
		ctx["highlight"] = ""
	}

	ctx["gosass_base_url"] = os.Getenv("gosass_URL")
	ctx["STATUS_QUEUED"] = models.STATUS_QUEUED
	ctx["STATUS_RUNNING"] = models.STATUS_RUNNING
	ctx["STATUS_FAILED"] = models.STATUS_FAILED
	ctx["STATUS_CONNECTING"] = models.STATUS_CONNECTING
	ctx["STATUS_COMPLETE"] = models.STATUS_COMPLETE

	ctx["_csrf_token"] = c.SetAndGetCSRFToken(gctx)

	view, err := gonja.Must(gonja.FromFile("templates/" + tpl_name + ".jinja")).Execute(ctx)
	if err != nil {
		fmt.Println(err)
	}

	ctx["view"] = view

	var MASTER_TPL = gonja.Must(gonja.FromFile(layoutTpl + ".jinja"))
	tpl, err := MASTER_TPL.Execute(ctx)
	if err != nil {
		fmt.Println(err)
	}

	gctx.Render(http.StatusOK, NinjaRender(tpl))
}

func (c *Controller) Render(tpl_name string, ctx gonja.Context, gctx *gin.Context) {
	c.RenderHtml(tpl_name, ctx, gctx, "templates/master")
}

func (c *Controller) RenderAuth(tpl_name string, ctx gonja.Context, gctx *gin.Context) {
	c.RenderHtml(tpl_name, ctx, gctx, "templates/auth")
}

func (c *Controller) RenderWithoutLayout(tpl_name string, ctx gonja.Context, gctx *gin.Context) {
	ctx["gosass_base_url"] = os.Getenv("gosass_URL")
	view, err := gonja.Must(gonja.FromFile("templates/" + tpl_name + ".jinja")).Execute(ctx)
	if err != nil {
		fmt.Println(err)
	}
	gctx.Render(http.StatusOK, NinjaRender(view))
}

func (c *Controller) FlashSuccess(gctx *gin.Context, msg string) {
	session := sessions.Default(gctx)
	session.AddFlash(msg, "success")
	session.Save()
}

func isTimestampWithin5Minutes(timestampStr string) bool {
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return true
	}

	timestampTime := time.Unix(timestamp, 0)
	timeDiff := time.Since(timestampTime)
	return timeDiff >= 2*time.Minute
}

func (c *Controller) SetAndGetCSRFToken(gctx *gin.Context) string {
	session := sessions.Default(gctx)
	sessSet := session.Get("csrfToken")
	setToken := true
	token := ""

	if sessSet != nil && sessSet.(string) != "" {

		token = sessSet.(string)
		tStamp := strings.Split(sessSet.(string), "|")
		if len(tStamp) == 2 {
			tStampStr := tStamp[1]
			setToken = isTimestampWithin5Minutes(tStampStr)
		}
	}

	if setToken {
		token := fmt.Sprintf("%s|%d", uuid.New().String(), time.Now().Unix())
		session.Set("csrfToken", token)
		session.Save()
	}

	return token
}

func (c *Controller) TestCSRFToken(gctx *gin.Context) bool {
	token := gctx.PostForm("_csrf_token")

	session := sessions.Default(gctx)
	return session.Get("csrfToken").(string) == token
}

func (c *Controller) FlashError(gctx *gin.Context, msg string) {
	session := sessions.Default(gctx)
	session.AddFlash(msg, "error")
	session.Save()
}

func (c *Controller) GetSessionUser(gctx *gin.Context) models.User {
	session := sessions.Default(gctx)
	userId := session.Get("user_id").(int64)
	var user models.User
	c.GetDB(gctx).Raw("SELECT id, name, email, verified, team_id FROM users where id = ?", userId).Scan(&user)
	return user
}

func (c *Controller) ShowGuide(gctx *gin.Context) {
	c.Render("general/guide", gonja.Context{
		"title":     "gosass help guide",
		"highlight": "help",
	}, gctx)
}

func (c *Controller) AccessDenied(gctx *gin.Context) {
	c.Render("general/permission", gonja.Context{
		"highlight": "",
	}, gctx)
}
