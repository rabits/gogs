// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package v1

import (
	"github.com/Unknwon/com"

	api "github.com/gogits/go-gogs-client"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/modules/base"
	"github.com/gogits/gogs/modules/middleware"
	"github.com/gogits/gogs/modules/setting"
)

// ToApiUser converts user to API format.
func ToApiUser(u *models.User) *api.User {
	return &api.User{
		Id:        u.Id,
		UserName:  u.Name,
		AvatarUrl: string(setting.Protocol) + u.AvatarLink(),
	}
}

func SearchUsers(ctx *middleware.Context) {
	opt := models.SearchOption{
		Keyword: ctx.Query("q"),
		Limit:   com.StrTo(ctx.Query("limit")).MustInt(),
	}
	if opt.Limit == 0 {
		opt.Limit = 10
	}

	us, err := models.SearchUserByName(opt)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}

	results := make([]*api.User, len(us))
	for i := range us {
		results[i] = &api.User{
			UserName:  us[i].Name,
			AvatarUrl: us[i].AvatarLink(),
			FullName:  us[i].FullName,
		}
	}

	ctx.Render.JSON(200, map[string]interface{}{
		"ok":   true,
		"data": results,
	})
}

// GET /users/:username
func GetUserInfo(ctx *middleware.Context) {
	u, err := models.GetUserByName(ctx.Params(":username"))
	if err != nil {
		if err == models.ErrUserNotExist {
			ctx.Error(404)
		} else {
			ctx.JSON(500, &base.ApiJsonErr{"GetUserByName: " + err.Error(), base.DOC_URL})
		}
		return
	}
	ctx.JSON(200, &api.User{u.Id, u.Name, u.FullName, u.Email, u.AvatarLink()})
}
