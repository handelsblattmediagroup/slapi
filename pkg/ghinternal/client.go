package ghinternal

import (
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v43/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog/log"
	"net/http"
	"serenitylabs.cloud/slapi/pkg/api"
)

var llog = log.With().Str("system", "github").Logger()

func Provider(config *api.Config) (api.VersionedRouterOut, error) {
	c := githubapp.Config{}

	c.V3APIURL = "https://api.github.com/"
	c.App.IntegrationID = config.Github.IntegrationID
	c.App.PrivateKey = config.Github.PrivateKey

	cc, err := githubapp.NewDefaultCachingClientCreator(c)

	if err != nil {
		return api.VersionedRouterOut{}, err
	}

	routerGroupFunc := func(group *gin.RouterGroup) {
		group.GET("/installations", func(c *gin.Context) {
			appClient, err := cc.NewAppClient()
			if err != nil {
				c.Error(err)
				return
			}
			installations, _, err := appClient.Apps.ListInstallations(c.Request.Context(), &github.ListOptions{})
			if err != nil {
				c.Error(err)
				return
			}
			c.JSON(http.StatusOK, struct {
				Installations []*github.Installation `json:"installations"`
			}{installations})
		})
	}

	return api.VersionedRouterOut{
		Router: &api.VersionedRouterSpec{
			Version:        0,
			Prefix:         "github",
			RegisterRouter: routerGroupFunc,
		},
	}, nil
}
