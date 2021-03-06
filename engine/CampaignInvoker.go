package engine

import (
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"gopkg.in/vmihailenco/msgpack.v2"
	"time"
)

func invokeCampaign(demand *Demand) {

	adRequest := demand.AdRequest
	beego.Debug("Start Invoke Campaign,bid:" + adRequest.Bid)

	adResponse := getCachedAdResponse(demand)

	if adResponse == nil {
		adResponse = initAdResponse(demand)

		campaigns, err := m.GetCampaigns(demand.AdspaceKey, time.Now().Format("2006-01-02"), demand.AdRequest.Width, demand.AdRequest.Height)
		if err != nil {
			beego.Error(err.Error)
			adResponse.StatusCode = lib.ERROR_CAMPAIGN_DB_ERROR
		} else {
			if len(campaigns) == 0 {
				adResponse.StatusCode = lib.ERROR_NOAD
			} else {
				random := lib.GetRandomNumber(0, len(campaigns))
				mapCampaign(adResponse, campaigns[random])
				setCachedAdResponse(generateCacheKey(demand), adResponse)
			}
		}
	}

	//reset cache data for the refreshed data
	adResponse.Bid = adRequest.Bid
	adResponse.Did = adRequest.Did
	adResponse.SetResponseTime(time.Now().Unix())

	go SendDemandLog(adResponse)
	demand.Result <- adResponse

}

func mapCampaign(adResponse *m.AdResponse, campaign *m.PmpCampaignCreative) {

	adResponse.StatusCode = lib.STATUS_SUCCESS

	adUnit := new(m.AdUnit)
	adResponse.Adunit = adUnit
	adUnit.Cid = lib.ConvertIntToString(campaign.Id)
	adUnit.ClickUrl = campaign.LandingUrl
	adUnit.CreativeUrls = []string{campaign.CreativeUrl}
	adUnit.AdWidth = campaign.Width
	adUnit.AdHeight = campaign.Height

}

func generateCacheKey(demand *Demand) string {
	strWidth := lib.ConvertIntToString(demand.AdRequest.Width)
	strHeight := lib.ConvertIntToString(demand.AdRequest.Height)
	return beego.AppConfig.String("runmode") + "_CAMPAIGN_" + demand.AdRequest.AdspaceKey + "_" + demand.AdspaceKey + "_" + strWidth + "_" + strHeight
}

func setCachedAdResponse(cacheKey string, adResponse *m.AdResponse) {
	c := lib.GetCachePool().Get()
	val, err := msgpack.Marshal(adResponse)

	if _, err = c.Do("SETEX", cacheKey, 1800, val); err != nil {
		beego.Error(err.Error())
	}

	//_, err = c.Do("EXPIRE", cacheKey, 1800)
	//if err != nil {
	//	beego.Error(err.Error())
	//}
}

func getCachedAdResponse(demand *Demand) (adResponse *m.AdResponse) {
	c := lib.GetCachePool().Get()

	key := generateCacheKey(demand)
	v, err := c.Do("GET", key)
	if err != nil {
		beego.Error(err.Error())
		return nil
	}

	if v == nil {
		return
	}

	adResponse = new(m.AdResponse)
	switch t := v.(type) {
	case []byte:
		err = msgpack.Unmarshal(t, adResponse)
	default:
		err = msgpack.Unmarshal(t.([]byte), adResponse)
	}

	if err != nil {
		beego.Error(err.Error())
	}
	return
}
