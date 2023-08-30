package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type newStudent struct {
	Student_id       uint32 `json:"student_id" binding:"required"`
	Student_name     string `json:"student_name" binding:"required"`
	Student_age      uint8  `json:"student_age" binding:"required"`
	Student_address  string `json:"student_address" binding:"required"`
	Student_phone_no string `json:"student_phone_no" binding:"required"`
}

func rowToStruct(rows *sql.Rows, dest interface{}) error {
	destv := reflect.ValueOf(dest).Elem()

	args := make([]interface{}, destv.Type().Elem().NumField())

	for rows.Next() {
		rowp := reflect.New(destv.Type().Elem())
		rowv := rowp.Elem()

		for i := 0; i < rowv.NumField(); i++ {
			args[i] = rowv.Field(i).Addr().Interface()
		}

		if err := rows.Scan(args...); err != nil {
			return err
		}
		destv.Set(reflect.Append(destv, rowv))
	}
	return nil
}

func postHandler(c *gin.Context, db *gorm.DB) {
	var newStudent newStudent

	// if c.Bind(&newStudent) == nil {
	// 	_, err := db.Exec("INSERT INTO students values (?,?,?,?,?)", newStudent.Student_id, newStudent.Student_name, newStudent.Student_age, newStudent.Student_address, newStudent.Student_phone_no)
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"massage": err.Error()})
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{"massage": "Success Create"})
	// }
	// c.JSON(http.StatusBadRequest, gin.H{"status": "Errorr"})

	c.Bind(&newStudent)
	db.Create(&newStudent)
	c.JSON(http.StatusOK, gin.H{
		"massage": "success create",
		"data":    newStudent,
	})
}

func getAllHandler(c *gin.Context, db *gorm.DB) {
	var newStudent []newStudent

	// 	row, err := db.Query("SELECT * FROM students")
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	}

	// 	rowToStruct(row, &newStudent)

	// 	if newStudent == nil {
	// 		c.JSON(http.StatusNotFound, gin.H{"Massage": "data not found"})
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{"data": newStudent})
	// }

	db.Find(&newStudent)
	c.JSON(http.StatusOK, gin.H{
		"massage": "success find all data",
		"data":    newStudent,
	})
}

func getHandler(c *gin.Context, db *gorm.DB) {
	// var newStudent []newStudent

	// studentId := c.Param("student_id")

	// row, err := db.Query("SELECT * FROM students WHERE student_id = ?", studentId)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
	// }

	// rowToStruct(row, &newStudent)

	// if newStudent == nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"Massage": "data not found"})
	// }
	// c.JSON(http.StatusOK, gin.H{"data": newStudent})

	var newStudent newStudent

	studentId := c.Param("student_id")

	if db.Find(&newStudent, "student_id=?", studentId).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{
			"massage": "data not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"massage": "success find by Id",
		"data":    newStudent,
	})

}

func putHandler(c *gin.Context, db *gorm.DB) {
	// var newStudent newStudent

	// studentId := c.Param("student_id")

	// if c.Bind(&newStudent) == nil {
	// 	_, err := db.Exec("UPDATE students SET student_name=?, student_age=?, student_address=?, student_phone_no=? WHERE student_id=?", newStudent.Student_name, newStudent.Student_age, newStudent.Student_address, newStudent.Student_phone_no, studentId)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{"massage": "Update Success"})
	// }

	var newStudent newStudent

	studentId := c.Param("student_id")

	if db.Find(&newStudent, "student_id=?", studentId).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{
			"massage": "not found",
		})
		return
	}

	var reqStudent = newStudent

	c.Bind(&reqStudent)

	db.Model(&newStudent).Where("student_id=?", studentId).Update(reqStudent)

	c.JSON(http.StatusOK, gin.H{
		"massaege": "update success",
		"data":     reqStudent,
	})
}

func delHandler(c *gin.Context, db *gorm.DB) {
	// studentId := c.Param("student_id")

	// _, err := db.Exec("DELETE FROM students WHERE student_id=?", studentId)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	// c.JSON(http.StatusOK, gin.H{"massage": "Delete Success"})

	var newStudent newStudent

	studentId := c.Param("student_id")

	db.Delete(&newStudent, "student_id=?", studentId)

	c.JSON(http.StatusOK, gin.H{
		"massage": "Delete Success",
	})
}

// HANDLER
func setupRouter() *gin.Engine {
	db, err := gorm.Open("mysql", "root:@/sampleapi")
	if err != nil {
		panic(err)
	}

	Migrate(db)

	r := gin.Default()

	r.POST("/student", func(ctx *gin.Context) {
		postHandler(ctx, db)
	})
	r.GET("/student", func(ctx *gin.Context) {
		getAllHandler(ctx, db)
	})
	r.GET("/student/:student_id", func(ctx *gin.Context) {
		getHandler(ctx, db)
	})
	r.PUT("/student/:student_id", func(ctx *gin.Context) {
		putHandler(ctx, db)
	})
	r.DELETE("/student/:student_id", func(ctx *gin.Context) {
		delHandler(ctx, db)
	})

	return r
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&newStudent{})

	data := newStudent{}
	if db.Find(&data).RecordNotFound() {
		fmt.Println("======= run seeder user ======")
		seederUser(db)
	}
}

func seederUser(db *gorm.DB) {
	data := newStudent{
		Student_id:       4,
		Student_name:     "jojo",
		Student_age:      45,
		Student_address:  "semarang",
		Student_phone_no: "08986765645",
	}
	db.Create(&data)
}

func main() {

	r := setupRouter()

	r.Run(":8080")

}
