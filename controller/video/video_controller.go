package c_video

import (
	"apiproject/cache"
	"apiproject/dao"
	d_video "apiproject/dao/video"
	"apiproject/log"
	m_video "apiproject/model/video"
	s_video "apiproject/service/video"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

/**
查询视频列表接口
*/

/************************start swagger api定义注解 **************/
// @Summary 查询视频列表接口
// @Description 查询视频列表接口
// @Tags 视频
// @Accept  json
// @Produce  json
// @Success 200 {object} gin.H
// @Router /api/video/findVideoList [get]
/************************end swagger api定义注解 **************/
func FindVideoList(ctx *gin.Context) {
	//先读取本地缓存
	videoList, exist := cache.LocalCache.Get("FindVideoList")
	if exist {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": videoList,
		})
		return
	}

	videoList, err := s_video.VideoService.FindVideoList()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": videoList,
	})

	//响应数据缓存到本地
	cache.LocalCache.SetDefault("FindVideoList", videoList)
}

/**
分页查询视频列表接口
*/

/************************start swagger api定义注解 **************/
// @Summary 分页查询视频列表接口
// @Description 分页查询视频列表接口
// @Tags 视频
// @Accept  json
// @Produce  json
// @Param page query ReqPage false "查询参数"
// @Success 200 {object} gin.H
// @Router /api/video/findVideoListPage [get]
/************************end swagger api定义注解 **************/
func FindVideoListPage(ctx *gin.Context) {
	page := ReqPage{}
	if err := ctx.ShouldBind(&page); err != nil {
		log.Logger.Error("绑定请求参数到对象异常", zap.Error(err))
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}
	log.Logger.Info("绑定请求参数到对象", zap.Any("page", page))

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": s_video.VideoService.FindVideoListPage(page.PageNo, page.PageSize),
	})
}

/**
按条件查询视频列表接口
*/

/************************start swagger api定义注解 **************/
// @Summary 按条件查询视频列表接口
// @Description 按条件查询视频列表接口
// @Tags 视频
// @Accept  json
// @Produce  json
// @Param video query ReqFindVideoByWhere true "查询参数"
// @Success 200 {object} gin.H
// @Router /api/video/findVideoByWhere [get]
/************************end swagger api定义注解 **************/
func FindVideoByWhere(ctx *gin.Context) {
	/**
	绑定query参数到对象.
	注意:默认情况下, "query参数或form参数"和对象成员的大小写需保持一致. 除非为对象补充tag定义, 定义规则如下:
	1. 对于query参数, 需要如下定义
	form:"siteId" binding:"required"
	2. 对于form参数, 需要如下定义
	form:"siteId"
	3. 对于body的json参数, 无需定义tag, 且不区分大小写

	建议用ShouldBind, 不要使用其他bind方法, 否则可能会返回400错误, 参考链接:https://learnku.com/docs/gin-gonic/2018/gin-readme/3819
	*/
	para := ReqFindVideoByWhere{}
	if err := ctx.ShouldBind(&para); err != nil {
		log.Logger.Error("绑定请求参数到对象异常", zap.Error(err))
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": nil,
			"msg":  "参数错误",
		})
		return
	}
	log.Logger.Info("绑定请求参数到对象", zap.Any("videoQuery", para))

	videoList := []*m_video.Video{}
	if err := d_video.VideoDao.FindList(dao.Db.Where("site_id = ? and title like %?%", para.SiteId, para.TitleLike), &videoList); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": videoList,
	})
}

/**
添加视频
*/

/************************start swagger api定义注解 **************/
// @Summary 添加视频
// @Description 添加视频
// @Tags 视频
// @Accept  json
// @Produce  json
// @Success 200 {object} gin.H
// @Router /api/video/addVideo [post]
/************************end swagger api定义注解 **************/
func AddVideo(ctx *gin.Context) {
	video := &m_video.Video{
		Title: "test333",
	}
	if err := d_video.VideoDao.Insert(dao.Db, video); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": nil,
	})
}

/**
批量添加视频
*/

/************************start swagger api定义注解 **************/
// @Summary 批量添加视频
// @Description 批量添加视频
// @Tags 视频
// @Accept  json
// @Produce  json
// @Success 200 {object} gin.H
// @Router /api/video/bulkAddVideo [post]
/************************end swagger api定义注解 **************/
func BulkAddVideo(ctx *gin.Context) {
	s_video.VideoService.BulkAddVideo()

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": nil,
	})
}

/**
更新视频
*/

/************************start swagger api定义注解 **************/
// @Summary 更新视频
// @Description 更新视频
// @Tags 视频
// @Accept  json
// @Produce  json
// @Param video body m_video.Video true "json对象"
// @Success 200 {object} gin.H
// @Router /api/video/updateVideo [post]
/************************end swagger api定义注解 **************/
func UpdateVideo(ctx *gin.Context) {
	//绑定参数到对象
	video := m_video.Video{}
	if err := ctx.ShouldBind(&video); err != nil {
		log.Logger.Error("绑定请求参数到对象异常", zap.Error(err))
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": nil,
			"msg":  "参数错误",
		})
		return
	}
	log.Logger.Info("绑定请求参数到对象", zap.Any("video", video))

	if err := d_video.VideoDao.Update(dao.Db.Model(&video).Where("id = ?", "id4"), &video); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": nil,
	})
}

/**
删除视频
*/

/************************start swagger api定义注解 **************/
// @Summary 删除视频
// @Description 删除视频
// @Tags 视频
// @Accept  json
// @Produce  json
// @Param videoId query ReqVideoId false "查询参数"
// @Success 200 {object} gin.H
// @Router /api/video/deleteVideo [delete]
/************************end swagger api定义注解 **************/
func DeleteVideo(ctx *gin.Context) {
	//绑定参数到对象
	video := ReqVideoId{}
	if err := ctx.ShouldBind(&video); err != nil {
		log.Logger.Error("绑定请求参数到对象异常", zap.Error(err))
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": nil,
			"msg":  "参数错误",
		})
		return
	}
	log.Logger.Info("绑定请求参数到对象", zap.Any("video", video))

	if err := d_video.VideoDao.Delete(dao.Db.Where("id = ?", video.ID), &m_video.Video{}); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": nil,
	})
}

/**
查询视频总数
*/

/************************start swagger api定义注解 **************/
// @Summary 查询视频总数
// @Description 查询视频总数
// @Tags 视频
// @Accept  json
// @Produce  json
// @Success 200 {object} gin.H
// @Router /api/video/getCount [get]
/************************end swagger api定义注解 **************/
func GetCount(ctx *gin.Context) {
	totalCount := d_video.VideoDao.GetCount(dao.Db.Model(&m_video.Video{}))

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": totalCount,
	})
}
