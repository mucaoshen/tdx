package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/injoyai/tdx"
	"github.com/injoyai/tdx/extend"
	"github.com/injoyai/tdx/protocol"
)

// ── 代码列表相关 ──────────────────────────────────────────────────────────

func handleGetCodes(w http.ResponseWriter, r *http.Request) {
	if tdx.DefaultCodes == nil {
		errorResponse(w, "代码缓存未初始化")
		return
	}
	limit := parsePositiveInt(r.URL.Query().Get("limit"))
	all := tdx.DefaultCodes.GetStocks(limit)
	successResponse(w, map[string]interface{}{
		"count": len(all),
		"list":  all,
	})
}

func handleBatchQuote(w http.ResponseWriter, r *http.Request) {
	codes := splitCodes(r.URL.Query().Get("codes"))
	if len(codes) == 0 {
		errorResponse(w, "codes不能为空")
		return
	}
	quotes, err := client.GetQuote(codes...)
	if err != nil {
		errorResponse(w, err.Error())
		return
	}
	successResponse(w, quotes)
}

func handleGetKlineHistory(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		errorResponse(w, "code不能为空")
		return
	}
	resp, err := getQfqKlineDay(code)
	if err != nil {
		errorResponse(w, err.Error())
		return
	}
	successResponse(w, resp)
}

func handleGetIndex(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		errorResponse(w, "code不能为空")
		return
	}
	resp, err := client.GetIndex(0, code, 0, 800)
	if err != nil {
		errorResponse(w, err.Error())
		return
	}
	successResponse(w, resp)
}

func handleGetIndexAll(w http.ResponseWriter, r *http.Request) {
	quotes, err := client.GetQuote()
	if err != nil {
		errorResponse(w, err.Error())
		return
	}
	successResponse(w, quotes)
}

func handleGetMarketStats(w http.ResponseWriter, r *http.Request) {
	resp, err := client.GetQuote()
	if err != nil {
		errorResponse(w, err.Error())
		return
	}
	up, down, flat := 0, 0, 0
	for _, q := range resp {
		if q.Kline == nil {
			flat++
			continue
		}
		switch {
		case q.Kline.Last == 0:
			flat++
		case q.Kline.Close > q.Kline.Last:
			up++
		default:
			down++
		}
	}
	successResponse(w, map[string]interface{}{
		"total": len(resp), "up": up, "down": down, "flat": flat,
	})
}

func handleGetMarketCount(w http.ResponseWriter, r *http.Request) {
	if tdx.DefaultCodes == nil {
		errorResponse(w, "代码缓存未初始化")
		return
	}
	n := len(tdx.DefaultCodes.GetStocks())
	successResponse(w, map[string]int{"count": n})
}

func handleGetStockCodes(w http.ResponseWriter, r *http.Request) {
	if tdx.DefaultCodes == nil {
		errorResponse(w, "代码缓存未初始化")
		return
	}
	codes := tdx.DefaultCodes.GetStockCodes()
	successResponse(w, codes)
}

func handleGetETFCodes(w http.ResponseWriter, r *http.Request) {
	if tdx.DefaultCodes == nil {
		errorResponse(w, "代码缓存未初始化")
		return
	}
	codes := tdx.DefaultCodes.GetETFCodes()
	successResponse(w, codes)
}

func handleGetETFList(w http.ResponseWriter, r *http.Request) {
	if tdx.DefaultCodes == nil {
		errorResponse(w, "代码缓存未初始化")
		return
	}
	etfs := tdx.DefaultCodes.GetETFs()
	successResponse(w, etfs)
}

func handleGetTradeHistory(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		errorResponse(w, "code不能为空")
		return
	}
	resp, err := client.GetMinuteTrade(code, 0, 1800)
	if err != nil {
		errorResponse(w, err.Error())
		return
	}
	successResponse(w, resp)
}

func handleGetMinuteTradeAll(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("20060102")
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		errorResponse(w, "code不能为空")
		return
	}
	resp, err := client.GetHistoryMinuteTradeDay(date, code)
	if err != nil {
		errorResponse(w, err.Error())
		return
	}
	successResponse(w, resp)
}

func handleGetTradeHistoryFull(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		errorResponse(w, "code不能为空")
		return
	}
	resp, err := client.GetMinuteTrade(code, 0, 1800)
	if err != nil {
		errorResponse(w, err.Error())
		return
	}
	successResponse(w, resp)
}

func handleGetKlineAllTDX(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		errorResponse(w, "code不能为空")
		return
	}
	resp, err := client.GetKlineDayAll(code)
	if err != nil {
		errorResponse(w, err.Error())
		return
	}
	successResponse(w, resp)
}

func handleGetKlineAllTHS(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		errorResponse(w, "code不能为空")
		return
	}
	klines, err := extend.GetTHSDayKline(code, extend.THS_QFQ)
	if err != nil {
		errorResponse(w, err.Error())
		return
	}
	successResponse(w, klines)
}

func handleGetWorkday(w http.ResponseWriter, r *http.Request) {
	if manager == nil || manager.Workday == nil {
		errorResponse(w, "交易日数据未初始化")
		return
	}
	today := time.Now().Format("2006-01-02")
	isWorkday := manager.Workday.TodayIs()
	successResponse(w, map[string]interface{}{
		"date":    today,
		"workday": isWorkday,
	})
}

func handleGetWorkdayRange(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	successResponse(w, map[string]string{"start": start, "end": end})
}

func handleGetIncome(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		errorResponse(w, "code不能为空")
		return
	}
	klines, err := client.GetKlineDayAll(code)
	if err != nil {
		errorResponse(w, err.Error())
		return
	}
	// ⚠️ extend.DoIncomes API 已变更 — 待适配
	_ = extend.DoIncomes
	_ = klines
	successResponse(w, map[string]string{
		"status": "income calculation not available in current API version",
	})
}

func handleGetServerStatus(w http.ResponseWriter, r *http.Request) {
	if tdx.DefaultCodes == nil {
		successResponse(w, map[string]string{"status": "initializing"})
		return
	}
	n := len(tdx.DefaultCodes.GetStocks())
	successResponse(w, map[string]interface{}{
		"status":      "running",
		"stock_count": n,
		"uptime":      time.Now().Unix(),
	})
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(`{"status":"ok"}`))
}

func getAllCodeModels() ([]*protocol.Code, error) {
	var models []*protocol.Code
	if tdx.DefaultCodes != nil {
		for code, model := range tdx.DefaultCodes.Iter() {
			_ = code
			models = append(models, &protocol.Code{
				Code: model.Code,
				Name: model.Name,
			})
		}
		return models, nil
	}
	return nil, fmt.Errorf("代码库未初始化")
}

func parsePositiveInt(s string) int {
	n := 0
	for _, c := range strings.TrimSpace(s) {
		if c < '0' || c > '9' {
			return n
		}
		n = n*10 + int(c-'0')
	}
	return n
}
