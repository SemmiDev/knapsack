package main

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

// Data terdiri atas field No, Weight, dan Profit bertipe int
type Data struct {
	No     int
	Weight int
	Profit int
}

// Server terdiri atas fiber instance
type Server struct {
	app *fiber.App
}

// run berfungsi untuk menjalankan server
func (s *Server) run(port string) {
	log.Fatal(s.app.Listen(port))
}

// setupHandler untuk register handler yang dibutuhkan
func (s *Server) setupRoutes() {
	s.app.Get("/", s.homeHandler)
	s.app.Post("/solve", s.solveHandler)
}

// newServer untuk create instance Server
func newServer(viewFolder, format string) *Server {
	// load folder file html didalam folder views untuk templating engine (html)
	config := fiber.Config{
		Views: html.New(viewFolder, format),
	}
	s := Server{fiber.New(config)}
	return &s
}

func main() {
	// create instance server
	app := newServer("./views", ".html")

	// setup routes
	app.setupRoutes()

	// jalankan aplikasi di localhost:9090
	app.run(":9090")
}

// homeHandler nge-render file index.html ke browser
func (s *Server) homeHandler(c *fiber.Ctx) error {
	return c.Render("index", nil)
}

// solveHandler nge-render file solve.html ke browser
func (s *Server) solveHandler(c *fiber.Ctx) error {
	// mengambil kapasitas knapsack dan banyak barang
	// kemudian di convert menjadi integer
	capacity, _ := strconv.Atoi(c.FormValue("capacity"))
	numbers, _ := strconv.Atoi(c.FormValue("numbers"))

	// memanggil method randomStuff untuk generate weights & profits
	// dalam rentang 1 sampai 100 sebanyak total yang inputkan (numbers)
	randWeights, randProfits := randomStuff(1, 100, numbers)

	// hasil random weights & profits di convert kedalam struct Data
	// agar bisa dipanggil di file solve.html
	data := wpToData(randWeights, randProfits)

	// init waktu greedy
	greedyTime := time.Now()

	// menghitung knapsack berdasarkan profit menggunakan greedy
	totalWeightGreedy, totalProfitGreedy, itemsGreedy := greedyKnapsackByProfit(capacity, randWeights, randProfits)

	// hitung waktu yang dihabiskan untuk menghitung knapsack problem menggunakan greedy
	greedyTimeElapsed := time.Since(greedyTime).Milliseconds()

	// init waktu dynamic programming
	dpTime := time.Now()

	// menghitung knapsack berdasarkan profit menggunakan dynamic programming
	totalProfitDP, _, itemsDP := dpKnapsackByProfit(capacity, randWeights, randProfits)

	// hitung waktu yang dihabiskan untuk menghitung knapsack problem menggunakan dynamic programming
	dpTimeElapsed := time.Since(dpTime).Milliseconds()

	// menghitung total weight untuk dynamic programming
	totalWeightDP := totalWeightDP(randWeights, itemsDP)

	// render file solve.html, kemudian sisipkan data untuk dikirim ke view solve.html
	return c.Render("solve", fiber.Map{
		"Data":              data,
		"KapasitasKnapsack": capacity,

		"TotalWeightGreedy": totalWeightGreedy,
		"TotalProfitGreedy": totalProfitGreedy,
		"ItemsGreedy":       itemsGreedy,
		"TimeGreedy":        greedyTimeElapsed,

		"TotalWeightDP": totalWeightDP,
		"TotalProfitDP": totalProfitDP,
		"ItemsDP":       itemsDP,
		"TimeDP":        dpTimeElapsed,
	})
}

// dpKnapsackByProfit solve knapsack problem menggunakan dynamic programming
func dpKnapsackByProfit(capacity int, weights []int, profits []int) (int, int, []int) {
	n := len(weights)
	dp := make([][]int, n+1)
	for i := 0; i < n+1; i++ {
		dp[i] = make([]int, capacity+1)
	}
	for i := 0; i < n+1; i++ {
		for j := 0; j < capacity+1; j++ {
			if i == 0 || j == 0 {
				dp[i][j] = 0
			} else if weights[i-1] <= j {
				dp[i][j] = max(dp[i-1][j], profits[i-1]+dp[i-1][j-weights[i-1]])
			} else {
				dp[i][j] = dp[i-1][j]
			}
		}
	}
	items := make([]int, 0)
	i := n
	j := capacity
	for i > 0 && j > 0 {
		if dp[i][j] != dp[i-1][j] {
			items = append(items, i-1)
			j -= weights[i-1]
		}
		i--
	}
	sort.Ints(items)
	return dp[n][capacity], dp[n][capacity] - dp[n-1][capacity], items
}

// fungsi untuk mencari nilai max dari dua nilai
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// totalWeightDP untuk menghitung total weight dari setiap items DP (index)
func totalWeightDP(weights []int, items []int) (total int) {
	for _, v := range items {
		total += weights[v]
	}
	return
}

// wpToData convert weights & profits ke dalam bentuk array dari struct Data
func wpToData(w []int, p []int) (wp []Data) {
	for i, v := range w {
		wp = append(wp,
			Data{
				No:     i,
				Weight: v,
				Profit: p[i],
			})
	}
	return
}

// greedyKnapsackByProfit solve knapsack problem menggunakan algoritma greedy
func greedyKnapsackByProfit(capacity int, weights []int, profits []int) (totalWeight int, totalProfit int, items []int) {
	for i := 0; i < len(profits); i++ {
		if totalWeight+weights[i] <= capacity {
			totalWeight += weights[i]
			totalProfit += profits[i]
			items = append(items, i)
		}
	}
	return
}

// randomStuff untuk random weights dan profits sebanyak n
func randomStuff(min, max int, n int) (w []int, p []int) {
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= n; i++ {
		w = append(w, min+rand.Intn(max-min))
		p = append(p, min+rand.Intn(max-min))
	}

	// untuk mengurutkan profit secara descending
	sort.Slice(p, func(i, j int) bool {
		return p[i] > p[j]
	})
	return
}
