package v1

import (
	"github.com/kataras/iris/v12/core/router"
	"webapi/controller"
	"webapi/internal/wrapper"
)

func RegisterCourseRouter(party router.Party) {
	party.Handle("POST", "/", wrapper.Handler(controller.CourseController{}.CreateCourse)) //添加课程
	party.Handle("DELETE", "/", wrapper.Handler(controller.CourseController{}.DeleteCourse)) //删除课程
	party.Handle("POST", "/update", wrapper.Handler(controller.CourseController{}.UpdateCourse)) //编辑课程信息
	party.Handle("GET","/",wrapper.Handler(controller.CourseController{}.CourseList)) //获取课程列表
}
