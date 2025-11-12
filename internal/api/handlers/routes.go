package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

var (
	teamAddPostfix          = "/team/add"
	teamGetPostfix          = "/team/get"
	usersSetIsActivePostfix = "/users/setIsActive"
	prCreatePostfix         = "/pullRequest/create"
	prMergePostfix          = "/pullRequest/merge"
	prReassignPostfix       = "/pullRequest/reassign"
	usersGetReviewPostfix   = "/users/getReview"
)

func Register(router *gin.Engine, db *sql.DB) {
	// handler := NewHandler(db)

	router.POST(teamAddPostfix)
	router.GET(teamGetPostfix)
	router.GET(usersSetIsActivePostfix)
	router.POST(prCreatePostfix)
	router.POST(prMergePostfix)
	router.POST(prReassignPostfix)
	router.GET(usersGetReviewPostfix)
}
