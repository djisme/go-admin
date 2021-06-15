package service

import (
	"errors"
	"go-admin/app/admin/models"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"

	"gorm.io/gorm"

	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"

	"github.com/go-admin-team/go-admin-core/sdk/service"
)

type SysDept struct {
	service.Service
}

// GetPage 获取SysDept列表
func (e *SysDept) GetPage(c *dto.SysDeptSearch, list *[]models.SysDept) error {
	var err error
	var data models.SysDept

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).
		Find(list).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Get 获取SysDept对象
func (e *SysDept) Get(d *dto.SysDeptById, model *models.SysDept) error {
	var err error
	var data models.SysDept

	db := e.Orm.Model(&data).
		First(model, d.GetId())
	err = db.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if db.Error != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建SysDept对象
func (e *SysDept) Insert(c *dto.SysDeptControl) error {
	var err error
	var data models.SysDept
	c.Generate(&data)
	err = e.Orm.Create(data).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	deptPath := pkg.IntToString(data.DeptId) + "/"
	if data.ParentId != 0 {
		var deptP models.SysDept
		e.Orm.First(&deptP, data.ParentId)
		deptPath = deptP.DeptPath + deptPath
	} else {
		deptPath = "/0/" + deptPath
	}
	var mp = map[string]string{}
	mp["dept_path"] = deptPath
	if err := e.Orm.Update("dept_path", deptPath).Error; err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Update 修改SysDept对象
func (e *SysDept) Update(c *dto.SysDeptControl) error {
	var err error
	var model = models.SysDept{}
	e.Orm.First(&model, c.GetId())
	c.Generate(&model)

	deptPath := pkg.IntToString(model.DeptId) + "/"
	if model.ParentId != 0 {
		var deptP models.SysDept
		e.Orm.First(&deptP, model.ParentId)
		deptPath = deptP.DeptPath + deptPath
	} else {
		deptPath = "/0/" + deptPath
	}
	model.DeptPath = deptPath
	db := e.Orm.Save(&model)
	if db.Error != nil {
		e.Log.Errorf("UpdateSysDept error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除SysDept
func (e *SysDept) Remove(d *dto.SysDeptById) error {
	var err error
	var data models.SysDept

	db := e.Orm.Model(&data).Delete(&data, d.GetId())
	if db.Error != nil {
		err = db.Error
		e.Log.Errorf("Delete error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		err = errors.New("无权删除该数据")
		return err
	}
	return nil
}

// GetSysDeptList 获取组织数据
func (e *SysDept) getList(c *dto.SysDeptSearch, list *[]models.SysDept) error {
	var err error
	var data models.SysDept

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).
		Find(list).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// SetDeptTree 设置组织数据
func (e *SysDept) SetDeptTree(c *dto.SysDeptSearch) (m []dto.DeptLabel, err error) {
	var list []models.SysDept
	err = e.getList(c, &list)

	m = make([]dto.DeptLabel, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		e := dto.DeptLabel{}
		e.Id = list[i].DeptId
		e.Label = list[i].DeptName
		deptsInfo := deptTreeCall(&list, e)

		m = append(m, deptsInfo)
	}
	return
}

// Call 递归构造组织数据
func deptTreeCall(deptList *[]models.SysDept, dept dto.DeptLabel) dto.DeptLabel {
	list := *deptList
	min := make([]dto.DeptLabel, 0)
	for j := 0; j < len(list); j++ {
		if dept.Id != list[j].ParentId {
			continue
		}
		mi := dto.DeptLabel{Id: list[j].DeptId, Label: list[j].DeptName, Children: []dto.DeptLabel{}}
		ms := deptTreeCall(deptList, mi)
		min = append(min, ms)
	}
	dept.Children = min
	return dept
}

// SetDeptPage 设置dept页面数据
func (e *SysDept) SetDeptPage(c *dto.SysDeptSearch) (m []models.SysDept, err error) {
	var list []models.SysDept
	err = e.getList(c, &list)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		info := e.deptPageCall(&list, list[i])
		m = append(m, info)
	}
	return
}

func (e *SysDept) deptPageCall(deptlist *[]models.SysDept, menu models.SysDept) models.SysDept {
	list := *deptlist
	min := make([]models.SysDept, 0)
	for j := 0; j < len(list); j++ {
		if menu.DeptId != list[j].ParentId {
			continue
		}
		mi := models.SysDept{}
		mi.DeptId = list[j].DeptId
		mi.ParentId = list[j].ParentId
		mi.DeptPath = list[j].DeptPath
		mi.DeptName = list[j].DeptName
		mi.Sort = list[j].Sort
		mi.Leader = list[j].Leader
		mi.Phone = list[j].Phone
		mi.Email = list[j].Email
		mi.Status = list[j].Status
		mi.CreatedAt = list[j].CreatedAt
		mi.Children = []models.SysDept{}
		ms := e.deptPageCall(deptlist, mi)
		min = append(min, ms)
	}
	menu.Children = min
	return menu
}

// GetRoleDeptId 获取角色的部门ID集合
func (e *SysDept) GetWithRoleId(roleId int) ([]int, error) {
	deptIds := make([]int, 0)
	deptList := make([]dto.DeptIdList, 0)
	if err := e.Orm.Table("sys_role_dept").
		Select("sys_role_dept.dept_id").
		Joins("LEFT JOIN sys_dept on sys_dept.dept_id=sys_role_dept.dept_id").
		Where("role_id = ? ", roleId).
		Where(" sys_role_dept.dept_id not in(select sys_dept.parent_id from sys_role_dept LEFT JOIN sys_dept on sys_dept.dept_id=sys_role_dept.dept_id where role_id =? )", roleId).
		Find(&deptList).Error; err != nil {
		return nil, err
	}
	for i := 0; i < len(deptList); i++ {
		deptIds = append(deptIds, deptList[i].DeptId)
	}
	return deptIds, nil
}

func (e *SysDept) SetDeptLabel() (m []dto.DeptLabel, err error) {
	list := make([]models.SysDept, 0)
	err = e.Orm.Find(&list).Error
	if err != nil {
		log.Error("find dept list error, %s", err.Error())
		return
	}
	m = make([]dto.DeptLabel, 0)
	var item dto.DeptLabel
	for i := range list {
		if list[i].ParentId != 0 {
			continue
		}
		item = dto.DeptLabel{}
		item.Id = list[i].DeptId
		item.Label = list[i].DeptName
		deptInfo := deptLabelCall(&list, item)
		m = append(m, deptInfo)
	}
	return
}

// deptLabelCall
func deptLabelCall(deptList *[]models.SysDept, dept dto.DeptLabel) dto.DeptLabel {
	list := *deptList
	var mi dto.DeptLabel
	min := make([]dto.DeptLabel, 0)
	for j := 0; j < len(list); j++ {
		if dept.Id != list[j].ParentId {
			continue
		}
		mi = dto.DeptLabel{Id: list[j].DeptId, Label: list[j].DeptName, Children: []dto.DeptLabel{}}
		ms := deptLabelCall(deptList, mi)
		min = append(min, ms)
	}
	dept.Children = min
	return dept
}
