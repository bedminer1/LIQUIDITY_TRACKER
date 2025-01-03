package riskassessment

import (
	"testing"

	"github.com/bedminer1/liquidity_tracker/internal/models"
)

func TestAssessLiquidity(t *testing.T) {
	currentRecords := []models.Record{
		{AssetType: "ETF", Timestamp: "2013-03-29T00:00:00Z", BidAskSpread: 0.0207, Volume: 620875, BidPrice: 117.46},
		{AssetType: "ETF", Timestamp: "2013-04-01T00:00:00Z", BidAskSpread: 0.022, Volume: 945668, BidPrice: 117.05},
		{AssetType: "ETF", Timestamp: "2013-04-02T00:00:00Z", BidAskSpread: 0.0401, Volume: 765685, BidPrice: 117.27},
		{AssetType: "ETF", Timestamp: "2013-04-03T00:00:00Z", BidAskSpread: 0.0279, Volume: 1016691, BidPrice: 117.57},
		{AssetType: "ETF", Timestamp: "2013-04-04T00:00:00Z", BidAskSpread: 0.027, Volume: 862085, BidPrice: 118.15},
		{AssetType: "ETF", Timestamp: "2013-04-05T00:00:00Z", BidAskSpread: 0.024, Volume: 1059950, BidPrice: 119.21},
		{AssetType: "ETF", Timestamp: "2013-04-08T00:00:00Z", BidAskSpread: 0.0223, Volume: 579526, BidPrice: 120.29},
		{AssetType: "ETF", Timestamp: "2013-04-09T00:00:00Z", BidAskSpread: 0.0177, Volume: 647444, BidPrice: 120.31},
		{AssetType: "ETF", Timestamp: "2013-04-10T00:00:00Z", BidAskSpread: 0.0289, Volume: 647328, BidPrice: 120.23},
		{AssetType: "ETF", Timestamp: "2013-04-11T00:00:00Z", BidAskSpread: 0.0198, Volume: 653664, BidPrice: 120.46},
	}

	predictions := []models.Record{
		{AssetType: "ETF", Timestamp: "2013-04-12T00:00:00Z", BidAskSpread: 0.0176, Volume: 806282, BidPrice: 120.6972},
		{AssetType: "ETF", Timestamp: "2013-04-15T00:00:00Z", BidAskSpread: 0.0504, Volume: 462017, BidPrice: 120.51},
		{AssetType: "ETF", Timestamp: "2013-04-16T00:00:00Z", BidAskSpread: 0.0506, Volume: 1375894, BidPrice: 120.35},
		{AssetType: "ETF", Timestamp: "2013-04-17T00:00:00Z", BidAskSpread: 0.0227, Volume: 1258238, BidPrice: 119.85},
		{AssetType: "ETF", Timestamp: "2013-04-18T00:00:00Z", BidAskSpread: 0.0356, Volume: 771955, BidPrice: 119.84},
		{AssetType: "ETF", Timestamp: "2013-04-19T00:00:00Z", BidAskSpread: 0.0294, Volume: 353449, BidPrice: 119.85},
		{AssetType: "ETF", Timestamp: "2013-04-22T00:00:00Z", BidAskSpread: 0.0898, Volume: 472687, BidPrice: 120.19},
		{AssetType: "ETF", Timestamp: "2013-04-23T00:00:00Z", BidAskSpread: 0.0234, Volume: 700594, BidPrice: 120.48},
		{AssetType: "ETF", Timestamp: "2013-04-24T00:00:00Z", BidAskSpread: 0.0201, Volume: 696207, BidPrice: 120.65},
	}

	windowSize := 20

	report := AssessLiquidity(currentRecords, predictions, windowSize)
	// fmt.Println(report.CurrentHighRiskCount, report.PredictedHighRiskCount, report.CurrentModerateRiskCount, report.PredictedModerateRiskCount)

	// Expected values
	expectedCurrentHighRiskCount := 0
	expectedPredictedHighRiskCount := 0
	expectedCurrentModerateRiskCount := 1
	expectedPredictedModerateRiskCount := 4
	expectedWarningMessageOne := "Predicted high risk for ETF at 2013-04-15T00:00:00Z: Spread=0.04% (MA=0.02%), Volume=462017 (MA=643347)"

	// Assertions
	if report.CurrentHighRiskCount != expectedCurrentHighRiskCount {
		t.Errorf("CurrentHighRiskCount: got %d, want %d", report.CurrentHighRiskCount, expectedCurrentHighRiskCount)
	}
	if report.PredictedHighRiskCount != expectedPredictedHighRiskCount {
		t.Errorf("PredictedHighRiskCount: got %d, want %d", report.PredictedHighRiskCount, expectedPredictedHighRiskCount)
	}
	if report.CurrentModerateRiskCount != expectedCurrentModerateRiskCount {
		t.Errorf("CurrentModerateRiskCount: got %d, want %d", report.CurrentModerateRiskCount, expectedCurrentModerateRiskCount)
	}
	if report.PredictedModerateRiskCount != expectedPredictedModerateRiskCount {
		t.Errorf("PredictedModerateRiskCount: got %d, want %d", report.PredictedModerateRiskCount, expectedPredictedModerateRiskCount)
	}
	if len(report.PredictedWarnings) > 0 && report.PredictedWarnings[0] != expectedWarningMessageOne {
		t.Errorf("PredictedWarnings: got %s, want %s", report.PredictedWarnings[0], expectedWarningMessageOne)
	}
}
