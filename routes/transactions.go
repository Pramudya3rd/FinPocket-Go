package routes

import (
	"strconv"
	"time"

	"finpocket.com/api/database"
	"finpocket.com/api/models"
	"github.com/gofiber/fiber/v2"
)

func GetTransactions(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	typeParam := c.Query("type")
	categoryParam := c.Query("category_id")
	month := c.Query("month")
	limit := c.Query("limit")

	// Konversi parameter type
	var transactionType *int
	if typeParam != "" {
		val, err := strconv.Atoi(typeParam)
		if err != nil || (val != 0 && val != 1) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "type harus bernilai 0 (spending) atau 1 (income)",
				"status":  "error",
				"data":    nil,
			})
		}
		transactionType = &val
	}

	// Konversi parameter category_id
	var categoryID *int
	if categoryParam != "" {
		val, err := strconv.Atoi(categoryParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "category_id harus berupa angka",
				"status":  "error",
				"data":    nil,
			})
		}
		categoryID = &val
	}

	// Konversi parameter limit
	var limitValue int
	if limit != "" {
		val, err := strconv.Atoi(limit)
		if err != nil || val < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "limit harus berupa angka positif",
				"status":  "error",
				"data":    nil,
			})
		}
		limitValue = val
	}

	// Mengambil transaksi berdasarkan filter
	var transactions []models.Transaction
	query := database.DBConn.Model(&models.Transaction{})

	// Filter berdasarkan type
	if transactionType != nil {
		if *transactionType == 1 {
			query = query.Where("amount > 0")
		} else {
			query = query.Where("amount < 0")
		}
	}

	// Filter berdasarkan category_id
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	// Filter berdasarkan bulan
	if month != "" {
		monthInt, err := strconv.Atoi(month)
		if err != nil || monthInt < 1 || monthInt > 12 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "bulan harus berupa angka antara 1 dan 12",
				"status":  "error",
				"data":    nil,
			})
		}

		now := time.Now()
		startDate := time.Date(now.Year(), time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0)

		query = query.Where("created_at >= ? AND created_at < ?", startDate, endDate)
	}

	// Mengurutkan berdasarkan tanggal terbaru
	query = query.Where("user_id = ?", user.ID).Order("created_at DESC")

	// Terapkan limit jika ada
	if limitValue > 0 {
		query = query.Limit(limitValue)
	}

	// Eksekusi query
	if err := query.Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data transaksi",
			"status":  "error",
			"data":    nil,
		})
	}

	if len(transactions) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Tidak ada transaksi yang ditemukan",
			"status":  "error",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil data transaksi",
		"status":  "success",
		"data":    transactions,
	})
}

func CreateTransaction(c *fiber.Ctx) error {
	// Struct untuk menerima input dari request body
	user := c.Locals("user").(*models.User)
	type Request struct {
		CategoryID  int64   `json:"category_id"`
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
	}

	var request Request

	// Parsing body request
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
			"data":    nil,
		})
	}

	// Validasi input
	if request.CategoryID == 0 || request.Amount == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "category_id dan amount tidak boleh kosong",
			"status":  "error",
			"data":    nil,
		})
	}

	// Ambil transaksi terakhir berdasarkan jenis (income atau spending)
	var lastTransaction models.Transaction
	transactionType := "income"
	if request.Amount < 0 {
		transactionType = "spending"
	}

	if transactionType == "income" {
		if err := database.DBConn.Where("amount > 0").
			Order("created_at DESC").
			First(&lastTransaction).Error; err != nil && err.Error() != "record not found" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal mengambil transaksi terakhir",
				"status":  "error",
				"data":    nil,
			})
		}
	} else {
		if err := database.DBConn.Where("amount < 0").
			Order("created_at DESC").
			First(&lastTransaction).Error; err != nil && err.Error() != "record not found" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal mengambil transaksi terakhir",
				"status":  "error",
				"data":    nil,
			})
		}
	}

	// Hitung total baru berdasarkan transaksi terakhir
	newTotal := lastTransaction.Total + request.Amount

	// Buat data transaksi baru
	transaction := models.Transaction{
		UserID:      user.ID,
		CategoryID:  request.CategoryID,
		Description: request.Description,
		Amount:      request.Amount,
		Total:       newTotal,
	}

	// Simpan transaksi ke dalam database
	if err := database.DBConn.Create(&transaction).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menambahkan transaksi",
			"status":  "error",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Transaksi berhasil ditambahkan",
		"status":  "success",
		"data": fiber.Map{
			"transaction": transaction,
			"total":       newTotal,
		},
	})
}

func GetTransactionSummaries(c *fiber.Ctx) error {
	// Ambil parameter bulan
	user := c.Locals("user").(*models.User)
	month := c.Query("month")

	// Validasi parameter bulan
	if month == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Parameter 'month' diperlukan",
			"status":  "error",
			"data":    nil,
		})
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Parameter 'month' harus berupa angka antara 1 dan 12",
			"status":  "error",
			"data":    nil,
		})
	}

	// Hitung rentang waktu bulan
	now := time.Now()
	startDate := time.Date(now.Year(), time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	// Variabel untuk menyimpan total income dan spending
	var totalIncome, totalSpending float64

	// Query untuk total income
	err = database.DBConn.Model(&models.Transaction{}).
		Select("SUM(amount)").
		Where("user_id = ?", user.ID).
		Where("amount > 0 AND created_at >= ? AND created_at < ?", startDate, endDate).
		Scan(&totalIncome).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghitung total income",
			"status":  "error",
			"data":    nil,
		})
	}

	// Query untuk total spending
	err = database.DBConn.Model(&models.Transaction{}).
		Select("SUM(amount)").
		Where("user_id = ?", user.ID).
		Where("amount < 0 AND created_at >= ? AND created_at < ?", startDate, endDate).
		Scan(&totalSpending).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghitung total spending",
			"status":  "error",
			"data":    nil,
		})
	}

	// Ubah totalSpending menjadi positif (karena disimpan sebagai negatif)
	totalSpending = -totalSpending

	// Kembalikan hasil
	return c.JSON(fiber.Map{
		"message": "Berhasil menghitung ringkasan transaksi",
		"status":  "success",
		"data": fiber.Map{
			"spending": totalSpending,
			"income":   totalIncome,
		},
	})
}
