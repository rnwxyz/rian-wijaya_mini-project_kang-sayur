package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/region/service"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
)

type regionController struct {
	service service.RegionService
}

func NewRegionController(regionService service.RegionService) *regionController {
	return &regionController{
		service: regionService,
	}
}

func (r *regionController) InitRoute(auth *echo.Group) {
	regions := auth.Group("/regions")
	regions.GET("/provinces", r.GetProvince)
	regions.GET("/regencies/:province_id", r.GetRegency)
	regions.GET("/districts/:regency_id", r.GetDistrict)
	regions.GET("/villages/:district_id", r.GetVillage)
}

func (r *regionController) GetProvince(c echo.Context) error {
	provinces, err := r.service.FindProvince(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get provinces success",
		"data":    provinces,
	})
}

func (r *regionController) GetRegency(c echo.Context) error {
	id := c.Param("province_id")
	regencies, err := r.service.FindRegency(&id, c.Request().Context())
	if err != nil {
		if err == customerrors.ErrInvalidId {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get regencies success",
		"data":    regencies,
	})
}

func (r *regionController) GetDistrict(c echo.Context) error {
	id := c.Param("regency_id")
	districts, err := r.service.FindDistrict(&id, c.Request().Context())
	if err != nil {
		if err == customerrors.ErrInvalidId {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get districts success",
		"data":    districts,
	})
}

func (r *regionController) GetVillage(c echo.Context) error {
	id := c.Param("district_id")
	villages, err := r.service.FindVillage(&id, c.Request().Context())
	if err != nil {
		if err == customerrors.ErrInvalidId {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "get villages success",
		"data":    villages,
	})
}
