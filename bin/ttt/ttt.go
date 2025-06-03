package main

import (
	"path/filepath"
	"task-time-tracker/lib/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/rs/zerolog/log"
)

func main() {
	var here string = utils.GetHereDirExe()
    var e error

    var webBuildDir string=filepath.Join(here,"../../task-time-tracker-web/build")

    var app *fiber.App=fiber.New(fiber.Config{
        CaseSensitive: true,
        ErrorHandler: func(c fiber.Ctx, err error) error {
            log.Warn().Msg("fiber error")
            log.Warn().Msgf("%v",err)
            return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
        },
    })

    app.Use("/*",static.New(webBuildDir))

    e=app.Listen(":4602")

    if e!=nil {
        panic(e)
    }
}