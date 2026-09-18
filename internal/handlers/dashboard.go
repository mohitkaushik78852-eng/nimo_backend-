package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Household struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Quality struct {
	TDS             int     `json:"tds"`
	TDSStatus       string  `json:"tds_status"`
	Turbidity       float64 `json:"turbidity"`
	TurbidityStatus string  `json:"turbidity_status"`
}

type Motor struct {
	Status string `json:"status"`
	Mode   string `json:"mode"`
}

type Alert struct {
	Type     string `json:"type"`
	Severity string `json:"severity"`
	Time     string `json:"time"`
}

type DashboardData struct {
	Online            bool    `json:"online"`
	WaterSupply       string  `json:"water_supply"`
	TankLevel         int     `json:"tank_level"`
	FlowLpm           float64 `json:"flow_lpm"`
	Motor             Motor   `json:"motor"`
	Quality           Quality `json:"quality"`
	ConsumptionTodayL int     `json:"consumption_today_l"`
	LastAlert         Alert   `json:"last_alert"`
}

var households = []Household{
	{ID: "hh_1", Name: "Kumar Residence - Sector 62"},
	{ID: "hh_2", Name: "Farmhouse - Dadri"},
}

var dashboards = map[string]DashboardData{
	"hh_1": {
		Online: true, WaterSupply: "AVAILABLE", TankLevel: 72, FlowLpm: 18,
		Motor:             Motor{Status: "ON", Mode: "AUTO"},
		Quality:           Quality{TDS: 180, TDSStatus: "NORMAL", Turbidity: 2.1, TurbidityStatus: "NORMAL"},
		ConsumptionTodayL: 340,
		LastAlert:         Alert{Type: "Tank filling started", Severity: "info", Time: "9:42 AM"},
	},
	"hh_2": {
		Online: false, WaterSupply: "UNAVAILABLE", TankLevel: 24, FlowLpm: 0,
		Motor:             Motor{Status: "OFF", Mode: "AUTO"},
		Quality:           Quality{TDS: 410, TDSStatus: "HIGH", Turbidity: 5.6, TurbidityStatus: "HIGH"},
		ConsumptionTodayL: 95,
		LastAlert:         Alert{Type: "No water supply - motor not started", Severity: "warning", Time: "7:15 AM"},
	},
}

func ListHouseholds(c *gin.Context) {
	c.JSON(http.StatusOK, households)
}

func GetHouseholdDashboard(c *gin.Context) {
	id := c.Param("id")
	data, ok := dashboards[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "household not found"})
		return
	}
	c.JSON(http.StatusOK, data)
}
