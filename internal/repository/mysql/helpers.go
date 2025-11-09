package mysql

import (
	"strings"

	"github.com/go-sql-driver/mysql"
)

// Mengecek apakah error adalah mysql key error
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	// Cek nysql error number
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return mysqlErr.Number == 1062
	}

	// Fallback jika bukan *mysql.MySQLError
	return strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "duplicate key")
}

// // Mengecek error apakah foreign key constraint violation
// func isForeignKeyError(err error) bool {
// 	if err == nil {
// 		return false
// 	}

// 	// Cek nysql error number
// 	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
// 		return mysqlErr.Number == 1452 || mysqlErr.Number == 1451
// 	}

// 	// Fallback jika bukan *mysql.MySQLError
// 	return strings.Contains(err.Error(), "foreign key constraint")
// }

// // Mengecek apakah error deadlock
// func isDeadlockError(err error) bool {
// 	if err == nil {
// 		return false
// 	}

// 	// Cek nysql error number
// 	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
// 		return mysqlErr.Number == 1213
// 	}

// 	// Fallback jika bukan *mysql.MySQLError
// 	return strings.Contains(err.Error(), "Deadlock")
// }
