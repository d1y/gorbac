package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/storyicon/grbac"
)

func QueryRolesByHeaders(ctx *context.Context) (roles []string, err error) {
	role := ctx.URLParamDefault("role", "")
	if role == "superadmin" {
		return []string{"superadmin"}, nil
	} else if len(role) >= 1 {
		return []string{role}, nil
	}
	return
}

func Authentication() iris.Handler {
	rbac, err := grbac.New(grbac.WithYAML("auth_rbac.yaml", time.Second*12))
	if err != nil {
		panic(err)
	}
	return func(c *context.Context) {
		roles, err := QueryRolesByHeaders(c)
		if err != nil {
			c.StatusCode(http.StatusInternalServerError)
			c.StopExecution()
			return
		}
		state, err := rbac.IsRequestGranted(c.Request(), roles)
		if err != nil {
			c.StatusCode(http.StatusInternalServerError)
			c.StopExecution()
			return
		}
		if !state.IsGranted() {
			c.StatusCode(http.StatusUnauthorized)
			c.StopExecution()
			return
		}
		c.Next()
	}
}

func main() {
	c := iris.New()
	c.Use(Authentication())
	c.Get("/", func(c iris.Context) {
		c.HTML("<h1>Hello World!</h1>")
	})
	c.PartyFunc("/public", func(u iris.Party) {
		u.Get("/1.txt", func(c iris.Context) {
			s := uuid.NewString()
			c.Text(s)
		})
	})
	c.PartyFunc("/api", func(u iris.Party) {
		u.Get("/user/list", func(c iris.Context) {
			c.Text("OK")
		})
		u.PartyFunc("/user", func(p iris.Party) {
			p.Get("/admin/list", func(c iris.Context) {
				s := uuid.NewString()
				c.Text(s)
			})
		})
	})
	c.Listen(":8080")
}
