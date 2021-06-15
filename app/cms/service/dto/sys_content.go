package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/app/cms/models"

	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysContentSearch struct {
	dto.Pagination `search:"-"`
	CateId         string `form:"cateId" search:"type:exact;column:cate_id;table:sys_content" comment:"分类"`
	Name           string `form:"name" search:"type:contains;column:name;table:sys_content" comment:"名称"`
	Status         string `form:"status" search:"type:exact;column:status;table:sys_content" comment:"状态"`
}

func (m *SysContentSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysContentSearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

func (m *SysContentSearch) Generate() dto.Index {
	o := *m
	return &o
}

type SysContentControl struct {
	Id      int    `uri:"Id" comment:""` //
	CateId  int    `json:"cateId" comment:""`
	Name    string `json:"name" comment:""`
	Status  int    `json:"status" comment:""`
	Img     string `json:"img" comment:""`
	Content string `json:"content" comment:""`
	Remark  string `json:"remark" comment:""`
	Sort    int    `json:"sort" comment:""`
}

func (s *SysContentControl) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBindUri(s)
	if err != nil {
		log.Debugf("ShouldBindUri error: %s", err.Error())
		return err
	}
	err = ctx.ShouldBind(s)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

func (s *SysContentControl) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysContentControl) GenerateM() (common.ActiveRecord, error) {
	return &models.SysContent{
		Model:   common.Model{s.Id},
		CateId:  s.CateId,
		Name:    s.Name,
		Status:  s.Status,
		Img:     s.Img,
		Content: s.Content,
		Remark:  s.Remark,
		Sort:    s.Sort,
	}, nil
}

func (s *SysContentControl) GetId() interface{} {
	return s.Id
}

type SysContentById struct {
	dto.ObjectById
}

func (s *SysContentById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysContentById) GenerateM() (common.ActiveRecord, error) {
	return &models.SysContent{}, nil
}
