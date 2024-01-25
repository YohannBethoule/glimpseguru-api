package router

type AnalyticsRequest struct {
	StartTime int64 `form:"start_time" binding:"required"`
	EndTime   int64 `form:"end_time" binding:"required"`
}
