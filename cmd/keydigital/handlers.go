package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/byuoitav/common/status"
	"github.com/byuoitav/keydigital"
	"github.com/labstack/echo"
)

type Handlers struct {
	CreateVideoSwitcher func(string) *keydigital.VideoSwitcher
}

func (h *Handlers) RegisterRoutes(group *echo.Group) {
	vs := group.Group("/:address")

	// TODO singleflight?

	// get state
	vs.GET("/output/:output/input", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("Getting inputs")

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		inputs, err := vs.AudioVideoInputs(ctx)
		if err != nil {
			l.Printf("unable to get inputs: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		out := c.Param("output")

		in, ok := inputs[out]
		if !ok {
			l.Printf("invalid output %q requested", out)
			return c.String(http.StatusBadRequest, "invalid output")
		}

		l.Printf("Got inputs: %+v", inputs)
		return c.JSON(http.StatusOK, status.Input{
			Input: fmt.Sprintf("%v:%v", in, out),
		})
	})

	vs.GET("/info", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("Getting info")

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		info, err := vs.Info(ctx)
		if err != nil {
			l.Printf("unable to get inputs: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Got info: %+v", info)
		return c.JSON(http.StatusOK, info)
	})

	vs.GET("/healthy", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("Getting healthy")

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		if err := vs.Healthy(ctx); err != nil {
			l.Printf("not healthy: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Healthy!")
		return c.NoContent(http.StatusOK)
	})

	// set state
	vs.GET("/output/:output/input/:input", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)
		out := c.Param("output")
		in := c.Param("input")

		l.Printf("Setting AV input on %q to %q", out, in)

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		err := vs.SetAudioVideoInput(ctx, out, in)
		if err != nil {
			l.Printf("unable to set AV input: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Set AV input")
		return c.JSON(http.StatusOK, status.Input{
			Input: fmt.Sprintf("%v:%v", in, out),
		})
	})
}
