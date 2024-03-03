package handler

import (
	"api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type getGoodsResponse struct {
	Data []models.Goods `json:"goods"`
}

// @Summary Создать товар
// @Tags Товары
// @Description create goods
// @ID create-goods
// @Accept  json
// @Produce  json
// @Param input body models.GoodsInput true "Goods info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/create/{projectId} [post]
func (h *Handler) createGood(c *gin.Context) {
	projectId, err := strconv.ParseInt(c.Param("projectId"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, 3, "projectId не определен"+err.Error())
		return
	}

	var input models.GoodsInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, 3, "Ошибка входящих данных. "+err.Error())
		return
	}

	good, err := h.services.Goods.Create(projectId, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, 3, err.Error())
		return
	}

	c.JSON(http.StatusCreated, good)
}

func (h *Handler) updateGood(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, 3, "Параметр id должен быть цифрой")
		return
	}

	projectId, err := strconv.ParseInt(c.Query("projectId"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, 3, "Параметр projectId должен быть цифрой")
		return
	}

	var input models.GoodsInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, 3, "Ошибка входящих данных. "+err.Error())
		return
	}

	good, err := h.services.Goods.Update(id, projectId, input)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, 3, err.Error())
		return
	}

	c.JSON(http.StatusCreated, good)
}

func (h *Handler) deleteGood(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, 3, "Параметр id должен быть цифрой")
		return
	}

	projectId, err := strconv.ParseInt(c.Query("projectId"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, 3, "Параметр projectId должен быть цифрой")
		return
	}

	good, err := h.services.Goods.Delete(id, projectId)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, 3, err.Error())
		return
	}

	c.JSON(http.StatusOK, good)
}

func (h *Handler) getGoods(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, 3, "Параметр limit должен быть цифрой")
		return
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, 3, "Параметр offset должен быть цифрой")
		return
	}

	var metaOutput struct {
		Meta  models.Meta    `json:"meta"`
		Goods []models.Goods `json:"goods"`
	}

	goods, meta, err := h.services.Goods.GetList(limit, offset)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, 3, err.Error())
		return
	}

	metaOutput.Meta = meta
	metaOutput.Goods = goods

	c.JSON(http.StatusOK, metaOutput)
}

func (h *Handler) reprioritize(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, 3, "Параметр id должен быть цифрой")
		return
	}

	projectId, err := strconv.ParseInt(c.Query("projectId"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, 3, "Параметр projectId должен быть цифрой")
		return
	}

	var priorityUpdate struct {
		NewPriority int64 `json:"newPriority" binding:"required"`
	}

	if err := c.BindJSON(&priorityUpdate); err != nil {
		newErrorResponse(c, http.StatusBadRequest, 3, "Ошибка входящих данных. "+err.Error())
		return
	}

	goods, err := h.services.Goods.UpdatePriority(id, projectId, priorityUpdate.NewPriority)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, 3, err.Error())
		return
	}

	c.JSON(http.StatusOK, struct {
		Priorities []models.GoodsPriority `json:"priorities"`
	}{
		goods,
	})
}
