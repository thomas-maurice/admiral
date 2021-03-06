package api

import (
	"net/http"
	"strings"

	"github.com/Devatoria/admiral/auth"
	"github.com/Devatoria/admiral/models"
	"github.com/Devatoria/admiral/token"

	"github.com/gin-gonic/gin"
)

// getToken returns a JWT bearer token to the registry containing the user accesses
func getToken(c *gin.Context) {
	service := c.Query("service")
	user, err := auth.GetCurrentUser(c)
	if err != nil {
		panic(err)
	}

	// Scope is empty only for authentication
	var claimsAccesses []token.ClaimsAccess
	scope := c.Query("scope")
	if scope != "" {
		// Parse scope: repository:samalba/my-app:pull,push
		scopeSplit := strings.SplitN(scope, ":", 3)
		if len(scopeSplit) != 3 {
			c.Status(http.StatusUnauthorized)
			return
		}

		switch scopeSplit[0] {
		case "repository":
			// Parse image name to retrieve namespace
			imageName := scopeSplit[1]
			imageSplit := strings.SplitN(imageName, "/", 2)
			if len(imageSplit) != 2 {
				c.Status(http.StatusUnauthorized)
				return
			}

			image := models.GetImageByName(imageName)
			nsName := imageSplit[0]
			var accesses []string

			// If image is part of the user namespace, give all accesses
			// If image is public and not part of the user namespace, give pull only
			if nsName == user.Username {
				accesses = append(accesses, "*")
			}
			if image.IsPublic {
				accesses = append(accesses, "pull")
			}

			claimsAccesses = append(claimsAccesses, token.ClaimsAccess{
				Type:    "repository",
				Name:    imageName,
				Actions: accesses,
			})
		}
	}

	t := token.NewToken(service, user.Username, claimsAccesses)
	tString, err := token.SignToken(t)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"token": tString})
}
