package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mauricioryan/pagosgo/pago"
	"github.com/mauricioryan/pagosgo/tools/errors"
)

type permission struct {
	Permissions []string `json:"permissions" binding:"required"`
}

// Pagos Devuelve una lista de todos los pagos
/**
 * @api {get} /v1/pagosu Listar Pagos
 * @apiName Listar Pagos
 * @apiGroup Pagos
 *
 * @apiDescription Obtiene información de todos los pagos.
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *     [{
 *        "id": "{Id pago}",
 *        fecha, autorizado, webhook,
 * 	      "autorizado": true|false
 *     }, ...]
 *
 * @apiUse no se que va aca
 */
func Pagos(c *gin.Context) {
	//dejar la autorizacion para mas adelante
	payload, err := validateAuthHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	pagoService, err := pago.NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	if !pagoService.Granted(payload.UserID.Hex(), "admin") {
		errors.Handle(c, errors.AccessLevel)
		return
	}

	user, err := pagoService.Pagos()

	if err != nil {
		errors.Handle(c, err)
		return
	}
	result := []gin.H{}
	for i := 0; i < len(user); i = i + 1 {
		result = append(result, gin.H{
			"id":         pago[i].ID.Hex(),
			"fechaPago":  pago[i].FechaPago,
			"autorizado": pago[i].Autorizado,
		})
	}

	c.JSON(200, result)
}
