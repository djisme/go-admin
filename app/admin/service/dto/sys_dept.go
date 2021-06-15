package dto

import (
	"go-admin/app/admin/models"
	common "go-admin/common/models"

	"go-admin/common/dto"
)

// SysDeptSearch 列表或者搜索使用结构体
type SysDeptSearch struct {
	dto.Pagination `search:"-"`
	DeptId         int    `form:"deptId" search:"type:exact;column:dept_id;table:sys_dept" comment:"id"`       //id
	ParentId       int    `form:"parentId" search:"type:exact;column:parent_id;table:sys_dept" comment:"上级部门"` //上级部门
	DeptPath       string `form:"deptPath" search:"type:exact;column:dept_path;table:sys_dept" comment:""`     //路径
	DeptName       string `form:"deptName" search:"type:exact;column:dept_name;table:sys_dept" comment:"部门名称"` //部门名称
	Sort           int    `form:"sort" search:"type:exact;column:sort;table:sys_dept" comment:"排序"`            //排序
	Leader         string `form:"leader" search:"type:exact;column:leader;table:sys_dept" comment:"负责人"`       //负责人
	Phone          string `form:"phone" search:"type:exact;column:phone;table:sys_dept" comment:"手机"`          //手机
	Email          string `form:"email" search:"type:exact;column:email;table:sys_dept" comment:"邮箱"`          //邮箱
	Status         string `form:"status" search:"type:exact;column:status;table:sys_dept" comment:"状态"`        //状态
}

func (m *SysDeptSearch) GetNeedSearch() interface{} {
	return *m
}


// SysDeptControl 增、改使用的结构体
type SysDeptControl struct {
	DeptId   int    `uri:"id" comment:"编码"`          // 编码
	ParentId int    `form:"parentId" comment:"上级部门"` //上级部门
	DeptPath string `form:"deptPath" comment:""`     //路径
	DeptName string `form:"deptName" comment:"部门名称"` //部门名称
	Sort     int    `form:"sort" comment:"排序"`       //排序
	Leader   string `form:"leader" comment:"负责人"`    //负责人
	Phone    string `form:"phone" comment:"手机"`      //手机
	Email    string `form:"email" comment:"邮箱"`      //邮箱
	Status   string `form:"status" comment:"状态"`     //状态
	common.ControlBy
}

// Generate 结构体数据转化 从 SysDeptControl 至 SysDept 对应的模型
func (s *SysDeptControl) Generate(model *models.SysDept) {
	if s.DeptId != 0 {
		model.DeptId = s.DeptId
	}
	model.DeptName = s.DeptName
	model.ParentId = s.ParentId
	model.DeptPath = s.DeptPath
	model.Sort = s.Sort
	model.Leader = s.Leader
	model.Phone = s.Phone
	model.Email = s.Email
	model.Status = s.Status
}

// GetId 获取数据对应的ID
func (s *SysDeptControl) GetId() interface{} {
	return s.DeptId
}

// SysDeptById 获取单个或者删除的结构体
type SysDeptById struct {
	Id  int   `uri:"id"`
	Ids []int `json:"ids"`
}

func (s *SysDeptById) Generate() *SysDeptById {
	cp := *s
	return &cp
}

func (s *SysDeptById) GetId() interface{} {
	if len(s.Ids) > 0 {
		s.Ids = append(s.Ids, s.Id)
		return s.Ids
	}
	return s.Id
}

func (s *SysDeptById) GenerateM() (*models.SysDept, error) {
	return &models.SysDept{}, nil
}

type DeptLabel struct {
	Id       int         `gorm:"-" json:"id"`
	Label    string      `gorm:"-" json:"label"`
	Children []DeptLabel `gorm:"-" json:"children"`
}