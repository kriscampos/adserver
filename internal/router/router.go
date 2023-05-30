package router

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kriscampos/adserver/internal/ad_engine"
	"github.com/kriscampos/adserver/internal/campaign"
)

type postAdDecisionRequest struct {
	Keywords []string `json:"keywords" binding:"required"`
}

type router struct {
	campaignService *campaign.CampaignService
	adEngine        *ad_engine.AdEngine
}

func newRouter(engine *ad_engine.AdEngine) *router {
	return &router{
		campaignService: campaign.NewCampaignService(),
		adEngine:        engine,
	}
}

func SetupRouter(adEngine *ad_engine.AdEngine) *gin.Engine {
	router := gin.Default()

	// Middleware goes here

	handler := newRouter(adEngine)
	router.POST("/campaign", handler.PostCampaign)
	router.POST("/addecision", handler.PostAdDecision)
	router.GET("/:impression-url", handler.GetImpressionURL)

	return router
}

func (r *router) PostCampaign(ctx *gin.Context) {
	var postCampaignRequest campaign.PostCampaignRequest
	if err := ctx.BindJSON(&postCampaignRequest); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	newCampaign := r.campaignService.CreateCampaign(&postCampaignRequest)
	r.adEngine.RegisterCampaign(newCampaign)
	responseData := gin.H{
		"campaign_id": newCampaign.ID,
	}
	ctx.IndentedJSON(http.StatusOK, responseData)
}

func (r *router) PostAdDecision(ctx *gin.Context) {
	var newAdDecisionRequest postAdDecisionRequest
	if err := ctx.BindJSON(&newAdDecisionRequest); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	campaign, ok := r.adEngine.RecommendCampaign(newAdDecisionRequest.Keywords)
	if !ok {
		return // returns status 200
	}
	responseData := gin.H{
		"campaign_id":    campaign.ID,
		"impression_url": campaign.ImpressionURL,
	}
	ctx.IndentedJSON(http.StatusOK, responseData)
}

func (r *router) GetImpressionURL(ctx *gin.Context) {
	impressionURL := ctx.Param("impression-url")
	log.Printf("Impression URL: %s\n", impressionURL)
	reachedMax, validURL := r.campaignService.IncrementImpression(impressionURL)
	if !validURL {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	if reachedMax {
		log.Printf("Reached Max!\n")
		r.adEngine.DeleteCampaign(impressionURL)
	}
}
