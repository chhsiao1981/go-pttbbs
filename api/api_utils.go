package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
)

func processResult(c *gin.Context, result interface{}, err error) {
	setHeader(c)

	switch err {
	case nil:
		c.JSON(200, result)

	case ErrInvalidHost:
		c.JSON(400, &ErrResult{err.Error()})
	case ErrInvalidRemoteAddr:
		c.JSON(400, &ErrResult{err.Error()})

	case ErrInvalidParams:
		c.JSON(400, &ErrResult{err.Error()})
	case ErrInvalidPath:
		c.JSON(400, &ErrResult{err.Error()})

	case bbs.ErrInvalidParams:
		c.JSON(400, &ErrResult{err.Error()})

	case ptttype.ErrUserIDAlreadyExists:
		c.JSON(400, &ErrResult{err.Error()})

	case ErrInvalidIDEmail:
		c.JSON(400, &ErrResult{err.Error()})

	// 401
	case ErrInvalidToken:
		c.JSON(401, &ErrResult{err.Error()})

	case ErrInvalidToken:
		c.JSON(401, &ErrResult{err.Error()})
	case ErrLoginFailed:
		c.JSON(401, &ErrResult{err.Error()})

	// 403
	case cache.ErrInvalidUID:
		c.JSON(403, &ErrResult{err.Error()})
	case ErrInvalidUser:
		c.JSON(403, &ErrResult{err.Error()})
	case ptttype.ErrInvalidUserID:
		c.JSON(403, &ErrResult{err.Error()})
	case ptt.ErrNotPermitted:
		c.JSON(403, &ErrResult{err.Error()})
	case ptt.ErrPermitNoPost:
		c.JSON(403, &ErrResult{err.Error()})
	case ptt.ErrBanned:
		c.JSON(403, &ErrResult{err.Error()})
	case ptt.ErrRestricted:
		c.JSON(403, &ErrResult{err.Error()})
	case ptt.ErrViolateLaw:
		c.JSON(403, &ErrResult{err.Error()})
	case bbs.ErrInvalidPermission:
		c.JSON(403, &ErrResult{err.Error()})
	default:
		c.JSON(500, &ErrResult{err.Error()})
	}
}

func setHeader(c *gin.Context) {
	if !types.IS_ALLOW_CROSSDOMAIN {
		return
	}

	origin := c.GetHeader("Origin")

	if origin == "" {
		return
	}

	requestHeaders := c.GetHeader("Access-Control-Request-Headers")

	c.Header("X-Frame-Options", "SAMEORIGIN")
	c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Origin", origin)
	if requestHeaders != "" {
		c.Header("Access-Control-Allow-Headers", requestHeaders)
	}
}
