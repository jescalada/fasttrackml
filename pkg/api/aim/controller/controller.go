package controller

import (
	"github.com/G-Research/fasttrackml/pkg/api/aim/services/app"
	"github.com/G-Research/fasttrackml/pkg/api/aim/services/dashboard"
	"github.com/G-Research/fasttrackml/pkg/api/aim/services/experiment"
	"github.com/G-Research/fasttrackml/pkg/api/aim/services/project"
	"github.com/G-Research/fasttrackml/pkg/api/aim/services/run"
	"github.com/G-Research/fasttrackml/pkg/api/aim/services/tag"
	"github.com/G-Research/fasttrackml/pkg/common/services/artifact"
)

// Controller handles all the input HTTP requests.
type Controller struct {
	tagService        *tag.Service
	appService        *app.Service
	runService        *run.Service
	artifactService   *artifact.Service
	projectService    *project.Service
	dashboardService  *dashboard.Service
	experimentService *experiment.Service
}

// NewController creates new Controller instance.
func NewController(
	tagService *tag.Service,
	appService *app.Service,
	runService *run.Service,
	artifactService *artifact.Service,
	projectService *project.Service,
	dashboardService *dashboard.Service,
	experimentService *experiment.Service,
) *Controller {
	return &Controller{
		tagService:        tagService,
		appService:        appService,
		runService:        runService,
		artifactService:   artifactService,
		projectService:    projectService,
		dashboardService:  dashboardService,
		experimentService: experimentService,
	}
}
