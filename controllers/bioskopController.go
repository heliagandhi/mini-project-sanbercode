package controllers

import (
	"database/sql"
	"mini-project-sanbercode/database"
	"mini-project-sanbercode/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// validasi field
func formatValidationError(err error) []string {
	var errors []string

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			switch e.Field() {
			case "Nama":
				errors = append(errors, "nama wajib diisi")
			case "Lokasi":
				errors = append(errors, "lokasi wajib diisi")
			}
		}
	} else {
		errors = append(errors, "format JSON tidak valid")
	}

	return errors
}

// validasi rating
func isValidRating(r float64) bool {
        return r >= 1 && r <= 5
}

// create bioskop
func CreateBioskop(ctx *gin.Context) {
	var bioskop models.Bioskop

	if err := ctx.ShouldBindJSON(&bioskop); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors": formatValidationError(err),
		})
		return
	}

	// validasi rating
	if !isValidRating(*bioskop.Rating) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "rating harus antara 1 sampai 5",
		})
		return
	}

	// cek nama uniq
	var exists bool
	err := database.DB.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM bioskops WHERE nama=$1)",
		bioskop.Nama,
	).Scan(&exists)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal cek data",
		})
		return
	}

	if exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "nama bioskop sudah digunakan",
		})
		return
	}

	query := `
		INSERT INTO bioskops (nama, lokasi, rating)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err = database.DB.QueryRow(
		query,
		bioskop.Nama,
		bioskop.Lokasi,
		bioskop.Rating,
	).Scan(&bioskop.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal insert data",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "bioskop berhasil ditambahkan",
		"data":    bioskop,
	})
}

// get all bioskop
func GetBioskops(ctx *gin.Context) {
	rows, err := database.DB.Query(
		"SELECT id, nama, lokasi, rating FROM bioskops ORDER BY id ASC",
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal mengambil data",
		})
		return
	}
	defer rows.Close()

	var bioskops []models.Bioskop

	for rows.Next() {
		var b models.Bioskop
		if err := rows.Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "gagal membaca data",
			})
			return
		}
		bioskops = append(bioskops, b)
	}

	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error saat iterasi data",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": bioskops,
	})
}

// get by id bioskop
func GetBioskopByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id harus berupa angka",
		})
		return
	}

	var b models.Bioskop

	err = database.DB.QueryRow(
		"SELECT id, nama, lokasi, rating FROM bioskops WHERE id=$1",
		id,
	).Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "data tidak ditemukan",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "gagal mengambil data",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": b,
	})
}

// update bioskop by id
func UpdateBioskop(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id harus berupa angka",
		})
		return
	}

	var bioskop models.Bioskop

	if err := ctx.ShouldBindJSON(&bioskop); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors": formatValidationError(err),
		})
		return
	}

	// validasi rating
	if !isValidRating(*bioskop.Rating) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "rating harus antara 1 sampai 5",
		})
		return
	}

	// cek nama unique (exclude dirinya sendiri)
	var exists bool
	err = database.DB.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM bioskops WHERE nama=$1 AND id <> $2)",
		bioskop.Nama,
		id,
	).Scan(&exists)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal cek data",
		})
		return
	}

	if exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "nama bioskop sudah digunakan",
		})
		return
	}

	result, err := database.DB.Exec(
		`UPDATE bioskops SET nama=$1, lokasi=$2, rating=$3 WHERE id=$4`,
		bioskop.Nama,
		bioskop.Lokasi,
		bioskop.Rating,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal update data",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "data tidak ditemukan",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "data berhasil diupdate",
	})
}

// delete bioskop by id
func DeleteBioskop(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id harus berupa angka",
		})
		return
	}

	result, err := database.DB.Exec(
		"DELETE FROM bioskops WHERE id=$1",
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal menghapus data",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "data tidak ditemukan",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "data berhasil dihapus",
	})
}